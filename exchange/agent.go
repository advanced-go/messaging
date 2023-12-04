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
	dir           *directory
	m             *Mailbox
	ctrlHandler   core.MessageHandler
	dataHandler   core.MessageHandler
	statusHandler core.MessageHandler
}

func NewAgent(uri string, ctrl, data, status core.MessageHandler) (Agent, runtime.Status) {
	dir := any(exchDir).(*directory)
	if dir == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, newAgentLocation, errors.New(fmt.Sprintf("Directory is not of *directory type")))
	}
	return newAgent(dir, uri, ctrl, data, status)
}

func newAgent(dir *directory, uri string, ctrl, data, status core.MessageHandler) (Agent, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, newAgentLocation, errors.New("invalid argument: uri is empty"))
	}
	m, status1 := dir.get(uri)
	if !status1.OK() {
		return nil, status1
	}
	cfg := new(agentCfg)
	cfg.dir = dir
	cfg.m = m
	cfg.ctrlHandler = ctrl
	cfg.dataHandler = data
	cfg.statusHandler = status
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
				if cfg.statusHandler != nil {
					go cfg.statusHandler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
				}
			case core.ResumeEvent:
				paused = false
				if cfg.statusHandler != nil {
					go cfg.statusHandler(core.Message{Event: msg.Event, Status: runtime.StatusOK()})
				}
			case core.ShutdownEvent:
				cfg.ctrlHandler(msg)
				status := cfg.dir.shutdown(cfg.m.uri)
				if cfg.statusHandler != nil {
					go cfg.statusHandler(core.Message{Event: msg.Event, Status: status})
				}
				return
			default:
				go cfg.ctrlHandler(msg)
			}
		default:
		}
	}
}
