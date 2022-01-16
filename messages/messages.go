package messages

type MessageType int

const (
	RsvpSuccess MessageType = iota
	NotifyPositiveCase
	ThankForSelfReporting
)

type Message struct {
	Recipient string
	Type      MessageType
	Message   string
}
