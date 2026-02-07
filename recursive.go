package gotypes

import "go/types"

// IsTypeRecursive checks the type, and all referenced types, for a reference
// back to itself.
func IsTypeRecursive(typ types.Type) bool {
	return isTypeRecursive(typ, map[types.Type]bool{})
}

func isTypeRecursive(typ types.Type, found map[types.Type]bool) bool {
	f, ok := found[typ]
	if ok {
		return f
	}

	found[typ] = len(found) == 0

	switch t := typ.(type) {
	case *types.Named:
		if params := t.TypeParams(); params != nil {
			if len(found) == 1 {
				clear(found)
			}

			for param := range params.TypeParams() {
				if isTypeRecursive(param.Constraint(), found) {
					return true
				}
			}
		}

		return isTypeRecursive(t.Underlying(), found)
	case *types.Struct:
		for field := range t.Fields() {
			if isTypeRecursive(field.Type(), found) {
				return true
			}
		}
	case *types.Pointer:
		return isTypeRecursive(t.Elem(), found)
	case *types.Map:
		if isTypeRecursive(t.Key(), found) {
			return true
		}

		return isTypeRecursive(t.Elem(), found)
	case *types.Array:
		return isTypeRecursive(t.Elem(), found)
	case *types.Slice:
		return isTypeRecursive(t.Elem(), found)
	case *types.Signature:
		for typ := range t.Params().Variables() {
			if isTypeRecursive(typ.Type(), found) {
				return true
			}
		}

		for typ := range t.Results().Variables() {
			if isTypeRecursive(typ.Type(), found) {
				return true
			}
		}
	case *types.Interface:
		for typ := range t.EmbeddedTypes() {
			if isTypeRecursive(typ, found) {
				return true
			}
		}

		for fn := range t.ExplicitMethods() {
			if isTypeRecursive(fn.Signature(), found) {
				return true
			}
		}
	}

	return false
}
