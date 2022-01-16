package smshandler

import (
	"exposures/messages"
	"exposures/test"
	"testing"
)

func TestHandlePositiveCase(t *testing.T) {
	client, ctx := test.SetupTestWithEntClient(t)
	defer client.Close()

	user1 := "+1(123)456-7890"
	user2 := "+1(234)567-8901"
	user3 := "+1(345)678-9012"
	sampleCode := "205202ac-c4e6-48ed-b469-9b3bcf592316"
	Rsvp(ctx, client, user1, sampleCode)
	Rsvp(ctx, client, user2, sampleCode)
	Rsvp(ctx, client, user3, sampleCode)

	ms := HandlePositiveCase(ctx, client, user1)
	if len(ms) != 3 {
		t.Fatalf("Expected 3 messages to be sent. Got: %v", len(ms))
	}

	for _, m := range ms {
		if m.Recipient == user1 && m.Type != messages.ThankForSelfReporting {
			t.Fatalf("Expected to thank user 1 for reporting. Found this message instead: %v", ms)
		}

		if m.Recipient == user2 && m.Type != messages.NotifyPositiveCase {
			t.Fatalf("Expected to notify user 2. Sent this message instead: %v", ms)
		}

		if m.Recipient == user3 && m.Type != messages.NotifyPositiveCase {
			t.Fatalf("Expected to notify user 3. Sent this message instead: %v", ms)
		}
	}

	users := client.User.Query().AllX(ctx)
	if len(users) != 3 {
		t.Errorf("Expected 3 users to be created and persisted. Found: %v", len(users))
	}
}

func TestOnlyNotifyForCorrectCases(t *testing.T) {
	client, ctx := test.SetupTestWithEntClient(t)
	defer client.Close()

	user1 := "+1(123)456-7890"
	user2 := "+1(234)567-8901"
	user3 := "+1(345)678-9012"
	sampleCode := "205202ac-c4e6-48ed-b469-9b3bcf592316"
	sampleCode2 := "e3fe8a5d-97ec-402b-99a0-04c04822aad2"
	Rsvp(ctx, client, user1, sampleCode)
	Rsvp(ctx, client, user2, sampleCode)
	Rsvp(ctx, client, user3, sampleCode2)

	ms := HandlePositiveCase(ctx, client, user1)

	for _, m := range ms {
		if m.Recipient == user3 {
			t.Fatalf("User 3 should not be exposed.")
		}

	}

	ms = HandlePositiveCase(ctx, client, user2)
	for _, m := range ms {
		if m.Recipient == user1 && m.Type != messages.NotifyPositiveCase {
			t.Fatalf("User 1 should get a positive case notification.")
		}

		if m.Recipient == user2 && m.Type != messages.ThankForSelfReporting {
			t.Fatalf("User2 should be thanked for self-reporting.")
		}
	}
}
