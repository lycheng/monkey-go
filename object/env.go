package object

// Environment for objects map
type Environment struct {
	store map[string]Object
}

// NewEnvironment returns new map for name and object
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get object from map
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set object into map
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
