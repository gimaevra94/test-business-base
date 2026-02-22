package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/database"
	"github.com/gimaevra94/test-business-base/errs"
	"github.com/gimaevra94/test-business-base/structs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getSession(r *http.Request) (int, string, string, error) {
	c, err := r.Cookie(consts.ID)
	if err != nil {
		return 0, "", "", errors.WithStack(err)
	}

	id, err := strconv.Atoi(c.Value)
	if err != nil {
		return 0, "", "", errors.WithStack(err)
	}

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

	return id, role, name, nil
}

func Home(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _, _, _ := getSession(r)
		if id > 0 {
			http.Redirect(w, r, consts.DashboardPath, http.StatusSeeOther)
			return
		}

		if err := tmpl.ExecuteTemplate(w, consts.LoginHTML, nil); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func Login(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := structs.LoginData{}

		switch r.Method {
		case http.MethodPost:
			id := r.FormValue(consts.ID)
			role := r.FormValue(consts.Role)
			name := r.FormValue(consts.Name)

			if id == "" || role == "" || name == "" {
				errs.RenderError(w, tmpl, consts.LoginHTML, &data, consts.BadInputMsg, errors.WithStack(errors.New(consts.BadInputMsg)))
				return
			}

			http.SetCookie(w, &http.Cookie{Name: consts.ID, Value: id})
			http.SetCookie(w, &http.Cookie{Name: consts.Role, Value: role})
			http.SetCookie(w, &http.Cookie{Name: consts.Name, Value: name})
			http.Redirect(w, r, consts.DashboardPath, http.StatusSeeOther)
			return

		case http.MethodGet:
			log.Print("get")
			users, err := db.GetUsers()
			if err != nil {
				errs.RenderError(w, tmpl, consts.LoginHTML, &data, consts.InternalErrorMsg, err)
				return
			}
			log.Print(users)
			data.Users = users
			if err := tmpl.ExecuteTemplate(w, consts.LoginHTML, data); err != nil {
				logrus.Error(err)
				return
			}
			return

		default:
			errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.NotAllowedMsg, errors.WithStack(errors.New(consts.NotAllowedMsg)))
			return
		}
	}
}

func GetLogin(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := structs.LoginData{}
		log.Print("get")
		users, err := db.GetUsers()
		if err != nil {
			errs.RenderError(w, tmpl, consts.LoginHTML, &data, consts.InternalErrorMsg, err)
			return
		}
		log.Print(users)
		data.Users = users
		if err := tmpl.ExecuteTemplate(w, consts.LoginHTML, data); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func Logout(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: consts.ID, Value: "", MaxAge: -1, Path: "/"})
		http.Redirect(w, r, consts.HomePath, http.StatusSeeOther)
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
				errs.RenderError(w, tmpl, consts.CreateHTML, &data, consts.BadInputMsg, errors.WithStack(errors.New(consts.BadInputMsg)))
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
				data.Msg = consts.InternalErrorMsg
				logrus.Error(err)
				if err := tmpl.ExecuteTemplate(w, consts.CreateHTML, data); err != nil {
					logrus.Error(err)
					return
				}
				errs.RenderError(w, tmpl, consts.CreateHTML, &data, consts.InternalErrorMsg, err)
				return
			}

			if !insertIsOk {
				errs.RenderError(w, tmpl, consts.CreateHTML, &data, consts.InternalErrorMsg, errors.WithStack(errors.New(consts.InternalErrorMsg)))
				return
			}

			http.Redirect(w, r, consts.DashboardPath, http.StatusSeeOther)
			return
		}

		if err := tmpl.ExecuteTemplate(w, consts.CreateHTML, data); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func Dashboard(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := structs.LoginData{}

		uid, role, name, err := getSession(r)
		if err != nil {
			errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, err)
			return
		}

		if uid == 0 {
			http.Redirect(w, r, consts.HomePath, http.StatusSeeOther)
			return
		}

		query := consts.DashboardSelectQuery
		args := []any{}
		if role == consts.Master {
			query += " AND assigned_to = ?"
			args = append(args, uid)
		}

		if status := r.URL.Query().Get(consts.Status); status != "" {
			query += " AND status = ?"
			args = append(args, status)
		}

		reqs, masters, err := db.Dashboard(query, args)
		if err != nil {
			errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, errors.WithStack(errors.New(consts.InternalErrorMsg)))
			return
		}

		dashboardData := map[string]interface{}{consts.User: name, consts.Role: role, consts.Requests: reqs, consts.Masters: masters}
		tmpl.ExecuteTemplate(w, consts.DashboardHTML, dashboardData)
	}
}

func Action(db *database.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := structs.LoginData{}

		_, role, _, cookieErr := getSession(r)
		if cookieErr != nil {
			errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, cookieErr)
			return
		}

		action := r.FormValue(consts.Action)
		StID := r.FormValue(consts.ID)

		if action == "" || StID == "" {
			errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.BadInputMsg, errors.WithStack(errors.New(consts.BadInputMsg)))
			return
		}

		id, strconvErr := strconv.Atoi(StID)
		if strconvErr != nil {
			errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, strconvErr)
			return
		}

		var err error
		var progressErr error

		switch action {
		case consts.Assign:
			if role != consts.Dispatcher {
				break
			}

			if err = db.AssignedStatusUpdate(id); err != nil {
				errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, err)
				return
			}

		case consts.Cancel:
			if role != consts.Dispatcher {
				break
			}

			if err = db.CanceledStatusUpdate(id); err != nil {
				errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, err)
				return
			}

		case consts.Start:
			if role != consts.Master {
				break
			}

			if progressErr = db.InProgressStatusUpdate(id); progressErr != nil {
				errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, progressErr)
				return
			}

		case consts.Finish:
			if role != consts.Master {
				break
			}

			if err = db.DoneStatusUpdate(id); err != nil {
				errs.RenderError(w, tmpl, consts.DashboardHTML, &data, consts.InternalErrorMsg, err)
				return
			}

		}

		if progressErr != nil {
			http.Redirect(w, r, consts.DashboardPath+"?error=conflict", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, consts.DashboardPath, http.StatusSeeOther)
	}
}
