package object

// Environment for objects map
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnclosedEnvironment Return new env with provided env as outer
func NewEnclosedEnvironment(env *Environment) *Environment {
	ne := NewEnvironment()
	ne.outer = env
	return ne
}

// NewEnvironment returns new map for name and object
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Get object from map
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set object into map
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
