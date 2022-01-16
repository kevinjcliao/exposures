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
	id := ps.ByName("uuid")

	if _, err := uuid.Parse(id); err != nil {
		t = template.Must(template.ParseFiles("templates/404.html"))
		err := t.Execute(w, IndexData{Uuid: id})
		if err != nil {
			log.Println("Failed to parse 404 template.")
		}
		return
	}
	err := t.Execute(w, IndexData{Uuid: id})
	if err != nil {
		log.Println("Failed to parse index template.")
	}
}
