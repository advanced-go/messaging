package exchange

import (
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
)

const (
	mailboxLoc = PkgPath + ":NewMailbox"
)

type Mailbox2 struct {
	dir    Directory
	public bool
	uri    string
	ctrl   chan core.Message
	data   chan core.Message
}

type Mailbox struct {
	public   bool
	uri      string
	ctrl     chan core.Message
	data     chan core.Message
	shutdown func() runtime.Status
}

func NewMailbox(uri string, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan core.Message, 16)
	m.data = data
	return m
}
