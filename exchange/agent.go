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
	m          *Mailbox
	cmdHandler core.MessageHandler
}

func NewCmdAgent(m *Mailbox, cmdHandler core.MessageHandler) (Agent, error) {
	if m == nil {
		return nil, errors.New(fmt.Sprintf("invalid argument: mailbox is nil"))
	}
	cfg := new(agentCfg)
	cfg.m = m
	cfg.cmdHandler = cmdHandler
	return cfg, nil
}

func (a *agentCfg) Run() {
	go run(a.m.cmd, a.cmdHandler)
}

func run(cmd chan core.Message, cmdHandler core.MessageHandler) {
	for {
		select {
		case msg, open := <-cmd:
			if !open {
				return
			}
			if msg.Event == core.ShutdownEvent {
				cmdHandler(msg)
				return
			}
			go cmdHandler(msg)
		default:
		}
	}
}
