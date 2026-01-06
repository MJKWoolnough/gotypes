# gotypes

[![CI](https://github.com/MJKWoolnough/gotypes/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/gotypes/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/gotypes.svg)](https://pkg.go.dev/vimagination.zapto.org/gotypes)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/gotypes)](https://goreportcard.com/report/vimagination.zapto.org/gotypes)

--
    import "vimagination.zapto.org/gotypes"

Package gotypes provides a helper function to parse Go types from source code.

## Highlights

 - Simple parsing of Go source code to enable reading of type information.
 - Automatically handles dependencies, using local cache, stdlib, and remote Go mod proxy.
 - Optionally omit specific source code files from parsing.

## Usage

```go
package gogotypes_test

import (
	"fmt"
	"go/types"

	"vimagination.zapto.org/gogotypes"
)

func Example() {
	pkg, err := gogotypes.ParsePackage(".")
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
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/gotypes
