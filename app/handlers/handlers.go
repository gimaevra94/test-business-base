package handlers

import (
	"net/http"
	"strconv"
)

func getSession(r *http.Request) (int, string, string) {
	c, err := r.Cookie("user_id")
	if err != nil {
		return 0, "", ""
	}
	uid, _ := strconv.Atoi(c.Value)
	c, err = r.Cookie("role")
	role := ""
	if err == nil {
		role = c.Value
	}
	c, err = r.Cookie("name")
	name := ""
	if err == nil {
		name = c.Value
	}
	return uid, role, name
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	uid, role, name := getSession(r)
	if uid > 0 {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uid, _ := strconv.Atoi(r.FormValue("user_id"))
		role := r.FormValue("role")
		name := r.FormValue("name")
		http.SetCookie(w, &http.Cookie{Name: "user_id", Value: strconv.Itoa(uid)})
		http.SetCookie(w, &http.Cookie{Name: "role", Value: role})
		http.SetCookie(w, &http.Cookie{Name: "name", Value: name})
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	// Load users for select
	rows, _ := db.Query("SELECT id, name, role FROM users")
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Role)
		users = append(users, u)
	}
	tmpl.ExecuteTemplate(w, "login.html", users)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_, err := db.Exec("INSERT INTO requests (client_name, phone, address, problem_text, status) VALUES (?, ?, ?, ?, 'new')",
			r.FormValue("client_name"), r.FormValue("phone"), r.FormValue("address"), r.FormValue("problem_text"))
		if err == nil { http.Redirect(w, r, "/dashboard", http.StatusSeeOther); return }
	}
	tmpl.ExecuteTemplate(w, "create.html", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	uid, role, name := getSession(r)
	if uid == 0 { http.Redirect(w, r, "/", http.StatusSeeOther); return }

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
	if err != nil { http.Error(w, err.Error(), 500); return }
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
		if role != "dispatcher" { break }
		mid := r.FormValue("master_id")
		res, err = db.Exec("UPDATE requests SET status = 'assigned', assigned_to = ? WHERE id = ? AND status = 'new'", mid, id)
	case "cancel":
		if role != "dispatcher" { break }
		res, err = db.Exec("UPDATE requests SET status = 'canceled' WHERE id = ? AND status IN ('new', 'assigned')", id)
	case "start":
		if role != "master" { break }
		// Race Condition Safe: Optimistic Locking
		res, err = db.Exec("UPDATE requests SET status = 'in_progress', version = version + 1 WHERE id = ? AND status = 'assigned' AND assigned_to = ?", id, uid)
		if err == nil {
			rows, _ := res.RowsAffected()
			if rows == 0 { err = sql.ErrNoRows } // Force conflict
		}
	case "finish":
		if role != "master" { break }
		res, err = db.Exec("UPDATE requests SET status = 'done' WHERE id = ? AND status = 'in_progress' AND assigned_to = ?", id, uid)
	}

	if err != nil || (action == "start" && err == sql.ErrNoRows) {
		http.Redirect(w, r, "/dashboard?error=conflict", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}