package evaluator

import (
	"testing"

	"github.com/lycheng/monkey-go/lexer"
	"github.com/lycheng/monkey-go/object"
	"github.com/lycheng/monkey-go/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Errorf("Eval get error: %s", err)
			break
		}
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func testEval(input string) (object.Object, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
