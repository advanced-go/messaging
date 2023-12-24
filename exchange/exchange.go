package exchange

const (
	sendLoc     = PkgPath + ":Send"
	registerLoc = PkgPath + ":Register"
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



// SendCtrl - send to command channel
func SendCtrl(msg core.Message) runtime.Status {
	status := Root.SendCtrl(msg)
	if !status.OK() {
		status.AddLocation(sendLoc)
	}
	return status
}

// SendData - send to data channel
func SendData(msg core.Message) runtime.Status {
	status := Root.SendData(msg)
	if !status.OK() {
		status.AddLocation(sendLoc)
	}
	return status
}

// Shutdown - send a shutdown message to all directory entries
func shutdown() {
	//Root.Shutdown()
}


*/
