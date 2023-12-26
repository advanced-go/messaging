package exchange

import (
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
)

const (
	newAgentLocation = PkgPath + ":NewDefaultAgent"
)

// Agent - interface for an AI Agent
type Agent interface {
	SendCtrl(msg core.Message)
	SendData(msg core.Message)
	Register(dir Directory) runtime.Status
	Run()
	Shutdown()
}

type agentCfg struct {
	m   *Mailbox
	run func(m *Mailbox)
}

// NewDefaultAgent - create an agent with only a control channel, registered with the HostDirectory,
// and using the default run function.
func NewDefaultAgent(uri string, ctrlHandler core.MessageHandler) (Agent, runtime.Status) {
	return newDefaultAgent(uri, ctrlHandler, HostDirectory)
}

func newDefaultAgent(uri string, ctrlHandler core.MessageHandler, dir Directory) (Agent, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("URI is empty"))
	}
	if ctrlHandler == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("controller message handler is nil"))
	}
	a := new(agentCfg)
	a.m = NewMailbox(uri, nil)
	a.run = func(m *Mailbox) {
		DefaultRun(m, ctrlHandler)
	}
	return a, a.Register(dir)
}

// Run - run the agent
func (a *agentCfg) Run() {
	go a.run(a.m)
}

// Shutdown - shutdown the agent
func (a *agentCfg) Shutdown() {
	a.m.SendCtrl(core.Message{Event: core.ShutdownEvent})
}

// SendCtrl - send a message to the control channel
func (a *agentCfg) SendCtrl(msg core.Message) {
	a.m.SendCtrl(msg)
}

// SendData - send a message to the data channel
func (a *agentCfg) SendData(msg core.Message) {
	a.m.SendData(msg)
}

// Register - register an agent with a directory
func (a *agentCfg) Register(dir Directory) runtime.Status {
	return dir.Add(a.m)
}

// DefaultRun - a simple run function that only handles control messages, and dispatches via a message handler
func DefaultRun(m *Mailbox, ctrlHandler core.MessageHandler) {
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
