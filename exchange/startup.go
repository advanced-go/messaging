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
	startupLocation = PkgPath + ":Startup"
)

// Startup - templated function to start all registered resources.
func Startup[E runtime.ErrorHandler](duration time.Duration, content core.Map) (status runtime.Status) {
	return startup[E](HostDirectory, duration, content)
}

func startup[E runtime.ErrorHandler](directory Directory, duration time.Duration, content core.Map) (status runtime.Status) {
	var e E
	var failures []string
	var count = directory.Count()

	if count == 0 {
		return runtime.StatusOK()
	}
	cache := core.NewMessageCache()
	toSend := createToSend(directory, content, core.NewMessageCacheHandler[E](cache))
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
			return runtime.StatusOK()
		}
		break
	}
	shutdownHost(core.Message{Event: core.ShutdownEvent})
	if len(failures) > 0 {
		handleErrors[E](failures, cache)
		return runtime.NewStatus(http.StatusInternalServerError)
	}
	return e.Handle(runtime.NewStatusError(runtime.StatusDeadlineExceeded, startupLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), directory.Count()))), "", "")
}

func createToSend(directory Directory, cm core.Map, fn core.MessageHandler) core.MessageMap {
	m := make(core.MessageMap)
	for _, k := range directory.List() {
		msg := core.Message{To: k, From: startupLocation, Event: core.StartupEvent, Status: nil, ReplyTo: fn}
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
		directory.SendCtrl(msgs[k])
	}
}

func handleErrors[E runtime.ErrorHandler](failures []string, cache core.MessageCache) {
	var e E
	for _, uri := range failures {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil && !msg.Status.OK() {
			loc := ""
			if msg.Status.Location() != nil && len(msg.Status.Location()) > 0 {
				loc = msg.Status.Location()[0]
			}
			e.Handle(runtime.NewStatusError(http.StatusInternalServerError, loc, msg.Status.Errors()...), "", "")
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
			fmt.Printf("startup successful: [%v] : %s\n", uri, msg.Status.Duration())
		}
	}
}
