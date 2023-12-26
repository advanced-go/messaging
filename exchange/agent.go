package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
)

const (
	newAgentLocation      = PkgPath + ":NewAgent"
	agentSendCtrlLocation = PkgPath + ":Agent/SendCtrl"
	agentSendDataLocation = PkgPath + ":Agent/SendData"
)

// RunFunc - type for an Agent run function
type RunFunc func(m *Mailbox, activity core.MessageHandler, handlers ...core.MessageHandler)

// Agent - interface for an AI Agent
type Agent interface {
	SendCtrl(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
	Register(dir Directory) runtime.Status
	Run(activity core.MessageHandler, handlers ...core.MessageHandler)
	Shutdown()
}

type agentCfg struct {
	m       *Mailbox
	runFunc RunFunc
}

// NewDefaultAgent - create an agent with only a control channel, registered with the HostDirectory,
// and using the default run function.
func NewDefaultAgent(uri string) (Agent, runtime.Status) {
	a, status := NewAgent(uri, DefaultRun, nil)
	if status.OK() {
		status = a.Register(HostDirectory)
	}
	return a, status
}

// NewAgent - create a new agent
func NewAgent(uri string, runFunc RunFunc, data chan core.Message) (Agent, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("URI is empty"))
	}
	a := new(agentCfg)
	a.m = NewMailbox(uri, data)
	a.runFunc = runFunc
	return a, runtime.StatusOK()
}

// Run - run the agent
func (a *agentCfg) Run(activity core.MessageHandler, handlers ...core.MessageHandler) {
	if activity == nil {
		activity = func(msg core.Message) {}
	}
	//if ctrlHandler == nil {
	//	ctrlHandler = func(msg core.Message) {}
	//}
	go a.runFunc(a.m, activity, handlers...)
}

// Shutdown - shutdown the agent's mailbox
func (a *agentCfg) Shutdown() {
	if a.m.ctrl != nil {
		a.m.ctrl <- core.Message{Event: core.ShutdownEvent}
	}
}

// SendCtrl - send a message to the control channel
func (a *agentCfg) SendCtrl(msg core.Message) runtime.Status {
	if a.m.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, agentSendCtrlLocation, errors.New(fmt.Sprintf("error: control channel is nil: [%v]", a.m.uri)))
	}
	a.m.ctrl <- msg
	return runtime.StatusOK()
}

// SendData - send a message to the data channel
func (a *agentCfg) SendData(msg core.Message) runtime.Status {
	if a.m.data == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, agentSendDataLocation, errors.New(fmt.Sprintf("error: data channel is nil: [%v]", a.m.uri)))
	}
	a.m.data <- msg
	return runtime.StatusOK()
}

// Register - register an agent with a directory
func (a *agentCfg) Register(dir Directory) runtime.Status {
	//a.m.public = makePublic
	return dir.Add(a.m)
}

// DefaultRun - a simple run function that only handles control messages, and dispatches via a message handler
func DefaultRun(m *Mailbox, _ core.MessageHandler, handlers ...core.MessageHandler) {
	ctrlHandler := func(msg core.Message) {}
	if len(handlers) > 0 {
		ctrlHandler = handlers[0]
	}
	for {
		select {
		case msg, open := <-m.ctrl:
			if !open {
				return
			}
			switch msg.Event {
			case core.ShutdownEvent:
				ctrlHandler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
				m.Close()
				return
			default:
				ctrlHandler(msg)
			}
		default:
		}
	}
}
