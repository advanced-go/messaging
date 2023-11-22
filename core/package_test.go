package core

import (
	"fmt"
	"reflect"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgPath = \"%v\"\n", reflect.TypeOf(any(pkg{})).PkgPath())

	//Output:
	//test: PkgPath = "github.com/advanced-go/messaging/core"

}
