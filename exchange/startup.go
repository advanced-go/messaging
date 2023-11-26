package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"net/http"
	"time"
)

const (
	startupLocation = PkgPath + "/Startup"
)

// Startup - templated function to start all registered resources.
func Startup[E runtime.ErrorHandler](duration time.Duration, content core.Map) (status runtime.Status) {
	return startup[E](exchDir, duration, content)
}

func startup[E runtime.ErrorHandler](directory Directory, duration time.Duration, content core.Map) (status runtime.Status) {
	var e E
	var failures []string
	var count = directory.Count()

	if count == 0 {
		return runtime.NewStatusOK()
	}
	cache := core.NewMessageCache()
	toSend := createToSend(directory, content, core.NewMessageCacheHandler(cache))
	sendMessages(directory, toSend)
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
	Shutdown()
	if len(failures) > 0 {
		handleErrors[E](failures, cache)
		return runtime.NewStatus(http.StatusInternalServerError)
	}
	//return e.Handle("", runLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), directory.Count()))).SetCode(runtime.StatusDeadlineExceeded)
	return e.Handle(runtime.NewStatusError(runtime.StatusDeadlineExceeded, startupLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), directory.Count()))), "", "")
}

func createToSend(directory Directory, cm core.Map, fn core.MessageHandler) core.MessageMap {
	m := make(core.MessageMap)
	for _, k := range directory.List() {
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

func sendMessages(directory Directory, msgs core.MessageMap) {
	for k := range msgs {
		directory.Send(msgs[k])
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