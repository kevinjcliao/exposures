package twilio

import (
	"log"

	"exposures/env"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(to string, body string) {
	if env.MustGetenv("EXPOSURES_ENV") != "PROD" {
		log.Println("Sending an SMS to: ", to, " with body: ", body)
		log.Println("Not actually sending an SMS because you're in test.")
		return
	}
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: env.MustGetenv("TWILIO_SID"),
		Password: env.MustGetenv("TWILIO_API_KEY"),
	})
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom("+1(417)668-2737")
	params.SetBody(body)
	_, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		log.Fatalf("Failed to send sms to: %s containing body: %s", to, body)
	}
}
