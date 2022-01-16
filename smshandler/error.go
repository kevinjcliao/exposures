package smshandler

import (
	"context"
	"exposures/ent"
	"exposures/messages"
)

func Error(ctx context.Context, client *ent.Client, from string, body string) []messages.Message {
	return []messages.Message{
		{
			Recipient: from,
			Type:      messages.Error,
			Message:   "Sorry, I don't understand this message. Please scan a QR code to check in, or send me POSITIVE if you've tested positive.",
		},
	}
}
