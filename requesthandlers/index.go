package requesthandlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/google/uuid"
)

// indexHandler responds to requests with our greeting.
type IndexData struct {
	Uuid string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
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
