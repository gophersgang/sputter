package api

// Value is the generic interface for all 'Values' in the VM
type Value interface {
}

// Iterable values can be used in loops and comprehensions
type Iterable interface {
	Iterate() Iterator
}

// Iterator interfaces are stateful iteration interfaces
type Iterator interface {
	Next() (Value, bool)
	Iterable() Iterable
}