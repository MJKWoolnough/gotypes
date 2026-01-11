package gotypes

import "go/types"

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
