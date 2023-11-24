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
		if e.c != nil {
			close(e.c)
		}
		delete(d.m, key)
	}
}

func get(uri string, d *directory) *entry {
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

	testDir.Add(uri, nil)
	fmt.Printf("test: Add(%v) -> : ok\n", uri)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	d2 = get(uri, testDir)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	testDir.Add(uri2, nil)
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
	//test: get(urn:test) -> : &{urn:test <nil>}
	//test: Add(urn:test:two) -> : ok
	//test: Count() -> : 2
	//test: get(urn:test:two) -> : &{urn:test:two <nil>}
	//test: List() -> : [urn:test urn:test:two]

}

func Example_SendError() {
	uri := "urn:test"
	empty(testDir)

	fmt.Printf("test: Send(%v) -> : %v\n", uri, testDir.Send(core.Message{To: uri}))

	testDir.Add(uri, nil)
	fmt.Printf("test: Add(%v) -> : ok\n", uri)
	fmt.Printf("test: Send(%v) -> : %v\n", uri, testDir.Send(core.Message{To: uri}))

	//Output:
	//test: Send(urn:test) -> : Invalid Argument [entry not found: [urn:test]]
	//test: Add(urn:test) -> : ok
	//test: Send(urn:test) -> : Invalid Content [entry channel is nil: [urn:test]]

}

func Example_Send() {
	uri1 := "urn:test-1"
	uri2 := "urn:test-2"
	uri3 := "urn:test-3"
	c := make(chan core.Message, 16)
	empty(testDir)

	testDir.Add(uri1, c)
	testDir.Add(uri2, c)
	testDir.Add(uri3, c)

	testDir.Send(core.Message{To: uri1, From: PkgPath, Event: core.StartupEvent})
	testDir.Send(core.Message{To: uri2, From: PkgPath, Event: core.StartupEvent})
	testDir.Send(core.Message{To: uri3, From: PkgPath, Event: core.StartupEvent})

	time.Sleep(time.Second * 1)
	resp1 := <-c
	resp2 := <-c
	resp3 := <-c
	fmt.Printf("test: <- c -> : [%v] [%v] [%v]\n", resp1.To, resp2.To, resp3.To)
	close(c)

	//Output:
	//test: <- c -> : [urn:test-1] [urn:test-2] [urn:test-3]

}
