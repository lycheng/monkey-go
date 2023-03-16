package evaluator

import (
	"fmt"

	"github.com/lycheng/monkey-go/ast"
	"github.com/lycheng/monkey-go/object"
)

var (
	nullObj  = &object.Null{}
	trueObj  = &object.Boolean{Value: true}
	falseObj = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return trueObj
	}
	return falseObj
}

// Eval returns the object of the AST node
func Eval(node ast.Node) (object.Object, error) {

	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value), nil
	case *ast.PrefixExpression:
		{
			right, err := Eval(node.Right)
			if err != nil {
				return nil, err
			}
			return evalPrefixExpression(node.Operator, right)
		}
	case *ast.InfixExpression:
		{
			left, err := Eval(node.Left)
			if err != nil {
				return nil, err
			}
			right, err := Eval(node.Right)
			if err != nil {
				return nil, err
			}
			return evalInfixExpression(node.Operator, left, right), nil
		}
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

func evalPrefixExpression(operator string, right object.Object) (object.Object, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right), nil
	default:
		return nullObj, nil
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case trueObj:
		return falseObj
	case falseObj:
		return trueObj
	case nullObj:
		return trueObj
	default:
		return falseObj
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return nullObj
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return nullObj
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return nullObj
	}
}