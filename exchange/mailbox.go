package exchange

import (
	"github.com/advanced-go/messaging/core"
)

type Mailbox struct {
	uri  string
	cmd  chan core.Message
	data chan core.Message
}

func NewMailbox(uri string, data bool) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.cmd = make(chan core.Message, 16)
	if data {
		m.data = make(chan core.Message, 16)
	}
	return m
}

func newMailbox(uri string, cmd, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.cmd = cmd
	m.data = data
	return m
}
