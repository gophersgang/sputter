package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

// Namespace is a special Namespace for built-in identifiers
var Namespace = a.GetNamespace(a.BuiltInDomain)

func registerAnnotated(v a.Annotated) {
	n := v.Metadata()[a.MetaName].(a.Name)
	Namespace.Put(n, v.(a.Value))
}

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalSequence(c, args)
}

func quote(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	//c.Put("unquote", a.NewMacro(unquote))
	return args.First()
}

func unquote(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if m, ok := v.(a.MakeExpression); ok {
		return m.Expression()
	}
	return v
}

func init() {
	registerAnnotated(
		a.NewFunction(do).WithMetadata(a.Metadata{
			a.MetaName: a.Name("do"),
			a.MetaDoc:  d.Get("do"),
		}),
	)

	registerAnnotated(
		a.NewMacro(quote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("quote"),
			a.MetaDoc:  d.Get("quote"),
		}),
	)
}
