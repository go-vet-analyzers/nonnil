package nonnil

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// These white-box tests exercise the defensive nil-guards in the helper
// functions. Those guards only fire on a type-incomplete AST node (an ident or
// composite literal with no recorded type), which the compiler never produces
// for code that type-checks — and analysistest only runs against packages that
// type-check. They are reached here directly by handing the helpers synthetic
// nodes that are absent from the (empty) TypesInfo maps, or a struct whose field
// set does not contain the requested name.

func emptyPass() *analysis.Pass {
	return &analysis.Pass{
		TypesInfo: &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{},
			Defs:  map[*ast.Ident]types.Object{},
		},
	}
}

// defType returns nil when the identifier has no entry in TypesInfo.Defs
// (only possible for an AST that did not fully type-check).
func TestDefTypeNoDef(t *testing.T) {
	id := ast.NewIdent("x")
	if got := defType(emptyPass(), id); got != nil {
		t.Fatalf("defType on undefined ident = %v, want nil", got)
	}
}

// typeOf returns nil for an expression with no recorded type, which makes
// checkComposite take its early-return guard.
func TestTypeOfNoType(t *testing.T) {
	lit := &ast.CompositeLit{}
	if got := typeOf(emptyPass(), lit); got != nil {
		t.Fatalf("typeOf on untyped expr = %v, want nil", got)
	}
}

// checkComposite returns immediately when the literal has no recorded type.
// It must not panic and must report nothing.
func TestCheckCompositeNoType(t *testing.T) {
	lit := &ast.CompositeLit{
		Elts: []ast.Expr{ast.NewIdent("nil")},
	}
	// A nil Report would panic if checkComposite tried to report; the empty
	// TypesInfo guarantees typeOf(lit) == nil, so it returns before that.
	checkComposite(emptyPass(), lit)
}

// structFieldType returns nil when the requested field name is not part of the
// struct (a keyed literal naming a non-existent field never compiles, so this
// guard is unreachable via analysistest).
func TestStructFieldTypeMissing(t *testing.T) {
	pkg := types.NewPackage("p", "p")
	field := types.NewField(token.NoPos, pkg, "X", types.Typ[types.Int], false)
	st := types.NewStruct([]*types.Var{field}, nil)

	if got := structFieldType(st, "X"); got == nil {
		t.Fatal("structFieldType on existing field = nil, want the field type")
	}
	if got := structFieldType(st, "Missing"); got != nil {
		t.Fatalf("structFieldType on missing field = %v, want nil", got)
	}
}
