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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

// Eval returns the object of the AST node
func Eval(node ast.Node) (object.Object, error) {

	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value), nil
	case *ast.ReturnStatement:
		val, err := Eval(node.ReturnValue)
		if err != nil {
			return nil, err
		}
		if isError(val) {
			return val, nil
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.PrefixExpression:
		{
			right, err := Eval(node.Right)
			if err != nil {
				return nil, err
			}
			if isError(right) {
				return right, nil
			}
			return evalPrefixExpression(node.Operator, right)
		}
	case *ast.InfixExpression:
		{
			left, err := Eval(node.Left)
			if err != nil {
				return nil, err
			}
			if isError(left) {
				return left, nil
			}
			right, err := Eval(node.Right)
			if err != nil {
				return nil, err
			}
			if isError(right) {
				return right, nil
			}
			return evalInfixExpression(node.Operator, left, right), nil
		}
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.Program:
		return evalProgram(node)
	}
	return nil, fmt.Errorf("Can not find the match object for %s", node.String())
}

func evalBlockStatements(block *ast.BlockStatement) (result object.Object, err error) {
	for _, statement := range block.Statements {
		result, err = Eval(statement)
		if err != nil {
			return nil, err
		}

		if result.Type() == object.RETURNVALUE || result.Type() == object.ERROR {
			return result, nil
		}
	}
	return result, nil
}

func evalStatements(stmts []ast.Statement) (result object.Object, err error) {
	for _, statement := range stmts {
		result, err = Eval(statement)
		if err != nil {
			return nil, err
		}

		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value, nil
		}
	}
	return result, nil
}

func evalProgram(program *ast.Program) (result object.Object, err error) {
	for _, statement := range program.Statements {
		result, err = Eval(statement)
		if err != nil {
			return nil, err
		}

		switch rv := result.(type) {
		case *object.ReturnValue:
			return rv.Value, nil
		case *object.Error:
			return rv, nil
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
		return newError("unknown operator: -%s", right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case nullObj:
		return false
	case trueObj:
		return true
	case falseObj:
		return false
	default:
		return true
	}
}

func evalIfExpression(ie *ast.IfExpression) (object.Object, error) {
	condition, err := Eval(ie.Condition)
	if isError(condition) {
		return condition, nil
	}
	if err != nil {
		return nil, err
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}
	return nullObj, nil
}
