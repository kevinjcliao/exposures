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
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseFiles("templates/index.html"))
	err := t.Execute(w, IndexData{Uuid: uuid.New().String()})
	if err != nil {
		log.Println("Failed to parse index template.")
	}
}
