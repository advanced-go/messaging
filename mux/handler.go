package mux

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"github.com/advanced-go/messaging/exchange"
	"net/http"
	"reflect"
)

const (
	PingResource  = "ping"
	ContentLength = "Content-Length"
	ContentType   = "Content-Type"
	//ContentTypeText       = "text/plain" //charset=utf-8; charset=us-ascii"
	writeStatusContentLoc = PkgPath + ":writeStatusContent"
	bytesLoc              = PkgPath + ":writeBytes"
)

type muxEntry struct {
	pattern string
	handler http.HandlerFunc
}

var routes []muxEntry

// Handle - add pattern and Http handler mux entry
// TO DO : panic on duplicate handler and pattern combination
func Handle(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	routes = append(routes, muxEntry{pattern: pattern, handler: handler})
}

// HttpHandler - handler for messaging
func HttpHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil || w == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, rt := range routes {
		nid, rsc, ok := uri.UprootUrn(r.URL.Path)
		if !ok {
			continue
		}
		if nid != rt.pattern {
			continue
		}
		if rsc == PingResource {
			ProcessPing[runtime.Log](w, nid)
			return
		}
		rt.handler(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func ProcessPing[E runtime.ErrorHandler](w http.ResponseWriter, nid string) {
	status := exchange.Ping[E](nil, nid)
	if status.OK() {
		status.SetContent(fmt.Sprintf("Ping status: %v, resource: %v", status, nid), false)
	}
	w.WriteHeader(status.Http())
	writeStatusContent[E](w, status)
}

func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status runtime.Status) {
	var e E

	if status.Content() == nil {
		return
	}
	buf, rc, status1 := writeBytes(status.Content())
	if !status1.OK() {
		e.Handle(status, status.RequestId(), writeStatusContentLoc)
		return
	}
	w.Header().Set(ContentType, rc)
	w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeStatusContentLoc, err), "", "")
	}
}

func writeBytes(content any) ([]byte, string, runtime.Status) {
	var buf []byte

	switch ptr := (content).(type) {
	case []byte:
		buf = ptr
	case string:
		buf = []byte(ptr)
	case error:
		buf = []byte(ptr.Error())
	default:
		return nil, "", runtime.NewStatusError(http.StatusInternalServerError, bytesLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
	}
	return buf, http.DetectContentType(buf), runtime.StatusOK()
}
