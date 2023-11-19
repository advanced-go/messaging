package startup

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/content"
	"github.com/advanced-go/messaging/core"
	"net/http"
	"time"
)

type messageMap map[string]core.Message

var runLocation = PkgUri + "/Run"

// Run - templated function to start all registered resources.
func Run[E runtime.ErrorHandler](duration time.Duration, content content.Map) (status runtime.Status) {
	var e E
	var failures []string
	var count = core.Directory.Count()

	if count == 0 {
		return runtime.NewStatusOK()
	}
	cache := core.NewMessageCache()
	toSend := createToSend(content, core.NewMessageCacheHandler(cache))
	sendMessages(toSend)
	for wait := time.Duration(float64(duration) * 0.25); duration >= 0; duration -= wait {
		time.Sleep(wait)
		// Check for completion
		if cache.Count() < count {
			continue
		}
		// Check for failed resources
		failures = cache.Exclude(core.StartupEvent, http.StatusOK)
		if len(failures) == 0 {
			handleStatus(cache)
			return runtime.NewStatusOK()
		}
		break
	}
	core.Shutdown()
	if len(failures) > 0 {
		handleErrors[E](failures, cache)
		return runtime.NewStatus(http.StatusInternalServerError)
	}
	//return e.Handle("", runLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), directory.Count()))).SetCode(runtime.StatusDeadlineExceeded)
	return e.Handle(runtime.NewStatusError(runtime.StatusDeadlineExceeded, runLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), core.Directory.Count()))), "", "")
}

func createToSend(cm content.Map, fn core.MessageHandler) messageMap {
	m := make(messageMap)
	for _, k := range core.Directory.Uri() {
		msg := core.Message{To: k, From: core.HostName, Event: core.StartupEvent, Status: nil, ReplyTo: fn}
		if cm != nil {
			if content, ok := cm[k]; ok {
				msg.Content = append(msg.Content, content...)
			}
		}
		m[k] = msg
	}
	return m
}

func sendMessages(msgs messageMap) {
	for k := range msgs {
		core.Directory.Send(msgs[k])
	}
}

func handleErrors[E runtime.ErrorHandler](failures []string, cache core.MessageCache) {
	var e E
	for _, uri := range failures {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			//.Handle("", msg.Status.Location()[0], msg.Status.Errors()...)
			e.Handle(runtime.NewStatusError(http.StatusInternalServerError, msg.Status.Location()[0], msg.Status.Errors()...), "", "")

		}
	}
}

func handleStatus(cache core.MessageCache) {
	for _, uri := range cache.Uri() {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			fmt.Printf("startup successful for startup [%v] : %s\n", uri, msg.Status.Duration())
		}
	}
}
