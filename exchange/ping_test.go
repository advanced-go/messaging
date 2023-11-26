package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"time"
)

var pingDir = any(NewDirectory()).(*directory)

func ExamplePing() {
	uri1 := "urn:ping:good"
	uri2 := "urn:ping:bad"
	uri3 := "urn:ping:error"
	uri4 := "urn:ping:delay"

	start = time.Now()
	empty(pingDir)

	c := make(chan core.Message, 16)
	pingDir.Add(uri1, c)
	go pingGood(c)
	status := ping[runtime.TestError](pingDir, nil, uri1)
	duration := status.Duration()
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan core.Message, 16)
	pingDir.Add(uri2, c)
	go pingBad(c)
	status = ping[runtime.TestError](pingDir, nil, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan core.Message, 16)
	pingDir.Add(uri3, c)
	go pingError(c, errors.New("ping depends error message"))
	status = ping[runtime.TestError](pingDir, nil, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan core.Message, 16)
	pingDir.Add(uri4, c)
	go pingDelay(c)
	status = ping[runtime.TestError](pingDir, nil, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status, duration)

	//Output:
	//test: Ping(good) -> [OK] [duration:0s]
	//test: Ping(bad) -> [Deadline Exceeded [ping response time out: [urn:ping:bad]]] [duration:0s]
	//test: Ping(error) -> [Internal Error [ping depends error message]] [duration:0s]
	//test: Ping(delay) -> [OK] [duration:0s]

}

func pingGood(c chan core.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}

func pingBad(c chan core.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second * 4)
			core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}

func pingError(c chan core.Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				core.ReplyTo(msg, runtime.NewStatusError(0, pingLocation, err).SetDuration(time.Since(start)))
			}
		default:
		}
	}
}

func pingDelay(c chan core.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second)
			core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}