package errs

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gimaevra94/test-business-base/structs"
	"github.com/sirupsen/logrus"
)

func RenderError(w http.ResponseWriter, tmpl *template.Template, templateName string, data *structs.LoginData, msg string, err error) {
	data.Msg = msg
	log.Printf("%+v", err)
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		logrus.Error(err)
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error rendering template"))
	}
}
