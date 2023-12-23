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

type RunFunc func(m *Mailbox, activity core.MessageHandler)

type Agent interface {
	Run()
	Shutdown()
	SendCtrl(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
	Register(dir Directory, makePublic bool) runtime.Status
}

type agentCfg struct {
	m        *Mailbox
	activity core.MessageHandler
	runFunc  RunFunc
}

func NewAgent(uri string, data chan core.Message, activity core.MessageHandler, runFunc RunFunc) (Agent, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("URI is empty"))
	}
	a := new(agentCfg)
	a.m = NewMailbox(uri, data)
	a.activity = activity
	if a.activity == nil {
		a.activity = func(msg core.Message) {}
	}
	a.runFunc = runFunc
	return a, runtime.StatusOK()
}

func (a *agentCfg) Run() {
	go a.runFunc(a.m, a.activity)
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

func (a *agentCfg) Register(dir Directory, makePublic bool) runtime.Status {
	a.m.public = makePublic
	return dir.Add(a.m)
}

func DefaultRun(m *Mailbox, _ core.Message, handler core.MessageHandler) {
	for {
		select {
		case msg, open := <-m.ctrl:
			if !open {
				return
			}
			switch msg.Event {
			case core.PauseEvent:
				handler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
			case core.ResumeEvent:
				handler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
			case core.ShutdownEvent:
				handler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
				if m.shutdown != nil {
					m.shutdown()
				}
				return
			default:
				handler(msg)
			}
		default:
		}
	}
}
