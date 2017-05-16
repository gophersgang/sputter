package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func vector(c a.Context, args a.Sequence) a.Value {
	if cnt, ok := args.(a.Counted); ok {
		l := cnt.Count()
		r := make([]a.Value, l)
		idx := 0
		for i := args; i.IsSequence(); i = i.Rest() {
			r[idx] = i.First().Eval(c)
			idx++
		}
		return a.NewVector(r...)
	}
	return vectorFromUncounted(c, args)
}

func vectorFromUncounted(c a.Context, args a.Sequence) a.Value {
	r := []a.Value{}
	for i := args; i.IsSequence(); i = i.Rest() {
		r = append(r, i.First().Eval(c))
	}
	return a.NewVector(r...)
}

func toVector(c a.Context, args a.Sequence) a.Value {
	return vector(c, concat(c, args).(a.Sequence))
}

func isVector(v a.Value) bool {
	if _, ok := v.(a.Vector); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(vector).WithMetadata(a.Metadata{
			a.MetaName: a.Name("vector"),
			a.MetaDoc:  d.Get("vector"),
		}),
	)

	registerAnnotated(
		a.NewFunction(toVector).WithMetadata(a.Metadata{
			a.MetaName: a.Name("to-vector"),
			a.MetaDoc:  d.Get("to-vector"),
		}),
	)

	registerSequencePredicate(isVector, a.Metadata{
		a.MetaName: a.Name("vector?"),
		a.MetaDoc:  d.Get("is-vector"),
	})
}
