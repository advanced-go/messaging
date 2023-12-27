package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"time"
)

type credentials func() (username string, password string, err error)

var credFn credentials = func() (string, string, error) {
	return "", "", nil
}

// accessCredentials - access function for Credentials in a message
func accessCredentials(msg *core.Message) credentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(credentials); ok {
			return fn
		}
	}
	return nil
}

func accessResource(msg *core.Message) resource {
	if msg == nil || msg.Content == nil {
		return resource{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(resource); ok {
			return url
		}
	}
	return resource{}
}

// Resource - struct for a resource
type resource struct {
	Uri string
}

func testRegister(dir *directory, uri string, cmd, data chan core.Message) runtime.Status {
	if cmd == nil {
		cmd = make(chan core.Message, 16)
	}
	return dir.Add(newMailbox(uri, false, cmd, data))
}

var start time.Time

func ExampleCreateToSend() {
	none := "startup/none"
	one := "startup/one"

	startupDir := any(NewDirectory()).(*directory)
	status := testRegister(startupDir, none, nil, nil)
	if !status.OK() {
		fmt.Printf("test: testRegister() -> [status:%v]\n", status)
	}
	status = testRegister(startupDir, one, nil, nil)
	if !status.OK() {
		fmt.Printf("test: testRegister() -> [status:%v]\n", status)
	}
	m := createToSend(startupDir, nil, nil)
	msg := m[none]
	fmt.Printf("test: createToSend(nil,nil) -> [to:%v] [from:%v]\n", msg.To, msg.From)

	cm := ContentMap{one: []any{credFn}}
	m = createToSend(startupDir, cm, nil)
	msg = m[one]
	fmt.Printf("test: createToSend(map,nil) -> [to:%v] [from:%v] [credentials:%v]\n", msg.To, msg.From, accessCredentials(&msg) != nil)

	//Output:
	//test: createToSend(nil,nil) -> [to:startup/none] [from:github.com/advanced-go/messaging/exchange:Startup]
	//test: createToSend(map,nil) -> [to:startup/one] [from:github.com/advanced-go/messaging/exchange:Startup] [credentials:true]

}

func ExampleStartup_Success() {
	uri1 := "github.com/startup/good"
	uri2 := "github.com/startup/bad"
	uri3 := "github.com/startup/depends"

	startupDir := any(NewDirectory()).(*directory)
	start = time.Now()

	c := make(chan core.Message, 16)
	testRegister(startupDir, uri1, c, nil)
	go startupGood(c)

	c = make(chan core.Message, 16)
	testRegister(startupDir, uri2, c, nil)
	go startupBad(c)

	c = make(chan core.Message, 16)
	testRegister(startupDir, uri3, c, nil)
	go startupDepends(c, nil)

	status := startup[runtime.Output](startupDir, time.Second*2, nil)

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
	startupDir := any(NewDirectory()).(*directory)

	start = time.Now()

	c := make(chan core.Message, 16)
	testRegister(startupDir, uri1, c, nil)
	go startupGood(c)

	c = make(chan core.Message, 16)
	testRegister(startupDir, uri2, c, nil)
	go startupBad(c)

	c = make(chan core.Message, 16)
	testRegister(startupDir, uri3, c, nil)
	go startupDepends(c, errors.New("startup failure error message"))

	status := startup[runtime.Output](startupDir, time.Second*2, nil)

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
			core.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
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
			core.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
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
				core.SendReply(msg, runtime.NewStatusError(0, startupLocation, err).SetDuration(time.Since(start)))
			} else {
				time.Sleep(time.Second + (time.Millisecond * 900))
				core.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
			}

		default:
		}
	}
}

var msgTest = core.Message{To: "to-uri", From: "from-uri", Content: []any{
	"text content",
	500,
	credentials(func() (username, password string, err error) { return "", "", nil }),
	time.Second,
	nil,
	//runtime.Handle[runtime.DebugError](),
	errors.New("this is a content error message"),
	func() bool { return false },
	runtime.NewStatusError(0, "location", errors.New("error message")).SetDuration(time.Second * 2),
	//runtime.HandleWithContext[runtime.DebugError](),
	resource{"postgres://username:password@database.cloud.timescale.com/database?sslmode=require"},
}}

func ExampleAccessCredentials() {
	fmt.Printf("test: AccessCredentials(nil) -> %v\n", accessCredentials(nil) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", accessCredentials(&core.Message{To: "to-uri"}) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", accessCredentials(&msgTest) != nil)

	//Output:
	//test: AccessCredentials(nil) -> false
	//test: AccessCredentials(msg) -> false
	//test: AccessCredentials(msg) -> true
}

func ExampleAccessResource() {
	fmt.Printf("test: AccessResource(nil) -> %v\n", accessResource(nil))
	fmt.Printf("test: AccessResource(msg) -> %v\n", accessResource(&core.Message{To: "to-uri"}))
	fmt.Printf("test: AccessResource(msg) -> %v\n", accessResource(&msgTest))

	//Output:
	//test: AccessResource(nil) -> {}
	//test: AccessResource(msg) -> {}
	//test: AccessResource(msg) -> {postgres://username:password@database.cloud.timescale.com/database?sslmode=require}

}
