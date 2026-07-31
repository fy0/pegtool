package ast

import (
	"reflect"
	"testing"
)

func TestInspectRecoveryAndThrowExpressions(t *testing.T) {
	pos := Pos{Line: 1, Col: 1}
	regular := NewThrowExpr(pos)
	regular.Label = "expected"
	recovery := NewRecoveryExpr(pos)
	recovery.Expr = regular
	recovery.RecoverExpr = NewLitMatcher(pos, "fallback")

	var visited []string
	Inspect(recovery, func(expr Expression) bool {
		switch expr.(type) {
		case *RecoveryExpr:
			visited = append(visited, "recovery")
		case *ThrowExpr:
			visited = append(visited, "throw")
		case *LitMatcher:
			visited = append(visited, "literal")
		}
		return true
	})

	want := []string{"recovery", "throw", "literal"}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited expressions = %v, want %v", visited, want)
	}
}
