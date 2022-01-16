package requesthandlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func EventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	uuid := ps.ByName("uuid")
	err := t.Execute(w, IndexData{Uuid: uuid})
	if err != nil {
		log.Println("Failed to parse index template.")
	}
}
