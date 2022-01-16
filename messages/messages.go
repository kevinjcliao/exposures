package messages

type MessageType int

const (
	RsvpSuccess MessageType = iota
)

type message struct {
	recipient string
	message   MessageType
}

var MessageBodies = map[MessageType]string{
	RsvpSuccess: "Thanks for checking in. You can let your friends/other attendees know if you tested positive by replying POSITIVE.",
}
