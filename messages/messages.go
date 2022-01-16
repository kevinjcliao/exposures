package messages

type MessageType int

const (
	RsvpSuccess MessageType = iota
	NotifyPositiveCase
	ThankForSelfReporting
	Error
)

type Message struct {
	Recipient string
	Type      MessageType
	Message   string
}
