package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/messaging/core"
)

var exchDir = NewDirectory()

// Register - add an entry to the exchange directory
func Register(uri string, c chan core.Message) error {
	if uri == "" {
		return errors.New("invalid argument: uri is empty")
	}
	if c == nil {
		return errors.New(fmt.Sprintf("invalid argument: channel is nil for [%v]", uri))
	}
	exchDir.Add(uri, c)
	return nil
}

// Send - send a message
func Send(msg core.Message) {
	exchDir.Send(msg)
}

// Shutdown - send a shutdown message to all directory entries
func Shutdown() {
	exchDir.Shutdown()
}
