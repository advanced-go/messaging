package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"time"
)

var credFn core.Credentials = func() (string, string, error) {
	return "", "", nil
}

func testRegister(uri string, cmd, data chan core.Message) error {
	add(startupDir, newMailbox(uri, cmd, data))
	return nil
}

var start time.Time

var startupDir = any(NewDirectory()).(*directory)

func ExampleCreateToSend() {
	none := "startup/none"
	one := "startup/one"

	testRegister(none, nil, nil)
	testRegister(one, nil, nil)

	m := createToSend(startupDir, nil, nil)
	msg := m[none]
	fmt.Printf("test: createToSend(nil,nil) -> [to:%v] [from:%v]\n", msg.To, msg.From)

	cm := core.Map{one: []any{credFn}}
	m = createToSend(startupDir, cm, nil)
	msg = m[one]
	fmt.Printf("test: createToSend(map,nil) -> [to:%v] [from:%v] [credentials:%v]\n", msg.To, msg.From, core.AccessCredentials(&msg) != nil)

	//Output:
	//test: createToSend(nil,nil) -> [to:startup/none] [from:github.com/advanced-go/messaging/exchange:Startup]
	//test: createToSend(map,nil) -> [to:startup/one] [from:github.com/advanced-go/messaging/exchange:Startup] [credentials:true]

}

func ExampleStartup_Success() {
	uri1 := "github.com/startup/good"
	uri2 := "github.com/startup/bad"
	uri3 := "github.com/startup/depends"

	start = time.Now()
	empty(startupDir)

	c := make(chan core.Message, 16)
	testRegister(uri1, c, nil)
	go startupGood(c)

	c = make(chan core.Message, 16)
	testRegister(uri2, c, nil)
	go startupBad(c)

	c = make(chan core.Message, 16)
	testRegister(uri3, c, nil)
	go startupDepends(c, nil)

	status := startup[runtime.TestError](startupDir, time.Second*2, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//startup successful: [github.com/startup/bad] : 0s
	//startup successful: [github.com/startup/depends] : 0s
	//startup successful: [github.com/startup/good] : 0s
	//test: Startup() -> [OK]

}

func ExampleStartup_Failure() {
	uri1 := "github.com/startup/good"
	uri2 := "github.com/startup/bad"
	uri3 := "github.com/startup/depends"

	start = time.Now()
	empty(startupDir)

	c := make(chan core.Message, 16)
	testRegister(uri1, c, nil)
	go startupGood(c)

	c = make(chan core.Message, 16)
	testRegister(uri2, c, nil)
	go startupBad(c)

	c = make(chan core.Message, 16)
	testRegister(uri3, c, nil)
	go startupDepends(c, errors.New("startup failure error message"))

	status := startup[runtime.TestError](startupDir, time.Second*2, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//test: Startup() -> [Internal Error]

}

func startupGood(c chan core.Message) {
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

func startupBad(c chan core.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second + time.Millisecond*100)
			core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}

func startupDepends(c chan core.Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				core.ReplyTo(msg, runtime.NewStatusError(0, startupLocation, err).SetDuration(time.Since(start)))
			} else {
				time.Sleep(time.Second + (time.Millisecond * 900))
				core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
			}

		default:
		}
	}
}
