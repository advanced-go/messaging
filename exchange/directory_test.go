package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

func Example_Add() {
	uri1 := "urn:test:one"

	testDir := any(NewDirectory()).(*directory)
	m1 := newDefaultMailbox(uri1)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m0, status := testDir.get(uri1)
	fmt.Printf("test: get(%v) -> : [mbox:%v] [status:%v]\n", uri1, m0, status)

	status = testDir.Add(m1)
	fmt.Printf("test: Add(%v) -> : [status:%v]\n", uri1, status)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m0, status = testDir.get(uri1)
	fmt.Printf("test: get(%v) -> : [mbox:%v] [status:%v]\n", uri1, m0, status)

	uri2 := "urn:test:two"

	m2 := newDefaultMailbox(uri2)
	status = testDir.Add(m2)
	fmt.Printf("test: Add(%v) -> : [status:%v]\n", uri2, status)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m0, status = testDir.get(uri2)
	fmt.Printf("test: get(%v) -> : [mbox:%v] [status:%v]\n", uri2, m0, status)

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 0
	//test: get(urn:test:one) -> : [mbox:<nil>] [status:Not Found [invalid URI: directory mailbox not found [urn:test:one]]]
	//test: Add(urn:test:one) -> : [status:OK]
	//test: Count() -> : 1
	//test: get(urn:test:one) -> : [mbox:urn:test:one] [status:OK]
	//test: Add(urn:test:two) -> : [status:OK]
	//test: Count() -> : 2
	//test: get(urn:test:two) -> : [mbox:urn:test:two] [status:OK]
	//test: List() -> : [urn:test:one urn:test:two]

}

func Example_SendError() {
	uri := "urn:test"
	testDir := any(NewDirectory()).(*directory)

	fmt.Printf("test: SendCtrl(%v) -> : %v\n", uri, testDir.SendCtrl(core.Message{To: uri}))

	m := newMailbox(uri, false, nil, nil)
	status := testDir.Add(m)
	fmt.Printf("test: Add(%v) -> : [status:%v]\n", uri, status)

	//Output:
	//test: SendCtrl(urn:test) -> : Not Found [invalid URI: directory mailbox not found [urn:test]]
	//test: Add(urn:test) -> : [status:Invalid Argument [invalid argument: mailbox command channel is nil]]

}

func Example_Send() {
	uri1 := "urn:test-1"
	uri2 := "urn:test-2"
	uri3 := "urn:test-3"
	c := make(chan core.Message, 16)
	testDir := any(NewDirectory()).(*directory)

	testDir.Add(newMailbox(uri1, false, c, nil))
	testDir.Add(newMailbox(uri2, false, c, nil))
	testDir.Add(newMailbox(uri3, false, c, nil))

	testDir.SendCtrl(core.Message{To: uri1, From: PkgPath, Event: core.StartupEvent})
	testDir.SendCtrl(core.Message{To: uri2, From: PkgPath, Event: core.StartupEvent})
	testDir.SendCtrl(core.Message{To: uri3, From: PkgPath, Event: core.StartupEvent})

	time.Sleep(time.Second * 1)
	resp1 := <-c
	resp2 := <-c
	resp3 := <-c
	fmt.Printf("test: <- c -> : [%v] [%v] [%v]\n", resp1.To, resp2.To, resp3.To)
	close(c)

	//Output:
	//test: <- c -> : [urn:test-1] [urn:test-2] [urn:test-3]

}

func Example_ListCount() {
	testDir := any(NewDirectory()).(*directory)

	testDir.Add(newDefaultMailbox("test:uri1"))
	testDir.Add(newDefaultMailbox("test:uri2"))

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 2
	//test: List() -> : [test:uri1 test:uri2]

}

func Example_Remove() {
	uri := "urn:test/one"

	m := newDefaultMailbox(uri)
	testDir := any(NewDirectory()).(*directory)

	status := testDir.Add(m)
	fmt.Printf("test: Add(%v) -> : [%v]\n", uri, status)

	status = testDir.SendCtrl(core.Message{To: uri, Event: core.PingEvent})
	fmt.Printf("test: SendCtrl(%v) -> : [%v]\n", uri, status)

	m.Close()

	status = testDir.SendCtrl(core.Message{To: uri, Event: core.PingEvent})
	fmt.Printf("test: SendCtrl(%v) -> : [%v]\n", uri, status)

	//Output:
	//test: Add(urn:test/one) -> : [OK]
	//test: SendCtrl(urn:test/one) -> : [OK]
	//test: SendCtrl(urn:test/one) -> : [Not Found [invalid URI: directory mailbox not found [urn:test/one]]]

}
