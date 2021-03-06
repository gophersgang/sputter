package builtins

import a "github.com/kode4food/sputter/api"

type (
	reduceFunc  func(prev a.Number, next a.Number) a.Number
	compareFunc func(prev a.Number, next a.Number) bool

	addFunction struct{ BaseBuiltIn }
	subFunction struct{ BaseBuiltIn }
	mulFunction struct{ BaseBuiltIn }
	divFunction struct{ BaseBuiltIn }
	modFunction struct{ BaseBuiltIn }
	eqFunction  struct{ BaseBuiltIn }
	neqFunction struct{ BaseBuiltIn }
	gtFunction  struct{ BaseBuiltIn }
	gteFunction struct{ BaseBuiltIn }
	ltFunction  struct{ BaseBuiltIn }
	lteFunction struct{ BaseBuiltIn }
)

func reduceNum(s a.Sequence, v a.Number, fn reduceFunc) a.Value {
	res := v
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		fv := a.AssertNumber(f)
		res = fn(res, fv)
	}
	return res
}

func fetchFirstNumber(args a.Sequence) (a.Number, a.Sequence) {
	a.AssertMinimumArity(args, 1)
	f, r, _ := args.Split()
	nv := a.AssertNumber(f)
	return nv, r
}

func compare(_ a.Context, s a.Sequence, fn compareFunc) a.Bool {
	cur, rest := fetchFirstNumber(s)
	for f, r, ok := rest.Split(); ok; f, r, ok = r.Split() {
		v := a.AssertNumber(f)
		if !fn(cur, v) {
			return a.False
		}
		cur = v
	}
	return a.True
}

func (*addFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return a.Zero
	}
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Add(n)
	})
}

func (*subFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Sub(n)
	})
}

func (*mulFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if !args.IsSequence() {
		return a.One
	}
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Mul(n)
	})
}

func (*divFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Div(n)
	})
}

func (*modFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	v, r := fetchFirstNumber(args)
	return reduceNum(r, v, func(p a.Number, n a.Number) a.Number {
		return p.Mod(n)
	})
}

func (*eqFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	})
}

func (*neqFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.EqualTo
	}).Not()
}

func (*gtFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.GreaterThan
	})
}

func (*gteFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.GreaterThan
	})
}

func (*ltFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		return p.Cmp(n) == a.LessThan
	})
}

func (*lteFunction) Apply(c a.Context, args a.Sequence) a.Value {
	return compare(c, args, func(p a.Number, n a.Number) bool {
		r := p.Cmp(n)
		return r == a.EqualTo || r == a.LessThan
	})
}

func isPosInfinity(v a.Value) bool {
	if n, ok := v.(a.Number); ok {
		return a.PosInfinity.Cmp(n) == a.EqualTo
	}
	return false
}

func isNegInfinity(v a.Value) bool {
	if n, ok := v.(a.Number); ok {
		return a.NegInfinity.Cmp(n) == a.EqualTo
	}
	return false
}

func isNaN(v a.Value) bool {
	if n, ok := v.(a.Number); ok {
		return n.IsNaN()
	}
	return true
}

func init() {
	var add *addFunction
	var sub *subFunction
	var mul *mulFunction
	var div *divFunction
	var mod *modFunction
	var eq *eqFunction
	var neq *neqFunction
	var gt *gtFunction
	var gte *gteFunction
	var lt *ltFunction
	var lte *lteFunction

	Namespace.Put("inf", a.PosInfinity)
	Namespace.Put("-inf", a.NegInfinity)

	RegisterBuiltIn("+", add)
	RegisterBuiltIn("-", sub)
	RegisterBuiltIn("*", mul)
	RegisterBuiltIn("/", div)
	RegisterBuiltIn("%", mod)
	RegisterBuiltIn("=", eq)
	RegisterBuiltIn("!=", neq)
	RegisterBuiltIn(">", gt)
	RegisterBuiltIn(">=", gte)
	RegisterBuiltIn("<", lt)
	RegisterBuiltIn("<=", lte)

	RegisterSequencePredicate("inf?", isPosInfinity)
	RegisterSequencePredicate("-inf?", isNegInfinity)
	RegisterSequencePredicate("nan?", isNaN)
}
