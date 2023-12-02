package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

var agentDir = any(NewDirectory()).(*directory)

func agentMessageHandler(msg core.Message) {
	fmt.Printf(fmt.Sprintf("test: NewAgent() -> %v\n", msg.Event))
}

func Example_NewAgent() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("test: NewAgent() -> [recovered:%v]\n", r)
		}
	}()
	uri := "github.com/advanced-go/example-domain/activity"
	c := make(chan core.Message, 16)
	err := add(agentDir, newMailbox(uri, c, nil))
	if err != nil {
		fmt.Printf("test: add() -> [err:%v]\n", err)
	}
	a, err1 := newAgent(agentDir, uri, agentMessageHandler, nil)
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
	//test: NewAgent() -> event:startup
	//test: NewAgent() -> event:ping
	//test: NewAgent() -> event:reconfigure
	//test: NewAgent() -> event:shutdown
	//test: NewAgent() -> [recovered:send on closed channel]

}
