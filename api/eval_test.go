package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

var helloName = a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
	i := a.Iterate(args)
	n, _ := i.Next()
	v := a.Eval(c, n)
	return s("Hello, " + string(v.(a.Str)) + "!")
}).WithMetadata(a.Metadata{
	a.MetaName: a.Name("hello"),
}).(a.Function)

func TestEvaluate(t *testing.T) {
	as := assert.New(t)

	l := a.NewList(helloName, s("World"))
	c := a.NewContext()
	r := a.Eval(c, l)

	as.String("Hello, World!", r)
}

func TestBlock(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put("hello", helloName)
	helloSym := a.NewLocalSymbol("hello")

	s1 := a.NewList(helloSym, s("World"))
	s2 := a.NewList(helloSym, s("Foo"))
	l := a.NewBlock(a.NewList(s1, s2))
	as.True(l.IsBlock())
	as.String(`(hello "World")(hello "Foo")`, l)

	r := a.Eval(c, l)
	as.String("Hello, Foo!", r)
}

func TestEvaluateBlock(t *testing.T) {
	as := assert.New(t)

	s1 := a.NewList(helloName, s("World"))
	s2 := a.NewList(helloName, s("Foo"))
	l := a.NewList(s1, s2)

	c := a.NewContext()
	r := a.EvalBlock(c, l)
	as.String("Hello, Foo!", r)
}

func TestAssertApplicable(t *testing.T) {
	as := assert.New(t)
	a.AssertApplicable(a.NewFunction(nil))

	defer as.ExpectError(a.Err(a.ExpectedApplicable, f(99)))
	a.AssertApplicable(f(99))
}
