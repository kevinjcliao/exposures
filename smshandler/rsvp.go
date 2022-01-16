package smshandler

import (
	"context"
	"exposures/ent"
	"exposures/ent/user"
	"exposures/messages"
	"log"
	"time"
)

func Rsvp(ctx context.Context, client *ent.Client, from string, body string) []messages.Message {
	user, err := client.User.Query().Where(user.PhoneNumberEQ(from)).Only(ctx)
	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			user = client.User.Create().SetPhoneNumber(from).SaveX(ctx)

		default:
			log.Println("Huh")
		}
	}
	_, err = client.Checkin.Create().
		SetCheckinTime(time.Now().Unix()).
		SetEventID(body).
		SetSender(user).
		Save(ctx)
	if err != nil {
		log.Printf("Error creating checkin: %v", err)
	}
	return []messages.Message{
		{
			Recipient: from,
			Type:      messages.RsvpSuccess,
			Message:   "Thanks for checking in. You can let your friends/other attendees know if you tested positive by replying POSITIVE.",
		},
	}
}
