package gotypes

import (
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
}

func TestLookup(t *testing.T) {
	pkg, err := ParsePackage(".")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	imps := Imports(pkg)

	for path, name := range map[string]string{
		"vimagination.zapto.org/gotypes":      "Imports",
		"vimagination.zapto.org/httpreaderat": "Request",
		"vimagination.zapto.org/cache":        "LRU",
		"archive/zip":                         "File",
	} {
		if typ, err := Lookup(imps, path+"."+name); err != nil {
			t.Errorf("import %s: unexpected error: %s", path, err)
		} else if typ.Name() != name {
			t.Errorf("expecting package %q to have type name %q, got %q", path, name, typ.Name())
		}
	}
}
