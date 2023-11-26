package mux

import (
	"github.com/advanced-go/core/http2"
	"net/http"
)

type route struct {
	pattern string
	handler http.HandlerFunc
}

var routes []route

func Handle(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	routes = append(routes, route{pattern: pattern, handler: handler})
}

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
