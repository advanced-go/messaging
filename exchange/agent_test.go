package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

func newAgentCtrlHandler(msg core.Message) {
	fmt.Printf(fmt.Sprintf("test: NewAgent_CtrlHandler() -> %v\n", msg.Event))
}

func Example_NewAgent() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("test: NewAgent() -> [recovered:%v]\n", r)
		}
	}()
	agentDir := any(NewDirectory()).(*directory)
	uri := "github.com/advanced-go/example-domain/activity"
	c := make(chan core.Message, 16)
	status := agentDir.add(newMailbox(uri, c, nil))
	if !status.OK() {
		fmt.Printf("test: add() -> [status:%v]\n", status)
	}
	a, status1 := newAgent(agentDir, uri, newAgentCtrlHandler, nil)
	if !status1.OK() {
		fmt.Printf("test: newAgent() -> [status:%v]\n", status1)
	}
	// 1 -10 Nanoseconds works for a direct send to a channel, sending via a directory needs a longer sleep time
	//d := time.Nanosecond * 10
	// Needed time.Nanoseconds * 50 for directory send with mutex
	// Needed time.Nanoseconds * 1 for directory send via sync.Map
	d := time.Nanosecond * 1
	a.Run()
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.StartupEvent})
	//c <- core.Message{To: "", From: "", Event: core.StartupEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.PauseEvent})
	//c <- core.Message{To: "", From: "", Event: core.PauseEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.ResumeEvent})
	//c <- core.Message{To: "", From: "", Event: core.ResumeEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.PingEvent})
	//c <- core.Message{To: "", From: "", Event: core.PingEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.ReconfigureEvent})
	//c <- core.Message{To: "", From: "", Event: core.ReconfigureEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.ShutdownEvent})
	//c <- core.Message{To: "", From: "", Event: core.ShutdownEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 100)

	// will panic
	c <- core.Message{}

	//Output:
	//test: NewAgent_CtrlHandler() -> event:startup
	//test: NewAgent_CtrlHandler() -> event:ping
	//test: NewAgent_CtrlHandler() -> event:reconfigure
	//test: NewAgent_CtrlHandler() -> event:shutdown
	//test: NewAgent() -> [recovered:send on closed channel]

}

func newAgentShutdownCtrlHandler(msg core.Message) {
	fmt.Printf(fmt.Sprintf("test: NewAgentShutdown_CtrlHandler() -> %v\n", msg.Event))
}

func newAgentShutdownDataHandler(msg core.Message) {
	fmt.Printf(fmt.Sprintf("test: NewAgentShutdown_DataHandler() -> %v\n", msg.RelatesTo))
}

func _Example_NewAgent_Shutdown() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("test: NewAgentShutdown() -> [recovered:%v]\n", r)
		}
	}()
	agentDir := any(NewDirectory()).(*directory)
	uri := "github.com/advanced-go/example-domain/activity"
	c := make(chan core.Message, 16)
	err := agentDir.add(newMailbox(uri, c, nil))
	if err != nil {
		fmt.Printf("test: add() -> [err:%v]\n", err)
	}
	a, err1 := newAgent(agentDir, uri, newAgentShutdownCtrlHandler, newAgentShutdownDataHandler)
	if err1 != nil {
		fmt.Printf("test: add() -> [err:%v]\n", err1)
	}
	// 1 Nanosecond is too short a wait
	d := time.Nanosecond * 10
	a.Run()
	c <- core.Message{To: "", From: "", Event: core.StartupEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	c <- core.Message{To: "", From: "", Event: core.PauseEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	c <- core.Message{To: "", From: "", Event: core.ResumeEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	c <- core.Message{To: "", From: "", Event: core.PingEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	c <- core.Message{To: "", From: "", Event: core.ReconfigureEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	c <- core.Message{To: "", From: "", Event: core.ShutdownEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 100)

	// will panic
	c <- core.Message{}

	//Output:
	//test: NewAgentShutdown_CtrlHandler() -> event:startup
	//test: NewAgentShutdown_CtrlHandler() -> event:ping
	//test: NewAgentShutdown_CtrlHandler() -> event:reconfigure
	//test: NewAgentShutdown_CtrlHandler() -> event:shutdown
	//test: NewAgentShutdown() -> [recovered:send on closed channel]

}
