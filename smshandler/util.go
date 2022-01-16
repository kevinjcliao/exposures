package smshandler

import (
	"context"
	"exposures/ent"
	"exposures/messages"
)

type SmsHandlerFunc func(ctx context.Context, client *ent.Client) messages.Message
