package startup

import (
	"fmt"
	"reflect"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgUri = \"%v\"\n", reflect.TypeOf(any(pkg{})).PkgPath())

	//Output:
	//test: PkgUri = "github.com/advanced-go/messaging/startup"

}

func Example_newStatusRequest() {
	req := newStatusRequest()
	fmt.Printf("test: newStatusRequest() -> %v\n", req.URL.Path)

	//Output:
	//test: newStatusRequest() -> /startup/status

}
