package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
)

const (
	newAgentLocation = PkgPath + ":NewAgent"
)

type Agent interface {
	Run()
	Shutdown()
}

type agentCfg struct {
	m           *Mailbox
	ctrlHandler core.MessageHandler
	dataHandler core.MessageHandler
}

func NewAgent(uri string, ctrlHandler, dataHandler core.MessageHandler) (Agent, runtime.Status) {
	dir := any(exchDir).(*directory)
	if dir == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, newAgentLocation, errors.New(fmt.Sprintf("Directory is not of *directory type")))
	}
	return newAgent(dir, uri, ctrlHandler, dataHandler)
}

func newAgent(dir *directory, uri string, ctrlHandler, dataHandler core.MessageHandler) (Agent, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("invalid argument: uri is empty"))
	}
	m, status := dir.get(uri)
	if !status.OK() {
		return nil, status
	}
	cfg := new(agentCfg)
	cfg.m = m
	cfg.ctrlHandler = ctrlHandler
	cfg.dataHandler = dataHandler
	return cfg, runtime.StatusOK()
}

func (a *agentCfg) Run() {
	go run(a)
}

func (a *agentCfg) Shutdown() {
	if a.m.ctrl != nil {
		a.m.ctrl <- core.Message{Event: core.ShutdownEvent}
	}
}

func run(cfg *agentCfg) {
	paused := false
	for {
		if cfg.m.data != nil && !paused {
			select {
			case msg, open := <-cfg.m.data:
				if open {
					go cfg.dataHandler(msg)
				}
			default:
			}
		}
		select {
		case msg, open := <-cfg.m.ctrl:
			if !open {
				return
			}
			switch msg.Event {
			case core.PauseEvent:
				paused = true
			case core.ResumeEvent:
				paused = false
			case core.ShutdownEvent:
				cfg.ctrlHandler(msg)
				close(cfg.m.ctrl)
				if cfg.m.data != nil {
					close(cfg.m.data)
				}
				return
			default:
				go cfg.ctrlHandler(msg)
			}
		default:
		}
	}
}
