package ast

import (
	"reflect"
	"testing"
)

func TestWalkRecoveryAndThrowExpressions(t *testing.T) {
	expr := &RecoveryExpr{
		Expr:        &LitMatcher{},
		RecoverExpr: &ThrowExpr{},
	}

	var visited []string
	Inspect(expr, func(expr Expression) bool {
		if expr == nil {
			return false
		}
		visited = append(visited, reflect.TypeOf(expr).String())
		return true
	})

	want := []string{"*ast.RecoveryExpr", "*ast.LitMatcher", "*ast.ThrowExpr"}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("unexpected walk order: want %v, got %v", want, visited)
	}
}
