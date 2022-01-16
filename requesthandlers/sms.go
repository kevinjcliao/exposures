package requesthandlers

import (
	"context"
	"exposures/ent"
	"exposures/ent/checkin"
	"exposures/ent/user"
	"exposures/messages"
	"exposures/twilio"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handlePositiveCase(ctx context.Context, client *ent.Client, from string) {
	twilio.SendSms(from, "Thanks for letting us know. Please take care of yourself and those around you. We hope you recover soon.")
	fourteenDaysAgo := time.Now().Unix() - (14 * 24 * 60 * 60)
	threeHours := int64(3 * 60 * 60)
	checkins := client.Checkin.Query().
		Where(
			checkin.And(
				checkin.HasSenderWith(user.PhoneNumberEQ(from)),
				checkin.CheckinTimeGT(fourteenDaysAgo),
			),
		).
		AllX(ctx)
	for _, positiveCheckin := range checkins {
		similarCheckins := client.Checkin.Query().Where(
			checkin.EventIDEQ(positiveCheckin.EventID),
		).
			Where(
				checkin.HasSenderWith(user.PhoneNumberNEQ(from)),
			).
			Where(
				checkin.Or(
					checkin.CheckinTimeGT(positiveCheckin.CheckinTime-threeHours),
					checkin.CheckinTimeLT(positiveCheckin.CheckinTime+threeHours),
				),
			).
			AllX(ctx)
		phoneNumbers := make(map[string]bool)
		for _, similarCheckin := range similarCheckins {
			sender := similarCheckin.QuerySender().OnlyX(ctx)
			phoneNumbers[sender.PhoneNumber] = true
		}
		for phone, _ := range phoneNumbers {
			date := time.Unix(positiveCheckin.CheckinTime, 0)
			body := fmt.Sprintf("Someone at an event you attended on: %s tested positive for COVID-19. You should get tested.", date)
			twilio.SendSms(phone, body)
		}
	}
}

func SmsHandler(ctx context.Context, entClient *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse request.")
		}
		body := r.Form.Get("Body")
		from := r.Form.Get("From")

		if body == "POSITIVE" {
			handlePositiveCase(ctx, entClient, from)
			return
		}

		rsvpMessage := messages.MessageBodies[messages.RsvpSuccess]
		twilio.SendSms(from, rsvpMessage)
		user, err := entClient.User.Query().Where(user.PhoneNumberEQ(from)).Only(ctx)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				user = entClient.User.Create().SetPhoneNumber(from).SaveX(ctx)

			default:
				log.Println("Huh")
			}
		}
		_, err = entClient.Checkin.Create().
			SetCheckinTime(time.Now().Unix()).
			SetEventID(body).
			SetSender(user).
			Save(ctx)
		if err != nil {
			log.Printf("Error creating checkin: %v", err)
		}

		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println("SMS Sent Successfully!")
		}
	}
}
