package requesthandlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// indexHandler responds to requests with our greeting.
type IndexData struct {
	Uuid string
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	uuid := uuid.New().String()
	err := t.Execute(w, IndexData{Uuid: uuid})
	if err != nil {
		log.Println("Failed to parse template.")
	}
}
