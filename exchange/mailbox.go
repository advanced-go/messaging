package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/messaging/core"
)

type Mailbox struct {
	uri  string
	ctrl chan core.Message
	data chan core.Message
}

func NewMailbox(uri string, data bool) (*Mailbox, error) {
	if len(uri) == 0 {
		return nil, errors.New(fmt.Sprintf("invalid argument: uri is empty"))
	}
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan core.Message, 16)
	if data {
		m.data = make(chan core.Message, 16)
	}
	return m, nil
}
