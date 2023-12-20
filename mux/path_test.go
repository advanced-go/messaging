package mux

import (
	"fmt"
	"reflect"
)

func Example_PkgUri() {
	pkgPath := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath)

	//Output:
	//test: PkgPath = "github.com/advanced-go/messaging/mux"

}
