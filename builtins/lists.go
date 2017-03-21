package builtins

import (
	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

func list(c a.Context, args a.Sequence) a.Value {
	s := u.NewStack()
	for i := args; i.IsSequence(); i = i.Rest() {
		s.Push(a.Eval(c, i.First()))
	}

	e, ok := s.Pop()
	if !ok {
		return a.EmptyList
	}

	var l = a.NewList(e)
	for v, ok := s.Pop(); ok; v, ok = s.Pop() {
		l = l.Prepend(v)
	}
	return l
}

func toList(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	arg := a.Eval(c, args.First())
	seq := a.AssertSequence(arg)
	return list(c, seq)
}

func isList(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if _, ok := a.Eval(c, v).(*a.List); ok {
		return a.True
	}
	return a.False
}

func init() {
	registerFunction(&a.Function{Name: "list", Exec: list})
	registerFunction(&a.Function{Name: "to-list", Exec: toList})
	registerPredicate(&a.Function{Name: "list?", Exec: isList})
}
