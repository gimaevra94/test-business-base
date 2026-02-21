package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gimaevra94/test-business-base/database"
	"github.com/gimaevra94/test-business-base/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

//go:embed templates/*.html
var templatesFS embed.FS

func main() {
	initEnv()

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))

	r := initRouter(db)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}

func initEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
		return
	}

	envVars := []string{
		"",
		"",
	}

	for _, v := range envVars {
		if os.Getenv(v) == "" {
			log.Fatal(v)
			return
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
	r.Get("/", handlers.Home(tmpl))
	r.MethodFunc("GET|POST", "/login", loginHandler(db))
	r.Get("/logout", logoutHandler(db))
	r.MethodFunc("GET|POST", "/create", createHandler(db))
	r.Get("/dashboard", dashboardHandler(db))
	r.Post("/action", actionHandler(db))
	return r
}
