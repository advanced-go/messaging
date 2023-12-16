package exchange

import "github.com/advanced-go/messaging/core"

func newMailbox(uri string, public bool, ctrl, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.public = public
	m.uri = uri
	m.ctrl = ctrl
	m.data = data
	return m
}
