package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
)

func Example_NewMailbox() {
	m := NewMailbox("github.com/advanced-go/messaging", nil)
	fmt.Printf("test: NewMailbox() -> %v", m)

	//Output:
	//test: NewMailbox() -> github.com/advanced-go/messaging

}

func newMailbox(uri string, public bool, ctrl, data chan core.Message) *Mailbox {
	m := new(Mailbox)
	m.public = public
	m.uri = uri
	m.ctrl = ctrl
	m.data = data
	return m
}

func newDefaultMailbox(uri string) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan core.Message, 16)
	return m
}
