package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/messaging/core"
)

type Agent interface {
	Run()
}

type agentCfg struct {
	m           *Mailbox
	ctrlHandler core.MessageHandler
	dataHandler core.MessageHandler
}

func NewAgent(uri string, ctrlHandler, dataHandler core.MessageHandler) (Agent, error) {
	return newAgent(exchDir, uri, ctrlHandler, dataHandler)
}

func newAgent(dir Directory, uri string, ctrlHandler, dataHandler core.MessageHandler) (Agent, error) {
	if len(uri) == 0 {
		return nil, errors.New(fmt.Sprintf("invalid argument: uri is empty"))
	}
	m, err := get(dir, uri)
	if err != nil {
		return nil, err
	}
	cfg := new(agentCfg)
	cfg.m = m
	cfg.ctrlHandler = ctrlHandler
	cfg.dataHandler = dataHandler
	return cfg, nil
}

func (a *agentCfg) Run() {
	go run(a)
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
