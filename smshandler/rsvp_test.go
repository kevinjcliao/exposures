package smshandler

import (
	"exposures/messages"
	"exposures/test"
	"testing"
)

func TestRsvp(t *testing.T) {
	client, ctx := test.SetupTestWithEntClient(t)
	defer client.Close()

	sampleCode := "205202ac-c4e6-48ed-b469-9b3bcf592316"
	ms := Rsvp(ctx, client, "+1(123)456-7890", sampleCode)
	if len(ms) != 1 {
		t.Errorf(
			"Expected RSVP to send 1 message. Sent %v.", ms,
		)
	}

	if ms[0].Type != messages.RsvpSuccess {
		t.Errorf(
			"Expected RSVP message, but actually got: %v", ms,
		)
	}

	users := client.User.Query().AllX(ctx)

	if len(users) != 1 {
		t.Errorf(
			"Expected 1 user to be created and persisted to DB, but actually got: %v", len(users),
		)
	}
}
