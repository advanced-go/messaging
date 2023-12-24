package exchange

import (
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
)

const (
	sendCtrlLoc = PkgPath + ":SendCtrl"
	sendDataLoc = PkgPath + ":SendData"

	//registerLoc = PkgPath + ":Register"
)

var HostDirectory = NewDirectory()

/*
// Register - add a mailbox to the exchange directory

func Register(m *Mailbox) runtime.Status {
	if m == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, registerLoc, errors.New("invalid argument: mailbox is nil"))
	}
	if len(m.uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, registerLoc, errors.New("invalid argument: mailbox uri is empty"))
	}
	if m.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, registerLoc, errors.New("invalid argument: mailbox command channel is nil"))
	}
	d, ok := any(Root).(*directory)
	if !ok {
		return runtime.NewStatusError(runtime.StatusInvalidContent, registerLoc, errors.New("invalid argument: Directory type is not *directory"))
	}
	//d.add(m)
	return runtime.StatusOK()
}



*/

// SendCtrl - send to command channel
func SendCtrl(msg core.Message) runtime.Status {
	status := HostDirectory.SendCtrl(msg)
	if !status.OK() {
		status.AddLocation(sendCtrlLoc)
	}
	return status
}

// SendData - send to data channel
func SendData(msg core.Message) runtime.Status {
	status := HostDirectory.SendData(msg)
	if !status.OK() {
		status.AddLocation(sendDataLoc)
	}
	return status
}

// Shutdown - send a shutdown message to all directory entries
//func shutdown() {
//HostDirectory.Shutdown()
//}

func shutdownHost(msg core.Message) runtime.Status {
	//TO DO: authentication and implementation

	return runtime.StatusOK()
}
