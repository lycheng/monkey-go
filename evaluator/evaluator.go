package evaluator

import (
	"fmt"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/object"
)

// Eval returns the object of the AST node
func Eval(node ast.Node) (object.Object, error) {

	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Program:
		return evalStatements(node.Statements)
	}
	return nil, fmt.Errorf("Can not find the match object for %s", node.String())
}

func evalStatements(stmts []ast.Statement) (result object.Object, err error) {
	for _, statement := range stmts {
		result, err = Eval(statement)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
