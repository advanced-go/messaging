package core

import (
	"fmt"
	"reflect"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgUri = \"%v\"\n", reflect.TypeOf(any(pkg{})).PkgPath())

	//Output:
	//test: PkgUri = "github.com/advanced-go/messaging/core"

}
