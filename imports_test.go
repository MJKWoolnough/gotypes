package gotypes

import (
	"fmt"
	"path"
	"testing"
)

func TestGetImports(t *testing.T) {
	pkg, err := ParsePackage(".")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	imps := Imports(pkg)

	for path, name := range map[string]string{
		"vimagination.zapto.org/gotypes":      "gotypes",
		"vimagination.zapto.org/httpreaderat": "httpreaderat",
		"vimagination.zapto.org/cache":        "cache",
		"archive/zip":                         "zip",
		"bytes":                               "bytes",
	} {
		if n := imps[path].Name(); n != name {
			t.Errorf("expecting package %q to have name %q, got %q", path, name, n)
		}
	}

	for p, pkg := range imps {
		if path.Base(p) != pkg.Name() {
			fmt.Println(p, pkg.Name())
		}
	}
}
