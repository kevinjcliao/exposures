package requesthandlers

import (
	"fmt"
	"log"
	"net/http"

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
	log.Println("Redirecting.")
	http.Redirect(w, r, fmt.Sprintf("/event/%s", uuid.New().String()), http.StatusMovedPermanently)
}
