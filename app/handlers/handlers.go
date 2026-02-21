package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/database"
	"github.com/gimaevra94/test-business-base/structs"
	"github.com/sirupsen/logrus"
)

func getSession(r *http.Request) (int, string, string) {
	c, err := r.Cookie(consts.UID)
	if err != nil {
		return 0, "", ""
	}
	uid, _ := strconv.Atoi(c.Value)

	c, err = r.Cookie(consts.Role)
	role := ""
	if err == nil {
		role = c.Value
	}

	c, err = r.Cookie(consts.Name)
	name := ""
	if err == nil {
		name = c.Value
	}
	return uid, role, name
}

func Home(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, _, _ := getSession(r)
		if uid > 0 {
			http.Redirect(w, r, consts.Dashboard, http.StatusSeeOther)
			return
		}
		tmpl.ExecuteTemplate(w, consts.LoginHTML, nil)
	}
}

func Login(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := structs.LoginData{}

		if r.Method == http.MethodPost {
			StrUID := r.FormValue(consts.UID)
			role := r.FormValue(consts.Role)
			name := r.FormValue(consts.Name)

			if StrUID == "" || role == "" || name == "" {
				data.Msg = consts.BadInput
				logrus.Info(consts.BadInput)
				return
			}

			uid, _ := strconv.Atoi(r.FormValue(consts.UID))

			http.SetCookie(w, &http.Cookie{Name: consts.UID, Value: strconv.Itoa(uid)})
			http.SetCookie(w, &http.Cookie{Name: consts.Role, Value: role})
			http.SetCookie(w, &http.Cookie{Name: consts.Name, Value: name})
			http.Redirect(w, r, consts.Dashboard, http.StatusSeeOther)
			return

		} else if r.Method == http.MethodGet {
			users, err := db.GetUsers()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					data.Msg = consts.EmptyDB
					logrus.Info(consts.EmptyDB)
				}
				data.Msg = consts.InternalError
				logrus.Error(err)
			}

			data.Users = users
			tmpl.ExecuteTemplate(w, consts.Login, data)
			return

		} else {
			data.Msg = consts.NotAllowed
			logrus.Info(consts.NotAllowed)
			return
		}
	}
}

func Logout(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: consts.UID, Value: "", MaxAge: -1})
		http.Redirect(w, r, consts.Home, http.StatusSeeOther)
	}
}

func Create(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := structs.LoginData{}

		if r.Method == http.MethodPost {
			clientName := r.FormValue(consts.ClientName)
			phone := r.FormValue(consts.Phone)
			address := r.FormValue(consts.Address)
			problemText := r.FormValue(consts.ProblemText)

			if clientName == "" || phone == "" || address == "" || problemText == "" {
				data.Msg = consts.BadInput
				logrus.Info(consts.BadInput)
				return
			}

			req := structs.Request{
				ClientName:  clientName,
				Phone:       phone,
				Address:     address,
				ProblemText: problemText,
			}

			err, insertIsOk := db.Create(&req)
			if err != nil {
				data.Msg = consts.InternalError
				logrus.Error(err)
				return
			}

			if !insertIsOk {
				data.Msg = consts.InternalError
				logrus.Error(consts.InternalError)
				return
			}

			http.Redirect(w, r, consts.Dashboard, http.StatusSeeOther)
			return
		}
		tmpl.ExecuteTemplate(w, "create.html", nil)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	uid, role, name := getSession(r)
	if uid == 0 {
		http.Redirect(w, r, consts.Home, http.StatusSeeOther)
		return
	}

	query := "SELECT id, client_name, phone, address, problem_text, status, assigned_to, version, created_at, updated_at FROM requests WHERE 1=1"
	args := []interface{}{}
	if role == "master" {
		query += " AND assigned_to = ?"
		args = append(args, uid)
	}
	if status := r.URL.Query().Get("status"); status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var reqs []Request
	for rows.Next() {
		var r Request
		rows.Scan(&r.ID, &r.ClientName, &r.Phone, &r.Address, &r.ProblemText, &r.Status, &r.AssignedTo, &r.Version, &r.CreatedAt, &r.UpdatedAt)
		reqs = append(reqs, r)
	}

	masters, _ := db.Query("SELECT id, name FROM users WHERE role = 'master'")
	defer masters.Close()
	var ms []User
	for masters.Next() {
		var u User
		masters.Scan(&u.ID, &u.Name)
		ms = append(ms, u)
	}

	data := map[string]interface{}{"User": name, "Role": role, "Requests": reqs, "Masters": ms}
	tmpl.ExecuteTemplate(w, "dashboard.html", data)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	uid, role, _ := getSession(r)
	id, _ := strconv.Atoi(r.FormValue("id"))
	action := r.FormValue("action")

	var res sql.Result
	var err error

	switch action {
	case "assign":
		if role != "dispatcher" {
			break
		}
		mid := r.FormValue("master_id")
		res, err = db.Exec("UPDATE requests SET status = 'assigned', assigned_to = ? WHERE id = ? AND status = 'new'", mid, id)
	case "cancel":
		if role != "dispatcher" {
			break
		}
		res, err = db.Exec("UPDATE requests SET status = 'canceled' WHERE id = ? AND status IN ('new', 'assigned')", id)
	case "start":
		if role != "master" {
			break
		}
		// Race Condition Safe: Optimistic Locking
		res, err = db.Exec("UPDATE requests SET status = 'in_progress', version = version + 1 WHERE id = ? AND status = 'assigned' AND assigned_to = ?", id, uid)
		if err == nil {
			rows, _ := res.RowsAffected()
			if rows == 0 {
				err = sql.ErrNoRows
			} // Force conflict
		}
	case "finish":
		if role != "master" {
			break
		}
		res, err = db.Exec("UPDATE requests SET status = 'done' WHERE id = ? AND status = 'in_progress' AND assigned_to = ?", id, uid)
	}

	if err != nil || (action == "start" && err == sql.ErrNoRows) {
		http.Redirect(w, r, "/dashboard?error=conflict", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
