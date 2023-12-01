package exchange

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"time"
)

var testDir = any(NewDirectory()).(*directory)

func empty(d *directory) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for key, e := range d.m {
		if e.cmd != nil {
			close(e.cmd)
		}
		if e.data != nil {
			close(e.data)
		}
		delete(d.m, key)
	}
}

func get(uri string, d *directory) *Mailbox {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.m[uri]
}

func Example_Add() {
	uri := "urn:test"
	uri2 := "urn:test:two"

	empty(testDir)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2 := get(uri, testDir)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	testDir.add(newMailbox(uri, nil, nil))
	fmt.Printf("test: Add(%v) -> : ok\n", uri)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2 = get(uri, testDir)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	testDir.add(newMailbox(uri2, nil, nil))
	fmt.Printf("test: Add(%v) -> : ok\n", uri2)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2 = get(uri2, testDir)
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
	empty(testDir)

	fmt.Printf("test: SendCmd(%v) -> : %v\n", uri, testDir.SendCmd(core.Message{To: uri}))

	testDir.add(newMailbox(uri, nil, nil))
	fmt.Printf("test: Add(%v) -> : ok\n", uri)
	fmt.Printf("test: SendCmd(%v) -> : %v\n", uri, testDir.SendCmd(core.Message{To: uri}))

	//Output:
	//test: SendCmd(urn:test) -> : Invalid Argument [entry not found: [urn:test]]
	//test: Add(urn:test) -> : ok
	//test: SendCmd(urn:test) -> : Invalid Content [entry command channel is nil: [urn:test]]

}

func Example_Send() {
	uri1 := "urn:test-1"
	uri2 := "urn:test-2"
	uri3 := "urn:test-3"
	c := make(chan core.Message, 16)
	empty(testDir)

	testDir.add(newMailbox(uri1, c, nil))
	testDir.add(newMailbox(uri2, c, nil))
	testDir.add(newMailbox(uri3, c, nil))

	testDir.SendCmd(core.Message{To: uri1, From: PkgPath, Event: core.StartupEvent})
	testDir.SendCmd(core.Message{To: uri2, From: PkgPath, Event: core.StartupEvent})
	testDir.SendCmd(core.Message{To: uri3, From: PkgPath, Event: core.StartupEvent})

	time.Sleep(time.Second * 1)
	resp1 := <-c
	resp2 := <-c
	resp3 := <-c
	fmt.Printf("test: <- c -> : [%v] [%v] [%v]\n", resp1.To, resp2.To, resp3.To)
	close(c)

	//Output:
	//test: <- c -> : [urn:test-1] [urn:test-2] [urn:test-3]

}
