package builtins_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	a := assert.New(t)
	l := s.NewList(big.NewFloat(1)).Cons(big.NewFloat(2)).Cons(big.NewFloat(10))
	f, _ := b.BuiltIns.Get("+")
	r := f.(*s.Function).Exec(b.BuiltIns, l)
	a.Equal(big.NewFloat(13.0), r)
}

func TestSubtract(t *testing.T) {
	a := assert.New(t)
	l := &s.Vector{big.NewFloat(7), big.NewFloat(3), big.NewFloat(1)}
	f, _ := b.BuiltIns.Get("-")
	r := f.(*s.Function).Exec(b.BuiltIns, l)
	a.Equal(big.NewFloat(3.0), r)
}

func TestMultiply(t *testing.T) {
	a := assert.New(t)
	l := s.NewList(big.NewFloat(12)).Cons(big.NewFloat(2)).Cons(big.NewFloat(5))
	f, _ := b.BuiltIns.Get("*")
	r := f.(*s.Function).Exec(b.BuiltIns, l)
	a.Equal(big.NewFloat(120), r)
}

func TestDivide(t *testing.T) {
	a := assert.New(t)
	l := &s.Vector{big.NewFloat(10), big.NewFloat(2), big.NewFloat(5)}
	f, _ := b.BuiltIns.Get("/")
	r := f.(*s.Function).Exec(b.BuiltIns, l)
	a.Equal(big.NewFloat(1.0), r)
}
