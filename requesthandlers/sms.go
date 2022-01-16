package requesthandlers

import (
	"context"
	"exposures/ent"
	"exposures/messages"
	"exposures/smshandler"
	"exposures/twilio"
	"log"
	"net/http"
)

func SmsHandler(ctx context.Context, entClient *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse request.")
		}
		body := r.Form.Get("Body")
		from := r.Form.Get("From")

		result := []messages.Message{}
		if body == "POSITIVE" {
			result = append(result, smshandler.HandlePositiveCase(ctx, entClient, from)...)
		} else {
			result = append(result, smshandler.Rsvp(ctx, entClient, from, body)...)
		}
		for _, x := range result {
			twilio.SendSms(x.Recipient, x.Message)
		}
	}
}
