package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

var agentDir = any(NewDirectory()).(*directory)

func agentMessageHandler(msg core.Message) {
	s := "test: NewCmdAgent() -> %v\n"
	fmt.Printf(fmt.Sprintf(s, msg.Event))
	/*switch msg.Event {
	case core.StartupEvent:
	case core.ShutdownEvent:
	case core.PauseEvent:
	case core.ResumeEvent:
	case core.PingEvent:
		fmt.Printf("test: agentMessageHandler() -> PauseEvent\n")
		fmt.Printf("test: agentMessageHandler() -> ResumeEvent\n")
		fmt.Printf("test: agentMessageHandler() -> PingEvent\n")
		fmt.Printf("test: agentMessageHandler() -> ReconfigureEvent\n")
	}
	*/
}

func Example_NewCmdAgent() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("test: NewCmdAgent() -> [recovered:%v]\n", r)
		}
	}()
	uri := "github.com/advanced-go/example-domain/activity"
	c := make(chan core.Message, 16)
	err := add(agentDir, newMailbox(uri, c, nil))
	if err != nil {
		fmt.Printf("test: add() -> [err:%v]\n", err)
	}
	a, err1 := newCmdAgent(agentDir, uri, agentMessageHandler)
	if err1 != nil {
		fmt.Printf("test: add() -> [err:%v]\n", err1)
	}

	a.Run()
	c <- core.Message{To: "", From: "", Event: core.StartupEvent, RelatesTo: nil, Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 500)
	c <- core.Message{To: "", From: "", Event: core.PauseEvent, RelatesTo: nil, Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 500)
	c <- core.Message{To: "", From: "", Event: core.ResumeEvent, RelatesTo: nil, Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 500)
	c <- core.Message{To: "", From: "", Event: core.PingEvent, RelatesTo: nil, Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 500)
	c <- core.Message{To: "", From: "", Event: core.ReconfigureEvent, RelatesTo: nil, Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 500)
	c <- core.Message{To: "", From: "", Event: core.ShutdownEvent, RelatesTo: nil, Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 500)

	// will panic
	c <- core.Message{}

	//Output:
	//test: NewCmdAgent() -> event:startup
	//test: NewCmdAgent() -> event:pause
	//test: NewCmdAgent() -> event:resume
	//test: NewCmdAgent() -> event:ping
	//test: NewCmdAgent() -> event:reconfigure
	//test: NewCmdAgent() -> event:shutdown
	//test: NewCmdAgent() -> [recovered:send on closed channel]

}
