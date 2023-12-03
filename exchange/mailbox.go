package exchange

import (
	"github.com/advanced-go/messaging/core"
)

const (
	mailboxLoc = PkgPath + ":NewMailbox"
)

type Mailbox struct {
	uri  string
	ctrl chan core.Message
	data chan core.Message
}

func NewMailbox(uri string, data bool) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan core.Message, 16)
	if data {
		m.data = make(chan core.Message, 16)
	}
	return m
}
