package builtins

import a "github.com/kode4food/sputter/api"

type (
	chanFunction    struct{ BaseBuiltIn }
	promiseFunction struct{ BaseBuiltIn }
	goFunction      struct{ BaseBuiltIn }
)

var (
	// MetaChannel is the key used to identify a Channel
	MetaChannel = a.NewKeyword("channel")

	// MetaEmitter is the key used to emit to a Channel
	MetaEmitter = a.NewKeyword("emit")

	// MetaSequence is the key used to retrieve the Sequence from a Channel
	MetaSequence = a.NewKeyword("seq")
)

var channelPrototype = a.Properties{
	MetaChannel: a.True,
}

func (*chanFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	e, s := a.NewChannel()

	return channelPrototype.Child(a.Properties{
		MetaEmitter:  bindWriter(e),
		MetaClose:    bindCloser(e),
		MetaSequence: s,
	})
}

func (*promiseFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if a.AssertArityRange(args, 0, 1) == 1 {
		p := a.NewPromise()
		p.Deliver(args.First())
		return p
	}
	return a.NewPromise()
}

func isPromise(v a.Value) bool {
	if _, ok := v.(a.Promise); ok {
		return true
	}
	return false
}

func (*goFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	go a.MakeBlock(args).Eval(a.ChildContext(c))
	return a.Nil
}

func init() {
	var _chan *chanFunction
	var promise *promiseFunction
	var _go *goFunction

	RegisterBuiltIn("chan", _chan)
	RegisterBuiltIn("promise", promise)
	RegisterBuiltIn("make-go", _go)

	RegisterSequencePredicate("promise?", isPromise)
}
