// [START gae_go111_app]

// Sample exposures is an App Engine app.
package main

// [START import].
import (
	"context"
	"exposures/db"
	"exposures/env"
	"exposures/requesthandlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// [END import]
// [START main_func]

func SmsEndpoint() string {
	return fmt.Sprintf("/%s", env.MustGetenv("SMS_ENDPOINT"))
}

func main() {
	client, err := db.BuildEntClient()
	if err != nil {
		log.Fatalf("Error creating ent client: %v", err)
	}

	defer client.Close()

	router := httprouter.New()
	router.GET("/", requesthandlers.IndexHandler)
	router.GET("/event/:uuid", requesthandlers.EventHandler)
	router.POST(SmsEndpoint(), requesthandlers.SmsHandler(context.Background(), client))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

// [END main_func]

// [START Handle Positive Notification].
