package gotypes

import "go/types"

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
