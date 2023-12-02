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
	cmd        chan core.Message
	cmdHandler core.MessageHandler
}

func NewCmdAgent(uri string, cmdHandler core.MessageHandler) (Agent, error) {
	return newCmdAgent(exchDir, uri, cmdHandler)
}

func newCmdAgent(dir Directory, uri string, cmdHandler core.MessageHandler) (Agent, error) {
	if len(uri) == 0 {
		return nil, errors.New(fmt.Sprintf("invalid argument: uri is empty"))
	}
	m, err := get(dir, uri)
	if err != nil {
		return nil, err
	}
	cfg := new(agentCfg)
	cfg.cmd = m.cmd
	cfg.cmdHandler = cmdHandler
	return cfg, nil
}

func (a *agentCfg) Run() {
	go run(a.cmd, a.cmdHandler)
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
				close(cmd)
				return
			}
			go cmdHandler(msg)
		default:
		}
	}
}
