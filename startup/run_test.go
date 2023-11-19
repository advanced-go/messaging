package startup

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/runtime/runtimetest"
	"github.com/advanced-go/messaging/content"
	"github.com/advanced-go/messaging/core"
	"time"
)

var credFn content.Credentials = func() (string, string, error) {
	return "", "", nil
}

var start time.Time

func ExampleCreateToSend() {
	none := "/startup/none"
	one := "/startup/one"

	core.RegisterUnchecked(none, nil)
	core.RegisterUnchecked(one, nil)

	m := createToSend(nil, nil)
	msg := m[none]
	fmt.Printf("test: createToSend(nil,nil) -> [to:%v] [from:%v]\n", msg.To, msg.From)

	cm := content.Map{one: []any{credFn}}
	m = createToSend(cm, nil)
	msg = m[one]
	fmt.Printf("test: createToSend(map,nil) -> [to:%v] [from:%v] [credentials:%v]\n", msg.To, msg.From, content.AccessCredentials(&msg) != nil)

	//Output:
	//test: createToSend(nil,nil) -> [to:/startup/none] [from:startup]
	//test: createToSend(map,nil) -> [to:/startup/one] [from:startup] [credentials:true]

}

func ExampleStartup_Success() {
	uri1 := "urn:startup:good"
	uri2 := "urn:startup:bad"
	uri3 := "urn:startup:depends"

	start = time.Now()
	core.Directory.Empty()

	c := make(chan core.Message, 16)
	core.Register(uri1, c)
	go startupGood(c)

	c = make(chan core.Message, 16)
	core.Register(uri2, c)
	go startupBad(c)

	c = make(chan core.Message, 16)
	core.Register(uri3, c)
	go startupDepends(c, nil)

	status := Run[runtimetest.DebugError](time.Second*2, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//test: Startup() -> [OK]

}

func ExampleStartup_Failure() {
	uri1 := "urn:startup:good"
	uri2 := "urn:startup:bad"
	uri3 := "urn:startup:depends"

	start = time.Now()
	core.Directory.Empty()

	c := make(chan core.Message, 16)
	core.Register(uri1, c)
	go startupGood(c)

	c = make(chan core.Message, 16)
	core.Register(uri2, c)
	go startupBad(c)

	c = make(chan core.Message, 16)
	core.Register(uri3, c)
	go startupDepends(c, errors.New("startup failure error message"))

	status := Run[runtimetest.DebugError](time.Second*2, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//{ "id":null, "l":"github.com/advanced-go/core/runtime/startup/Run", "o":null "err" : [ "startup failure error message" ] }
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
				core.ReplyTo(msg, runtime.NewStatusError(0, runLocation, err).SetDuration(time.Since(start)))
			} else {
				time.Sleep(time.Second + (time.Millisecond * 900))
				core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
			}

		default:
		}
	}
}
