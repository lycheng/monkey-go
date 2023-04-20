package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	BUILTIN     = "BUILTIN"
	ARRAY       = "ARRAY"
	HASH        = "HASH"
)

// Type for object type
type Type string

// Object is a wrapper for object interface
type Object interface {
	Type() Type
	Inspect() string
}

// Hashable for objects that can be used as key of hash object
type Hashable interface {
	HashKey() HashKey
}

// Integer for int object
type Integer struct {
	Value int64
}

// Inspect returns the int value of Integer
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type returns the object type
func (i *Integer) Type() Type { return INTEGER }

// HashKey for Integer type as hash type's key
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Boolean for boolean object
type Boolean struct {
	Value bool
}

// Type returns the object type
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns the boolean value of Boolean
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// HashKey for Boolean type as hash type's key
func (b *Boolean) HashKey() HashKey {
	var val uint64
	if b.Value {
		val = 1
	} else {
		val = 0
	}
	return HashKey{Type: b.Type(), Value: val}
}

// String for String object
type String struct {
	Value string
}

// Type returns STRING
func (s *String) Type() Type { return STRING }

// Inspect returns string value
func (s *String) Inspect() string { return s.Value }

// HashKey for string type as hash type's key
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

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

// BuiltinFunction for Built-In function definition
type BuiltinFunction func(args ...Object) Object

// Builtin for Built-In function object
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns BUILTIN
func (b *Builtin) Type() Type { return BUILTIN }

// Inspect returns built in message
func (b *Builtin) Inspect() string { return "Built-In function" }

// Array object
type Array struct {
	Elements []Object
}

// Type returns ARRAY
func (ao *Array) Type() Type { return ARRAY }

// Inspect returns array literal value
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// HashKey for hash type's key
type HashKey struct {
	Type  Type
	Value uint64
}

// HashPair for hash type's value
type HashPair struct {
	Key Object
	Val Object
}

// Hash for hash type
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns ARRAY
func (h *Hash) Type() Type { return HASH }

// Inspect returns hash type literal value
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Val.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
