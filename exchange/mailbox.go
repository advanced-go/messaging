package exchange

import (
	"errors"
	"github.com/advanced-go/messaging/core"
)

type Mailbox struct {
	uri  string
	cmd  chan core.Message
	data chan core.Message
}

func NewMailbox(uri string, data bool) (*Mailbox, error) {
	if len(uri) == 0 {
		return nil, errors.New("invalid argument: uri is empty")
	}
	m := new(Mailbox)
	m.uri = uri
	m.cmd = make(chan core.Message, 16)
	if data {
		m.data = make(chan core.Message, 16)
	}
	return m, nil
}

func newMailbox(uri string, cmd, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.cmd = cmd
	m.data = data
	return m
}
