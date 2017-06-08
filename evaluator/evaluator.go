package evaluator

import a "github.com/kode4food/sputter/api"

// Read converts the raw source into unexpanded data structures
func Read(c a.Context, src a.Str) a.Sequence {
	l := NewLexer(src)
	return NewReader(c, l)
}

// EvalBlock evaluates data structures that have not yet been expanded
func EvalBlock(c a.Context, s a.Sequence) a.Value {
	ex := Expand(c, s).(a.Sequence)
	return a.NewBlock(ex).Eval(c)
}

// EvalStr evaluates the specified raw Source
func EvalStr(c a.Context, src a.Str) a.Value {
	r := Read(c, src)
	return EvalBlock(c, r)
}

// EvalExpand evaluates a block and then expands its macros
func EvalExpand(c a.Context, s a.Sequence) a.Value {
	ev := a.NewBlock(s).Eval(c)
	return Expand(c, ev)
}

// NewEvalContext creates a new Context instance that
// chains up to the UserDomain Context for special forms
func NewEvalContext() a.Context {
	ns := a.GetNamespace(a.UserDomain)
	c := a.ChildContext(ns)
	c.Put(a.ContextDomain, ns)
	return c
}
