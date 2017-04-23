package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestConstStrings(t *testing.T) {
	as := assert.New(t)

	as.String("true", a.True)
	as.String("false", a.False)
	as.String("nil", a.Nil)
}

func TestStr(t *testing.T) {
	as := assert.New(t)

	s1 := s("hello")
	as.Equal(5, a.Count(s1))
	as.String("h", s1.First())
	as.String("ello", s1.Rest())

	as.True(s1.IsSequence())
	as.False(s("").IsSequence())

	s2 := s1.Prepend(s("s"))
	as.Equal(6, a.Count(s2))
	as.String("shello", s2)

	s3 := s1.Conjoin(s("z"))
	as.Equal(6, a.Count(s3))
	as.String("helloz", s3)

	l1 := s1.Prepend(f(99))
	as.Equal(6, a.Count(l1))
	as.Equal(`(99 "h" "e" "l" "l" "o")`, l1)

	v1 := s1.Conjoin(f(99))
	as.Equal(6, a.Count(v1))
	as.Equal(`["h" "e" "l" "l" "o" 99]`, v1)

	s4 := s("thér\\再e")
	as.Equal(7, a.Count(s4))

	s5 := string(s4.Str())
	r1 := []rune(s5)
	as.Equal(10, len(r1))
	as.Equal(`"`, string(r1[0]))

	c, ok := s1.Get(1)
	as.True(ok)
	as.String("e", c)

	c, ok = s1.Get(5)
	as.False(ok)
	as.Nil(c)

	as.String("e", s1.Apply(a.NewContext(), a.Vector{f(1)}))

	s6 := s("再见!")
	as.Equal(3, a.Count(s6))
	as.String("再", s6.First())
	as.String("见!", s6.Rest())
}