package smshandler

import (
	"context"
	"exposures/ent"
	"exposures/ent/checkin"
	"exposures/ent/user"
	"exposures/messages"
	"fmt"
	"time"
)

func HandlePositiveCase(ctx context.Context, client *ent.Client, from string) []messages.Message {
	results := []messages.Message{
		{Recipient: from, Type: messages.ThankForSelfReporting, Message: "Thanks for letting us know. Please take care of yourself and those around you. We hope you recover soon."},
	}
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
		for phone := range phoneNumbers {
			loc, _ := time.LoadLocation("America/Los_Angeles")
			date := time.Unix(positiveCheckin.CheckinTime, 0).In(loc).Format("2002 Jan 06 15:04 PST")
			body := fmt.Sprintf("Someone at an event you attended on: %s reported testing positive for COVID-19. You should get tested.", date)

			results = append(results, messages.Message{
				Recipient: phone,
				Message:   body,
				Type:      messages.NotifyPositiveCase,
			})
		}
	}
	return results
}
