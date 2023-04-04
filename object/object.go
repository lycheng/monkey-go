package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lycheng/monkey-go/ast"
)

// object types
const (
	INTEGER     = "INTEGER"
	BOOLEAN     = "BOOLEAN"
	STRING      = "STRING"
	NULL        = "NULL"
	RETURNVALUE = "RETURN_VALUE"
	ERROR       = "ERROR"
	FUNCTION    = "FUNCTION"
)

// Type for object type
type Type string

// Object is a wrapper for object interface
type Object interface {
	Type() Type
	Inspect() string
}

// Integer for int object
type Integer struct {
	Value int64
}

// Inspect returns the int value of Integer
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type returns the object type
func (i *Integer) Type() Type { return INTEGER }

// Boolean for boolean object
type Boolean struct {
	Value bool
}

// Type returns the object type
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns the boolean value of Boolean
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// String for String object
type String struct {
	Value string
}

// Type returns STRING
func (s *String) Type() Type { return STRING }

// Inspect returns string value
func (s *String) Inspect() string { return s.Value }

// Null for null type
type Null struct{}

// Type returns the null
func (n *Null) Type() Type { return NULL }

// Inspect returns the null
func (n *Null) Inspect() string { return "null" }

// ReturnValue for return value
type ReturnValue struct {
	Value Object
}

// Type returns the type of ReturnValue
func (rv *ReturnValue) Type() Type { return RETURNVALUE }

// Inspect returns string of the value
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error struct for eval errors
type Error struct {
	Message string
}

// Type returns the ERROR object
func (e *Error) Type() Type { return ERROR }

// Inspect returns the error message
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns FUNCTION
func (f *Function) Type() Type { return FUNCTION }

// Inspect returns function definition as string
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
