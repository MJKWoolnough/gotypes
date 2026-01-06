package gotypes_test

import (
	"fmt"
	"go/types"

	"vimagination.zapto.org/gotypes"
)

func Example() {
	pkg, err := gotypes.ParsePackage(".")
	if err != nil {
		fmt.Println(err)

		return
	}

	z := pkg.Scope().Lookup("zipFS")
	fmt.Println(z)

	for field := range z.Type().Underlying().(*types.Struct).Fields() {
		fmt.Println(field)
	}

	// Output:
	// type vimagination.zapto.org/gotypes.zipFS struct{*archive/zip.Reader; base string}
	// field Reader *archive/zip.Reader
	// field base string
}
