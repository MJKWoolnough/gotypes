package gotypes

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"runtime"
	"testing"
)

func TestIsRecursive(t *testing.T) {
	for n, test := range [...]struct {
		input       string
		isRecursive bool
	}{
		{ // 1
			"package a\n\ntype a struct { B b }\n\ntype b struct {c *b}",
			false,
		},
		{ // 2
			"package a\n\ntype a struct {B *a}",
			true,
		},
		{ // 3
			"package a\n\ntype a struct {B map[string]int}",
			false,
		},
		{ // 4
			"package a\n\ntype a struct {B map[*a]int}",
			true,
		},
		{ // 5
			"package a\n\ntype a struct {B map[string]a}",
			true,
		},
		{ // 6
			"package a\n\ntype a struct {B []a}",
			true,
		},
		{ // 7
			"package a\n\ntype a struct {B [2]*a}",
			true,
		},
		{ // 8
			"package a\n\ntype a struct {b func() int}",
			false,
		},
		{ // 9
			"package a\n\ntype a struct {b func() a}",
			true,
		},
		{ // 10
			"package a\n\ntype a struct {b func(int) }",
			false,
		},
		{ // 11
			"package a\n\ntype a struct {b func(a) }",
			true,
		},
		{ // 12
			"package a\n\ntype a struct { a b }\ntype b interface {C() int}",
			false,
		},
		{ // 13
			"package a\n\ntype a struct { a b }\ntype b interface {C() b}",
			false,
		},
		{ // 14
			"package a\n\ntype a interface { A() a }",
			true,
		},
		{ // 15
			"package a\n\ntype a interface { A() b }\ntype b interface { B() a\n}",
			true,
		},
		{ // 16
			"package a\n\ntype a interface { A() b }\ntype b struct { B a\n}",
			true,
		},
		{ // 17
			"package a\n\ntype a[T b] struct { A T }\ntype b interface { C() bool\n}",
			false,
		},
		{ // 18
			"package a\n\ntype a[T b] struct { A T }\ntype b interface { C() *a[d]\n}\ntype d struct {}\nfunc(d) C() *a[d]{return nil}",
			true,
		},
	} {
		if self := parseType(t, test.input); IsTypeRecursive(self) != test.isRecursive {
			t.Errorf("test %d: didn't get expected recursive value: %v", n+1, test.isRecursive)
		}
	}
}

func parseType(t *testing.T, input string) types.Type {
	t.Helper()

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "a.go", input, parser.AllErrors|parser.ParseComments)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	conf := types.Config{
		GoVersion: runtime.Version(),
		Importer:  importer.ForCompiler(fset, runtime.Compiler, nil),
	}

	pkg, err := conf.Check("a", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	return pkg.Scope().Lookup("a").Type()
}
