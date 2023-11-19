package startup

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"net/http"
	"time"
)

const (
	maxWait = time.Second * 2
)

var pingLocation = PkgUri + "/Ping"

// Ping - templated function to "ping" a startup
func Ping[E runtime.ErrorHandler](ctx context.Context, uri string) (status runtime.Status) {
	var e E

	if uri == "" {
		//return e.Handle(runtime.RequestId(ctx), pingLocation, errors.New("invalid argument: startup uri is empty"))
		return e.Handle(runtime.NewStatusError(http.StatusInternalServerError, pingLocation, errors.New("invalid argument: startup uri is empty")),
			runtime.RequestId(ctx), "")

	}
	cache := core.NewMessageCache()
	msg := core.Message{To: uri, From: core.HostName, Event: core.PingEvent, Status: nil, ReplyTo: core.NewMessageCacheHandler(cache)}
	err := core.Directory.Send(msg)
	if err != nil {
		//return e.Handle(runtime.RequestId(ctx), pingLocation, err)
		return e.Handle(runtime.NewStatusError(http.StatusInternalServerError, pingLocation, err), runtime.RequestId(ctx), "")

	}
	duration := maxWait
	for wait := time.Duration(float64(duration) * 0.20); duration >= 0; duration -= wait {
		time.Sleep(wait)
		result, err1 := cache.Get(uri)
		if err1 != nil {
			continue //return e.HandleWithContext(ctx, pingLocation, err1)
		}
		if result.Status == nil {
			//return e.Handle(runtime.RequestId(ctx), pingLocation, errors.New(fmt.Sprintf("ping response status not available: [%v]", uri))).SetCode(runtime.StatusNotProvided)
			return e.Handle(runtime.NewStatusError(http.StatusInternalServerError, pingLocation, errors.New(fmt.Sprintf("ping response status not available: [%v]", uri))), runtime.RequestId(ctx), "")

		}
		return result.Status
	}
	//return e.Handle(runtime.RequestId(ctx), pingLocation, errors.New(fmt.Sprintf("ping response time out: [%v]", uri))).SetCode(runtime.StatusDeadlineExceeded)
	return e.Handle(runtime.NewStatusError(runtime.StatusDeadlineExceeded, pingLocation, errors.New(fmt.Sprintf("ping response time out: [%v]", uri))), runtime.RequestId(ctx), "")

}
