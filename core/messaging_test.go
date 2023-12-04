package core

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
)

func handler(msg Message) {
	fmt.Printf(msg.Event)
}

func Example_ReplyTo() {
	msg := Message{To: "test", Event: "startup", ReplyTo: handler}
	SendReply(msg, runtime.StatusOK())

	msg = Message{To: "test", Event: "startup", ReplyTo: nil}
	SendReply(msg, runtime.StatusOK())

	//Output:
	//startup

}
