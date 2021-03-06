package builtins

import a "github.com/kode4food/sputter/api"

type (
	lazySequenceFunction struct{ BaseBuiltIn }
	concatFunction       struct{ BaseBuiltIn }
	filterFunction       struct{ BaseBuiltIn }
	mapFunction          struct{ BaseBuiltIn }
	reduceFunction       struct{ BaseBuiltIn }
	takeFunction         struct{ BaseBuiltIn }
	dropFunction         struct{ BaseBuiltIn }
	rangeFunction        struct{ BaseBuiltIn }
	forEachFunction      struct{ BaseBuiltIn }
)

func makeLazyResolver(c a.Context, f a.Applicable) a.LazyResolver {
	return func() (a.Value, a.Sequence, bool) {
		r := f.Apply(c, a.EmptyList)
		if s, ok := r.(a.Sequence); ok {
			if f, r, ok := s.Split(); ok {
				return f, r, true
			}
		}
		if r == a.Nil {
			return a.Nil, a.EmptyList, false
		}
		panic(a.ErrStr(a.ExpectedSequence, r))
	}
}

func makeValueFilter(c a.Context, f a.Applicable) a.ValueFilter {
	return func(v a.Value) bool {
		return a.Truthy(f.Apply(c, a.NewVector(v)))
	}
}

func makeValueMapper(c a.Context, f a.Applicable) a.ValueMapper {
	return func(v a.Value) a.Value {
		return f.Apply(c, a.NewVector(v))
	}
}

func makeValueReducer(c a.Context, f a.Applicable) a.ValueReducer {
	return func(l, r a.Value) a.Value {
		return f.Apply(c, a.NewVector(l, r))
	}
}

func (*lazySequenceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fn := NewBlockFunction(args)
	return a.NewLazySequence(makeLazyResolver(c, fn))
}

func (*concatFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if a.AssertMinimumArity(args, 1) == 1 {
		return a.AssertSequence(args.First())
	}
	return a.Concat(args)
}

func (*filterFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	fn := a.AssertApplicable(f)
	s := a.Concat(r)
	return a.Filter(s, makeValueFilter(c, fn))
}

func (*mapFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	fn := a.AssertApplicable(f)
	s := a.Concat(r)
	return a.Map(s, makeValueMapper(c, fn))
}

func (*reduceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	fn := a.AssertApplicable(f)
	s := a.Concat(r)
	return a.Reduce(s, makeValueReducer(c, fn))
}

func (*takeFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	n := a.AssertInteger(f)
	s := a.Concat(r)
	return a.Take(s, n)
}

func (*dropFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	n := a.AssertInteger(f)
	s := a.Concat(r)
	return a.Drop(s, n)
}

func (*rangeFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 3)
	f, r, _ := args.Split()
	low := a.AssertNumber(f)
	rf, rr, _ := r.Split()
	high := a.AssertNumber(rf)
	step := a.AssertNumber(rr.First())
	return a.NewRange(low, high, step)
}

func (*forEachFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	b := a.AssertVector(f)
	bc := b.Count()
	if bc%2 != 0 {
		panic(a.ErrStr(ExpectedBindings))
	}

	var proc forProc
	depth := bc / 2
	for i := depth - 1; i >= 0; i-- {
		o := i * 2
		s, _ := b.ElementAt(o)
		e, _ := b.ElementAt(o + 1)
		n := a.AssertUnqualified(s).Name()
		if i == depth-1 {
			proc = makeTerminal(n, e, r)
		} else {
			proc = makeIntermediate(n, e, proc)
		}
	}

	proc(c)
	return a.Nil
}

func makeIntermediate(n a.Name, e a.Value, next forProc) forProc {
	return func(c a.Context) {
		s := a.AssertSequence(a.Eval(c, e))
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			l := a.ChildContext(c)
			l.Put(n, f)
			next(l)
		}
	}
}

func makeTerminal(n a.Name, e a.Value, s a.Sequence) forProc {
	bl := a.MakeBlock(s)
	return func(c a.Context) {
		s := a.AssertSequence(a.Eval(c, e))
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			l := a.ChildContext(c)
			l.Put(n, f)
			bl.Eval(l)
		}
	}
}

func init() {
	var lazySequence *lazySequenceFunction
	var concat *concatFunction
	var filter *filterFunction
	var _map *mapFunction
	var reduce *reduceFunction
	var take *takeFunction
	var drop *dropFunction
	var _range *rangeFunction
	var forEach *forEachFunction

	RegisterBuiltIn("make-lazy-seq", lazySequence)
	RegisterBuiltIn("concat", concat)
	RegisterBuiltIn("filter", filter)
	RegisterBuiltIn("map", _map)
	RegisterBuiltIn("reduce", reduce)
	RegisterBuiltIn("take", take)
	RegisterBuiltIn("drop", drop)
	RegisterBuiltIn("make-range", _range)
	RegisterBuiltIn("for-each", forEach)
}
