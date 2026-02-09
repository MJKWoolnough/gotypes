package gotypes

import (
	"go/types"
	"os"
	"path/filepath"
	"testing"
)

func TestIsTypeReferable(t *testing.T) {
	tmp := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/main\n\ngo 1.25.5"), 0600); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for n, test := range [...]struct {
		input       string
		typeName    string
		isReferable bool
	}{
		{
			"type A struct{}",
			"example.com/main.A",
			true,
		},
		{
			"type a struct{}",
			"example.com/main.a",
			true,
		},
		{
			"type A struct{b}\ntype b struct{C int}",
			"example.com/main.A",
			true,
		},
		{
			"type A struct{b error}",
			"example.com/main.A",
			true,
		},
		{
			"import \"io\"\ntype A struct{lr io.LimitedReader}",
			"example.com/main.A",
			true,
		},
		{
			"import \"io\"\ntype A struct{lr io.LimitedReader}",
			"io.LimitedReader",
			true,
		},
		{
			"import \"os\"\ntype A struct{lr *os.File}",
			"example.com/main.A",
			true,
		},
		{
			"import \"os\"\ntype A struct{lr *os.File}",
			"os.File",
			false,
		},
		{
			"import \"os\"\ntype A = os.File",
			"example.com/main.A",
			true,
		},
	} {
		if err := os.WriteFile(filepath.Join(tmp, "a.go"), []byte("package a\n"+test.input), 0600); err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}

		pkg, err := ParsePackage(tmp)
		if err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}

		obj, err := Lookup(Imports(pkg), test.typeName)
		if err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}

		var typ types.Type

		if alias, ok := obj.Type().(*types.Alias); ok {
			typ = alias.Rhs()
		} else {
			typ = obj.Type().Underlying()
		}

		if isReferable := IsTypeReferable(pkg, typ); isReferable != test.isReferable {
			t.Errorf("test %d: expecting isRepresentable to be %v, got %v", n+1, test.isReferable, isReferable)
		}
	}
}
