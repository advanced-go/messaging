package core

import "fmt"

func handler(msg Message) {
	fmt.Printf(msg.Event)
}

func Example_ReplyTo() {
	msg := Message{To: "test", Event: "startup", ReplyTo: handler}
	msg.ReplyTo(msg)

	//Output:

}
