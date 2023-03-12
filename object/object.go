package object

import "fmt"

// object types
const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
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

// Null for null type
type Null struct{}

// Type returns the null
func (n *Null) Type() Type { return NULL }

// Inspect returns the null
func (n *Null) Inspect() string { return "null" }
