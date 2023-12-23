package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

func Example_Add() {
	uri := "urn:test"
	uri2 := "urn:test:two"

	testDir := any(NewDirectory()).(*directory)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2, _ := testDir.get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	testDir.Add(newMailbox(uri, false, nil, nil))
	fmt.Printf("test: Add(%v) -> : ok\n", uri)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2, _ = testDir.get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	testDir.Add(newMailbox(uri2, false, nil, nil))
	fmt.Printf("test: Add(%v) -> : ok\n", uri2)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2, _ = testDir.get(uri2)
	fmt.Printf("test: get(%v) -> : %v\n", uri2, d2)

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 0
	//test: get(urn:test) -> : <nil>
	//test: Add(urn:test) -> : ok
	//test: Count() -> : 1
	//test: get(urn:test) -> : &{urn:test <nil> <nil>}
	//test: Add(urn:test:two) -> : ok
	//test: Count() -> : 2
	//test: get(urn:test:two) -> : &{urn:test:two <nil> <nil>}
	//test: List() -> : [urn:test urn:test:two]

}

func Example_SendError() {
	uri := "urn:test"
	testDir := any(NewDirectory()).(*directory)

	fmt.Printf("test: SendCtrl(%v) -> : %v\n", uri, testDir.SendCtrl(core.Message{To: uri}))

	testDir.Add(newMailbox(uri, false, nil, nil))
	fmt.Printf("test: Add(%v) -> : ok\n", uri)
	fmt.Printf("test: SendCtrl(%v) -> : %v\n", uri, testDir.SendCtrl(core.Message{To: uri}))

	//Output:
	//test: SendCtrl(urn:test) -> : Not Found [invalid URI: directory mailbox not found [urn:test]]
	//test: Add(urn:test) -> : ok
	//test: SendCtrl(urn:test) -> : Invalid Content [entry control channel is nil: [urn:test]]

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

	testDir.Add(newMailbox("test:uri1", false, nil, nil))
	testDir.Add(newMailbox("test:uri2", false, nil, nil))

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 2
	//test: List() -> : [test:uri1 test:uri2]

}
