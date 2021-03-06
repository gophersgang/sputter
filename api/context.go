package api

import "sync"

const (
	// AlreadyBound is thrown when an attempt is made to rebind a Name
	AlreadyBound = "symbol is already bound in this context: %s"

	defaultContextEntries = 16
)

type (
	// Context represents a mutable variable scope
	Context interface {
		Value
		Get(Name) (Value, bool)
		Has(Name) (Context, bool)
		Put(Name, Value)
		Delete(Name)
	}

	context struct {
		sync.RWMutex
		parent Context
		vars   Variables
	}

	rootContext struct {
		*context
	}
)

// NewContext creates a new independent Context instance
func NewContext() Context {
	return &rootContext{
		context: &context{
			vars: make(Variables, defaultContextEntries),
		},
	}
}

// ChildContext creates a new child Context of the provided parent
func ChildContext(parent Context) Context {
	return &context{
		parent: parent,
		vars:   make(Variables, defaultContextEntries),
	}
}

// ChildContextVars creates a new child Context with Variables
func ChildContextVars(parent Context, vars Variables) Context {
	return &context{
		parent: parent,
		vars:   vars,
	}
}

// Get retrieves a value from the Context
func (c *rootContext) Get(n Name) (Value, bool) {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.vars[n]; ok {
		return v, ok
	}
	return Nil, false
}

// Has looks up the Context in which a value exists
func (c *rootContext) Has(n Name) (Context, bool) {
	c.RLock()
	defer c.RUnlock()
	if _, ok := c.vars[n]; ok {
		return c, ok
	}
	return nil, false
}

// Get retrieves a value from the Context chain
func (c *context) Get(n Name) (Value, bool) {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.vars[n]; ok {
		return v, true
	}
	return c.parent.Get(n)
}

// Has looks up the Context in which a value exists
func (c *context) Has(n Name) (Context, bool) {
	c.RLock()
	defer c.RUnlock()
	if _, ok := c.vars[n]; ok {
		return c, true
	}
	return c.parent.Has(n)
}

// Put puts a Value into the immediate Context
func (c *context) Put(n Name, v Value) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.vars[n]; ok {
		panic(ErrStr(AlreadyBound, n))
	}
	c.vars[n] = v
}

// Delete removes a Value from the immediate Context
func (c *context) Delete(n Name) {
	c.Lock()
	defer c.Unlock()
	delete(c.vars, n)
}

// Str converts this Value into a Str
func (c *context) Str() Str {
	return MakeDumpStr(c)
}
