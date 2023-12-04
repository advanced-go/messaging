package core

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
)

const (
	StartupEvent     = "event:startup"
	ShutdownEvent    = "event:shutdown"
	PingEvent        = "event:ping"
	ReconfigureEvent = "event:reconfigure"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive
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
		Event:     msg.Event,
		Status:    status,
		Content:   nil,
		ReplyTo:   nil,
	})
}

// NewMessageCacheHandler - handler to receive messages into a cache.
func NewMessageCacheHandler(cache MessageCache) MessageHandler {
	return func(msg Message) {
		err := cache.Add(msg)
		if err != nil {
			fmt.Printf("error on MessageCache.Add() -> [to:%v] [from:%v] [err:%v]\n", msg.To, msg.From, err)
		}
	}
}
