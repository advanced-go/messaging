package core

import (
	"github.com/advanced-go/core/runtime"
)

const (
	StartupEvent     = "event:startup"
	ShutdownEvent    = "event:shutdown"
	PingEvent        = "event:ping"
	ReconfigureEvent = "event:reconfigure"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive

	msgHandlerLoc = "github.com/advanced-go/messaging/core:MessageCacheHandler"
)

// MessageMap - map of messages
type MessageMap map[string]Message

// MessageHandler - function type to process a Message
type MessageHandler func(msg Message)

// Message - message payload
type Message struct {
	To        string
	From      string
	Event     string
	RelatesTo string
	Status    runtime.Status
	Content   []any
	ReplyTo   MessageHandler
}

// SendReply - function used by message recipient to reply with a runtime.Status
func SendReply(msg Message, status runtime.Status) {
	if msg.ReplyTo == nil {
		return
	}
	msg.ReplyTo(Message{
		To:        msg.From,
		From:      msg.To,
		RelatesTo: msg.RelatesTo,
		Event:     "",
		Status:    status,
		Content:   nil,
		ReplyTo:   nil,
	})
}

// NewMessageCacheHandler - handler to receive messages into a cache.
func NewMessageCacheHandler[E runtime.ErrorHandler](cache MessageCache) MessageHandler {
	return func(msg Message) {
		err := cache.Add(msg)
		if err != nil {
			var e E
			status := runtime.NewStatusError(runtime.StatusInvalidArgument, msgHandlerLoc, err)
			e.Handle(status, "", "")
		}
	}
}
