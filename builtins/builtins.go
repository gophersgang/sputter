package builtins

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
)

type (
	// BuiltInFunction is an interface that identifies a built-in
	BuiltInFunction interface {
		a.Function
		IsBuiltIn() bool
	}

	// BaseBuiltIn is the base structure for built-in Functions
	BaseBuiltIn struct {
		a.BaseFunction
		concrete reflect.Type
	}
)

var (
	// Namespace is a special Namespace for built-in identifiers
	Namespace = a.GetNamespace(a.BuiltInDomain)

	builtInFuncs = map[a.Name]a.Function{}

	baseZeroValue BaseBuiltIn
	baseFieldName = reflect.TypeOf(baseZeroValue).Name()
)

// MakeBuiltIn uses reflection to instantiate a built-in Function
func MakeBuiltIn(f BuiltInFunction) a.Function {
	t := reflect.TypeOf(f).Elem()
	return newBuiltInWithBase(t, a.DefaultBaseFunction)
}

// IsBuiltIn returns whether or not this Function is a built-in
func (f *BaseBuiltIn) IsBuiltIn() bool {
	return true
}

// WithMetadata creates a copy of this Function with additional Metadata
func (f *BaseBuiltIn) WithMetadata(md a.Object) a.AnnotatedValue {
	b := f.Extend(md)
	return newBuiltInWithBase(f.concrete, b)
}

func newBuiltInWithBase(t reflect.Type, b a.BaseFunction) a.Function {
	ptr := reflect.New(t)
	v := reflect.Indirect(ptr)
	f := reflect.Indirect(v).FieldByName(baseFieldName)
	f.Set(reflect.ValueOf(BaseBuiltIn{
		BaseFunction: b,
		concrete:     t,
	}))
	return ptr.Interface().(a.Function)
}

// RegisterFunction registers a built-in Function by Name
func RegisterFunction(n a.Name, f a.Function) {
	builtInFuncs[n] = f
}

// RegisterBuiltIn registers a Base-derived Function by Name
func RegisterBuiltIn(n a.Name, f BuiltInFunction) {
	RegisterFunction(n, MakeBuiltIn(f))
}

// GetFunction returns a registered built-in function
func GetFunction(n a.Name) (a.Function, bool) {
	if f, ok := builtInFuncs[n]; ok {
		return f, true
	}
	return nil, false
}

func isBuiltInDomain(s a.Symbol) bool {
	return s.Domain() == a.BuiltInDomain
}

func isBuiltInCall(n a.Name, v a.Value) (a.List, bool) {
	if l, ok := v.(a.List); ok && l.Count() > 0 {
		if s, ok := l.First().(a.Symbol); ok {
			return l, isBuiltInDomain(s) && s.Name() == n
		}
	}
	return nil, false
}

func defBuiltIn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	f, r, _ := args.Split()
	n := a.AssertUnqualified(f).Name()
	if f, ok := GetFunction(n); ok {
		var md a.Object = toProperties(a.SequenceToAssociative(r))
		r := f.WithMetadata(md)
		a.GetContextNamespace(c).Put(n, r)
		return r
	}
	panic(a.ErrStr(a.KeyNotFound, n))
}

func init() {
	Namespace.Put("def-builtin",
		a.NewExecFunction(defBuiltIn).WithMetadata(a.Properties{
			a.SpecialKey: a.True,
		}),
	)
}
