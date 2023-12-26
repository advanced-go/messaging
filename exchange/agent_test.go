package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

func newAgentCtrlHandler(msg core.Message) {
	fmt.Printf(fmt.Sprintf("test: NewDefaultAgent_CtrlHandler() -> %v\n", msg.Event))
}

func Example_NewDefaultAgent() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("test: NewDefaultAgent() -> [recovered:%v]\n", r)
		}
	}()
	agentDir := NewDirectory() //any(NewDirectory()).(*directory)
	uri := "github.com/advanced-go/example-domain/activity"
	//c := make(chan core.Message, 16)
	a, status := newDefaultAgent(uri, newAgentCtrlHandler, agentDir)
	if !status.OK() {
		fmt.Printf("test: NewDefaultAgent() -> [status:%v]\n", status)
	}
	//status = a.Register(agentDir)
	//if !status.OK() {
	//	fmt.Printf("test: Register() -> [status:%v]\n", status)
	//}
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
	a.SendCtrl(core.Message{To: uri, From: "", Event: core.ResumeEvent})
	//c <- core.Message{To: "", From: "", Event: core.ResumeEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	agentDir.SendCtrl(core.Message{To: uri, From: "", Event: core.PingEvent})
	//c <- core.Message{To: "", From: "", Event: core.PingEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.SendCtrl(core.Message{To: uri, From: "", Event: core.ReconfigureEvent})
	//c <- core.Message{To: "", From: "", Event: core.ReconfigureEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.Shutdown() //.SendCtrl(core.Message{To: uri, From: "", Event: core.ShutdownEvent})
	//c <- core.Message{To: "", From: "", Event: core.ShutdownEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 100)

	// will panic
	//c <- core.Message{}

	//Output:
	//test: NewDefaultAgent_CtrlHandler() -> event:startup
	//test: NewDefaultAgent_CtrlHandler() -> event:pause
	//test: NewDefaultAgent_CtrlHandler() -> event:resume
	//test: NewDefaultAgent_CtrlHandler() -> event:ping
	//test: NewDefaultAgent_CtrlHandler() -> event:reconfigure
	//test: NewDefaultAgent_CtrlHandler() -> event:shutdown

}
