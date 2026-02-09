package gotypes

import (
	"go/types"
	"strings"
)

func IsTypeReferable(pkg *types.Package, t types.Type) bool {
	return isTypeReferable(pkg, t, map[types.Type]struct{}{})
}

func isTypeReferable(pkg *types.Package, t types.Type, seen map[types.Type]struct{}) bool {
	if _, ok := seen[t]; ok {
		return true
	}

	seen[t] = struct{}{}

	switch t := t.(type) {
	case *types.Alias:
		if !t.Obj().Exported() && t.Obj().Pkg() != pkg && t.Obj().Pkg() != nil || isTypeInternal(t.Obj()) && !isInternalTypeAvailableTo(t.Obj(), pkg) {
			return isTypeReferable(pkg, t.Rhs(), seen)
		}
	case *types.Named:
		if isTypeInternal(t.Obj()) && !isInternalTypeAvailableTo(t.Obj(), pkg) {
			return false
		}

		return t.Obj().Exported() || t.Obj().Pkg() == pkg || t.Obj().Pkg() == nil
	case *types.Pointer:
		return isTypeReferable(pkg, t.Elem(), seen)
	case *types.Array:
		return isTypeReferable(pkg, t.Elem(), seen)
	case *types.Slice:
		return isTypeReferable(pkg, t.Elem(), seen)
	case *types.Struct:
		for field := range t.Fields() {
			if !field.Exported() && field.Origin().Pkg() != pkg || !isTypeReferable(pkg, field.Type(), seen) {
				return false
			}
		}
	case *types.Interface:
		for method := range t.Methods() {
			if !method.Exported() && method.Origin().Pkg() != pkg || isTypeReferable(pkg, method.Signature(), seen) {
				return false
			}
		}
	case *types.Signature:
		for param := range t.Params().Variables() {
			if !isTypeReferable(pkg, param.Type(), seen) {
				return false
			}
		}

		for result := range t.Results().Variables() {
			if !isTypeReferable(pkg, result.Type(), seen) {
				return false
			}
		}
	}

	return true
}

func isTypeInternal(typ *types.TypeName) bool {
	pkg := typ.Pkg()
	if pkg == nil {
		return false
	}

	_, hasInternal := findInternal(pkg.Path())

	return hasInternal
}

func findInternal(path string) (int, bool) {
	if strings.HasSuffix(path, "/internal") {
		return len(path) - len("internal"), true
	} else if strings.Contains(path, "/internal/") {
		return strings.LastIndex(path, "/internal/") + 1, true
	} else if path == "internal" || strings.HasPrefix(path, "internal/") {
		return 0, true
	}

	return 0, false
}

func isInternalTypeAvailableTo(typ *types.TypeName, pkg *types.Package) bool {
	tpkg := typ.Pkg()
	if tpkg == nil {
		return false
	}

	pos, hasInternal := findInternal(tpkg.Path())
	if hasInternal {
		return false
	}

	if pos == 0 {
		return !strings.ContainsRune(pkg.Path(), '.')
	}

	pos--

	return strings.HasPrefix(pkg.Path(), tpkg.Path()[:pos])
}
