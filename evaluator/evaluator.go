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
func Eval(node ast.Node, env *object.Environment) object.Object {

	switch node := node.(type) {
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.PrefixExpression:
		{
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalPrefixExpression(node.Operator, right)
		}
	case *ast.InfixExpression:
		{
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalInfixExpression(node.Operator, left, right)
		}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Program:
		return evalProgram(node, env)
	}
	return nullObj
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) (result object.Object) {
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result.Type() == object.RETURNVALUE || result.Type() == object.ERROR {
			return result
		}
	}
	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) (result object.Object) {
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch rv := result.(type) {
		case *object.ReturnValue:
			return rv.Value
		case *object.Error:
			return rv
		}
	}
	return result
}

func evalProgram(program *ast.Program, env *object.Environment) (result object.Object) {
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch rv := result.(type) {
		case *object.ReturnValue:
			return rv.Value
		case *object.Error:
			return rv
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nullObj
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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return nullObj
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}
