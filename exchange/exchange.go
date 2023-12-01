package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
)

const (
	sendLoc = PkgPath + ":Send"
)

var exchDir = NewDirectory()

// Register - add an entry to the exchange directory
func Register(m *Mailbox) error {
	if m == nil {
		return errors.New(fmt.Sprintf("invalid argument: mailbox is nil"))
	}
	exchDir.add(m)
	return nil
}

// SendCmd - send to command channel
func SendCmd(msg core.Message) runtime.Status {
	status := exchDir.SendCmd(msg)
	if !status.OK() {
		status.AddLocation(sendLoc)
	}
	return status
}

// SendData - send to data channel
func SendData(msg core.Message) runtime.Status {
	status := exchDir.SendData(msg)
	if !status.OK() {
		status.AddLocation(sendLoc)
	}
	return status
}

// Shutdown - send a shutdown message to all directory entries
func Shutdown() {
	exchDir.Shutdown()
}
