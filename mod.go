package gotypes

import (
	"os"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type ModFile struct {
	Module  string
	Path    string
	Imports map[string]module.Version
}

func ParseModFile(path string) (*ModFile, error) {
	return parseModFile(&osFS{os.DirFS(path).(statReadDirFileFS)}, path)
}

func parseModFile(fsys filesystem, path string) (*ModFile, error) {
	data, err := fsys.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}

	imports := make(map[string]module.Version, len(f.Require))

	for _, r := range f.Require {
		imports[r.Mod.Path] = r.Mod
	}

	for _, r := range f.Replace {
		if m, ok := imports[r.Old.Path]; ok && (r.Old.Version == "" || r.Old.Version == m.Version) {
			imports[r.Old.Path] = r.New
		}
	}

	return &ModFile{
		Module:  f.Module.Mod.Path,
		Path:    path,
		Imports: imports,
	}, nil
}
