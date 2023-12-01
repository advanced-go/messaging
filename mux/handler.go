package mux

import (
	"github.com/advanced-go/core/http2"
	"net/http"
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
	for _, rt := range routes {
		nid, _, ok := http2.UprootUrn(r.URL.Path)
		if !ok {
			continue
		}
		if nid == rt.pattern {
			rt.handler(w, r)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
