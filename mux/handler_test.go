package mux

import (
	"fmt"
	"github.com/advanced-go/core/http2"
	"net/http"
)

func appHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func Example_HttpHandler() {
	pattern := "github.com/advanced-go/example-domain/activity"
	r, _ := http.NewRequest("PUT", "http://localhost:8080/github.com/advanced-go/example-domain/activity:entry", nil)

	Handle(pattern, appHttpHandler)

	rec := http2.NewRecorder()

	HttpHandler(rec, r)

	fmt.Printf("test: HttpHandler() -> %v\n", rec.Result().StatusCode)

	//Output:
	//test: HttpHandler() -> 504

}
