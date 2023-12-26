package exchange

import (
	"github.com/advanced-go/messaging/core"
)

const (
	mailboxLoc = PkgPath + ":NewMailbox"
)

type Mailbox struct {
	public     bool
	uri        string
	ctrl       chan core.Message
	data       chan core.Message
	unregister func()
}

// NewMailbox - create a new mailbox
func NewMailbox(uri string, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan core.Message, 16)
	m.data = data
	return m
}

// Uri - return the mailbox uri
func (m *Mailbox) Uri() string {
	return m.uri
}

// String - return the mailbox uri
func (m *Mailbox) String() string {
	return m.uri
}

// SendCtrl - send a message to the control channel
func (m *Mailbox) SendCtrl(msg core.Message) {
	if m.ctrl != nil {
		m.ctrl <- msg
	}
}

// SendData - send a message to the data channel
func (m *Mailbox) SendData(msg core.Message) {
	if m.data != nil {
		m.data <- msg
	}
}

// Close - close the mailbox channels and unregsiter the mailbox with a Directory
func (m *Mailbox) Close() {
	if m.unregister != nil {
		m.unregister()
	}
	if m.data != nil {
		close(m.data)
	}
	if m.ctrl != nil {
		close(m.ctrl)
	}
}
