package api

import u "github.com/kode4food/sputter/util"

// ExpectedGetter is thrown if the Value is not a Mapped
const ExpectedGetter = "expected a propertied value: %s"

var keywords = u.NewCache()

// Keyword is a Value that represents a Name that resolves to itself
type Keyword interface {
	Value
	Applicable
	Named
	IsKeyword() bool
}

type keyword struct {
	name Name
}

// NewKeyword returns an interned instance of a Keyword
func NewKeyword(n Name) Keyword {
	return keywords.Get(n, func() u.Any {
		return &keyword{name: n}
	}).(Keyword)
}

func (k *keyword) IsKeyword() bool {
	return true
}

func (k *keyword) Name() Name {
	return k.name
}

func (k *keyword) Apply(c Context, args Sequence) Value {
	AssertArity(args, 1)
	return k.Get(args.First())
}

func (k *keyword) Get(v Value) Value {
	if g, ok := v.(Mapped); ok {
		if r, ok := g.Get(k); ok {
			return r
		}
		panic(ErrStr(KeyNotFound, k))
	}
	panic(ErrStr(ExpectedGetter, v))
}

func (k *keyword) Str() Str {
	return ":" + Str(k.name)
}
