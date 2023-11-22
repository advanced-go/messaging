package startup

import (
	"fmt"
	"net/http"
)

type pkg struct{}

const (
	PkgPath    = "github.com/advanced-go/messaging/startup"
	StatusPath = "/startup/status"
)

var StatusRequest = newStatusRequest()

func newStatusRequest() *http.Request {
	req, err := http.NewRequest("", StatusPath, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return req
}
