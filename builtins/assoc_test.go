package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func TestAssoc(t *testing.T) {
	as := assert.New(t)
	c := e.NewEvalContext()

	assoc := getBuiltIn("assoc")
	a1 := assoc.Apply(c, args(kw("hello"), s("foo")))
	m1 := a.AssertMapped(a1)
	v1, ok := m1.Get(kw("hello"))
	as.True(ok)
	as.String("foo", v1)

	isAssoc := getBuiltIn("assoc?")
	as.True(isAssoc.Apply(c, args(a1)))
	as.False(isAssoc.Apply(c, args(f(99))))

	isMapped := getBuiltIn("mapped?")
	as.True(isMapped.Apply(c, args(a1)))
	as.False(isMapped.Apply(c, args(f(99))))
}
