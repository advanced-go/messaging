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

type RunFunc func(m *Mailbox, activity core.MessageHandler, ctrlHandler core.MessageHandler)

type Agent interface {
	Run(activity core.MessageHandler, ctrlHandler core.MessageHandler)
	Shutdown()
	SendCtrl(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
	Register(dir Directory) runtime.Status
}

type agentCfg struct {
	m       *Mailbox
	runFunc RunFunc
}

func NewDefaultAgent(uri string) (Agent, runtime.Status) {
	return NewAgent(uri, DefaultRun, nil)
}

func NewAgent(uri string, runFunc RunFunc, data chan core.Message) (Agent, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("URI is empty"))
	}
	a := new(agentCfg)
	a.m = NewMailbox(uri, data)
	a.runFunc = runFunc
	return a, runtime.StatusOK()
}

func (a *agentCfg) Run(activity core.MessageHandler, ctrlHandler core.MessageHandler) {
	if activity == nil {
		activity = func(msg core.Message) {}
	}
	if ctrlHandler == nil {
		ctrlHandler = func(msg core.Message) {}
	}
	go a.runFunc(a.m, activity, ctrlHandler)
}

func (a *agentCfg) Shutdown() {
	if a.m.ctrl != nil {
		a.m.ctrl <- core.Message{Event: core.ShutdownEvent}
	}
}

func (a *agentCfg) SendCtrl(msg core.Message) runtime.Status {
	if a.m.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, agentSendCtrlLocation, errors.New(fmt.Sprintf("error: control channel is nil: [%v]", a.m.uri)))
	}
	a.m.ctrl <- msg
	return runtime.StatusOK()
}

func (a *agentCfg) SendData(msg core.Message) runtime.Status {
	if a.m.data == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, agentSendDataLocation, errors.New(fmt.Sprintf("error: data channel is nil: [%v]", a.m.uri)))
	}
	a.m.data <- msg
	return runtime.StatusOK()
}

func (a *agentCfg) Register(dir Directory) runtime.Status {
	//a.m.public = makePublic
	return dir.Add(a.m)
}

func DefaultRun(m *Mailbox, _ core.MessageHandler, ctrlHandler core.MessageHandler) {
	for {
		select {
		case msg, open := <-m.ctrl:
			if !open {
				return
			}
			switch msg.Event {
			case core.ShutdownEvent:
				ctrlHandler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
				if m.shutdown != nil {
					status := m.shutdown()
					if !status.OK() {
						fmt.Println(status)
					}
				}
				return
			default:
				ctrlHandler(msg)
			}
		default:
		}
	}
}
