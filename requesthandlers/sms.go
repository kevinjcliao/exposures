package requesthandlers

import (
	"context"
	"exposures/ent"
	"exposures/messages"
	"exposures/smshandler"
	"exposures/twilio"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func SmsHandler(ctx context.Context, entClient *ent.Client) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse request.")
		}
		body := r.Form.Get("Body")
		from := r.Form.Get("From")

		result := []messages.Message{}
		if strings.ToUpper(body) == "POSITIVE" {
			result = append(result, smshandler.HandlePositiveCase(ctx, entClient, from)...)
		} else if _, err := uuid.Parse(body); err == nil {
			result = append(result, smshandler.Rsvp(ctx, entClient, from, body)...)
		} else {
			result = append(result, smshandler.Error(ctx, entClient, from))
		}
		for _, x := range result {
			twilio.SendSms(x.Recipient, x.Message)
		}
	}
}
