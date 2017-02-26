package api

import "fmt"

// NonFinite is thrown if count is called against a non-finite sequence
const NonFinite = "sequence is not finite and can't be counted"

var (
	// True is a value that represents any value other than False
	True = &Atom{Label: "true"}

	// False is a value that represents either itself or nil
	False = &Atom{Label: "false"}
)

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Variables are represents a mapping from Name to Value
type Variables map[Name]Value

// Sequence interfaces expose a one dimensional set of Values
type Sequence interface {
	Iterate() Iterator
}

// Finite interfaces allow a Sequence item to be retrieved by index
type Finite interface {
	Count() int
	Get(index int) Value
}

// Iterator interfaces are stateful iteration interfaces
type Iterator interface {
	Next() (Value, bool)
	Rest() Sequence
}

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	switch {
	case v == Nil || v == False || v == nil || v == false:
		return false
	default:
		return true
	}
}

// Count will either use Finite.Count() or iterate over the Sequence
func Count(s Sequence) int {
	if f, ok := s.(Finite); ok {
		return f.Count()
	}
	panic(NonFinite)
}

// String either calls the String() method or tries to convert
func String(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return v.(string)
}
