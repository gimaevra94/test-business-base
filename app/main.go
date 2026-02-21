package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/database"
	"github.com/gimaevra94/test-business-base/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed templates/*.html
var templatesFS embed.FS

func main() {
	if err := initEnv(); err != nil {
		logrus.Fatal(err)
		return
	}

	db, err := initDB()
	if err != nil {
		logrus.Fatal(err)
		return
	}
	defer db.DB.Close()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))

	r := initRouter(db, tmpl)
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatal(err)
		return
	}

}

func initEnv() error {
	if err := godotenv.Load("../.env"); err != nil {
		return errors.WithStack(err)
	}

	envVars := []string{
		"",
		"",
	}

	for _, v := range envVars {
		if os.Getenv(v) == "" {
			spErr := fmt.Sprintf("Missing %s env var", v)
			err := errors.New(spErr)
			return errors.WithStack(err)
		}
	}
}

func initDB() (*database.DB, error) {
	driver, cfg := os.Getenv("CONNECTION_CFG"), os.Getenv("SQL_DRIVER")
	db, err := database.DBConn(driver, cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}

func initRouter(db *database.DB, tmpl *template.Template) *chi.Mux {
	r := chi.NewRouter()
	r.Get(consts.Home, handlers.Home(tmpl))
	r.MethodFunc("GET|POST", consts.Login, handlers.Login(db, tmpl))
	r.Get(consts.Logout, handlers.Logout(db, tmpl))
	r.MethodFunc("GET|POST", consts.Create, handlers.Create(db, tmpl))
	r.Get("/dashboard", dashboardHandler(db))
	r.Post("/action", actionHandler(db))
	return r
}
