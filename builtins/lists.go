package builtins

import a "github.com/kode4food/sputter/api"

func isListCommand(c *a.Context, args a.Iterable) a.Value {
	AssertArity(args, 1)
	i := args.Iterate()
	if v, ok := i.Next(); ok {
		if _, ok := a.Evaluate(c, v).(*a.List); ok {
			return a.True
		}
	}
	return a.False
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "list?", Exec: isListCommand})
}