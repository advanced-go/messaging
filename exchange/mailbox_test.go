package exchange

import "github.com/advanced-go/messaging/core"

func newMailbox(uri string, cmd, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.cmd = cmd
	m.data = data
	return m
}
