package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/database"
	"github.com/gimaevra94/test-business-base/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed templates/*.html
var templatesFS embed.FS

func main() {
	db, err := initDB()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.DB.Close()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))

	r := initRouter(db, tmpl)
	if err := http.ListenAndServe(":8081", r); err != nil {
		logrus.Fatal(err)
	}

}

func initDB() (*database.DB, error) {
	db, err := database.DBConn("root:root@tcp(localhost:3306)/repair_service")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}

func initRouter(db *database.DB, tmpl *template.Template) *chi.Mux {
	r := chi.NewRouter()
	r.Get(consts.HomePath, handlers.Home(tmpl))
	r.Get(consts.GetLoginPath, handlers.GetLogin(db, tmpl))
	r.Post(consts.LoginPath, handlers.Login(db, tmpl))
	r.Get(consts.LogoutPath, handlers.Logout(db, tmpl))
	r.Get(consts.CreatePath, handlers.Create(db, tmpl))
	r.Post(consts.CreatePath, handlers.Create(db, tmpl))
	r.Get(consts.DashboardPath, handlers.Dashboard(db, tmpl))
	r.Post(consts.ActionPath, handlers.Action(db, tmpl))
	return r
}
