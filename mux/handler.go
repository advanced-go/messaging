package mux

import (
	"github.com/advanced-go/core/http2"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/exchange"
	"net/http"
)

const (
	PingResource = "ping"
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
	// Check for ping request
	//if strings.HasPrefix(r.URL.Path, PingPrefix) {
	//	processPing(w, r)
	//	return
	//}
	for _, rt := range routes {
		nid, rsc, ok := http2.UprootUrn(r.URL.Path)
		if !ok {
			continue
		}
		if nid != rt.pattern {
			continue
		}
		if rsc == PingResource {
			processPing[runtime.LogError](w, nid)
			return
		}
		rt.handler(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func processPing[E runtime.ErrorHandler](w http.ResponseWriter, nid string) {
	status := exchange.Ping[E](nil, nid)
	http2.WriteResponse[E](w, nil, status, nil)
}
