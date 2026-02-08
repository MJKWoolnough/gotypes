package gotypes

import (
	"errors"
	"fmt"
	"go/types"
	"strings"
)

// Imports walks the dependency tree for the given package and returns every
// dependent package, recursively; the key of the returned map is the import
// path.
func Imports(pkg *types.Package) map[string]*types.Package {
	imps := map[string]*types.Package{pkg.Path(): pkg}

	getAllImports(pkg.Imports(), imps)

	return imps
}

func getAllImports(imports []*types.Package, imps map[string]*types.Package) {
	for _, imp := range imports {
		if _, ok := imps[imp.Path()]; ok {
			continue
		}

		imps[imp.Path()] = imp

		getAllImports(imp.Imports(), imps)
	}
}

func Lookup(imports map[string]*types.Package, typeName string) (types.Object, error) {
	var pkg *types.Package

	if pos := strings.LastIndexByte(typeName, '.'); pos >= 0 {
		pkg = imports[typeName[:pos]]
		if pkg == nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidImport, typeName)
		}

		typeName = typeName[pos+1:]
	} else {
		return nil, fmt.Errorf("%w: %s", ErrInvalidImport, typeName)
	}

	obj := pkg.Scope().Lookup(typeName)
	if obj == nil {
		return nil, ErrUnknownIdentifier
	}

	return obj, nil
}

var (
	ErrInvalidImport     = errors.New("invalid import path")
	ErrUnknownIdentifier = errors.New("unknown identifier")
)
