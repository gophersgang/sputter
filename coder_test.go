package sputter

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCoder(t *testing.T) {
	a := assert.New(t)
	l := NewLexer("99")
	c := NewCoder(l)
	a.NotNil(c)
}

func TestCodeInteger(t *testing.T) {
	a := assert.New(t)
	l := NewLexer("99")
	c := NewCoder(l)
	v := c.Next()
	f, ok := v.(*big.Float)
	a.True(ok)
	a.Equal(0, f.Cmp(big.NewFloat(99)))
}

func TestCodeList(t *testing.T) {
	a := assert.New(t)
	l := NewLexer(`(99 "hello" 55.12)`)
	c := NewCoder(l)
	v := c.Next()
	list, ok := v.(*List)
	a.True(ok)

	next := list.Iterate()
	value, found := next()
	a.True(found)
	a.Equal(0, big.NewFloat(99).Cmp(value.(*big.Float)))

	value, found = next()
	a.True(found)
	a.Equal("hello", value)

	value, found = next()
	a.True(found)
	a.Equal(0, big.NewFloat(55.12).Cmp(value.(*big.Float)))

	value, found = next()
	a.False(found)
}

func TestCodeNestedList(t *testing.T) {
	a := assert.New(t)
	l := NewLexer(`(99 ("hello" "there") 55.12)`)
	c := NewCoder(l)
	v := c.Next()
	list, ok := v.(*List)
	a.True(ok)

	next1 := list.Iterate()
	value, found := next1()
	a.True(found)
	a.Equal(0, big.NewFloat(99).Cmp(value.(*big.Float)))

	// get nested list
	value, found = next1()
	a.True(found)
	list2, ok := value.(*List)
	a.True(ok)

	// iterate over the rest of top-level list
	value, found = next1()
	a.True(found)
	a.Equal(0, big.NewFloat(55.12).Cmp(value.(*big.Float)))

	value, found = next1()
	a.False(found)

	// iterate over the nested list
	next2 := list2.Iterate()
	value, found = next2()
	a.True(found)
	a.Equal("hello", value)

	value, found = next2()
	a.True(found)
	a.Equal("there", value)

	value, found = next2()
	a.False(found)
}

func TestUnclosedList(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			a.Equal(r, ListNotClosed, "unclosed list")
			return
		}
		a.Fail("unclosed list didn't panic")
	}()

	l := NewLexer(`(99 ("hello" "there") 55.12`)
	c := NewCoder(l)
	c.Next()
}

func TestLiteral(t *testing.T) {
	a := assert.New(t)

	l := NewLexer(`'99`)
	c := NewCoder(l)
	v := c.Next()

	literal, ok := v.(*Literal)
	a.True(ok)

	value, ok := literal.value.(*big.Float)
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value))
}
