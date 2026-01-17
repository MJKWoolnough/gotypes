package gotypes

import "testing"

func TestParseModFile(t *testing.T) {
	tfs := testFS{
		"go.mod": `module vimagination.zapto.org/marshal

go 1.25.5

require (
	golang.org/x/mod v0.31.0
	golang.org/x/tools v0.40.0
)

require golang.org/x/sync v0.19.0 // indirect

replace golang.org/x/tools => somewhere.org/tools v0.1.0
`,
	}

	if pkg, err := parseModFile(tfs, ""); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if pkg.Module != "vimagination.zapto.org/marshal" {
		t.Errorf("expecting path %q, got %q", "vimagination.zapto.org/marshal", pkg.Module)
	} else if len(pkg.Imports) != 3 {
		t.Errorf("expecting 3 imports, got %d", len(pkg.Imports))
	} else if m := pkg.Imports["golang.org/x/mod"]; m.Path != "golang.org/x/mod" {
		t.Errorf("expecting url for %q to be %q, got %q", "golang.org/x/mod", "golang.org/x/mod", m.Path)
	} else if m.Version != "v0.31.0" {
		t.Errorf("expecting version for %q to be %q, got %q", "golang.org/x/mod", "v0.31.0", m.Version)
	} else if m = pkg.Imports["golang.org/x/tools"]; m.Path != "somewhere.org/tools" {
		t.Errorf("expecting url for %q to be %q, got %q", "golang.org/x/tools", "somewhere.org/tools", m.Path)
	} else if m.Version != "v0.1.0" {
		t.Errorf("expecting version for %q to be %q, got %q", "golang.org/x/tools", "v0.1.0", m.Version)
	}
}
