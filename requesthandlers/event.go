package requesthandlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func EventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	err := t.Execute(w, IndexData{Uuid: uuid.New().String()})
	if err != nil {
		log.Println("Failed to parse index template.")
	}
}
