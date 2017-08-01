package builtins

import (
	"strings"

	a "github.com/kode4food/sputter/api"
)

const (
	// InvalidRestArgument is thrown if you include more than one rest argument
	InvalidRestArgument = "rest-argument not well-formed: %s"

	// ExpectedArguments is thrown if argument patterns don't match
	ExpectedArguments = "expected arguments of the form: %s"
)

type (
	lambdaFunction struct{ a.ReflectedFunction }
	applyFunction  struct{ a.ReflectedFunction }

	argProcessor func(a.Context, a.Sequence) (a.Context, bool)

	functionSignature struct {
		args a.Vector
		body a.Sequence
	}

	functionSignatures []*functionSignature

	functionDefinition struct {
		name a.Name
		sigs functionSignatures
		meta a.Object
	}

	argProcessorMatch struct {
		args argProcessor
		body a.Block
	}
)

var (
	emptyMetadata = a.Properties{}
	restMarker    = a.Name("&")
)

func makeArgProcessor(cl a.Context, s a.Sequence) argProcessor {
	an := []a.Name{}
	for i := s; i.IsSequence(); i = i.Rest() {
		n := a.AssertUnqualified(i.First()).Name()
		if n == restMarker {
			rn := parseRestArg(i)
			return makeRestArgProcessor(cl, an, rn)
		}
		an = append(an, n)
	}
	return makeFixedArgProcessor(cl, an)
}

func parseRestArg(s a.Sequence) a.Name {
	r := s.Rest()
	if r.IsSequence() {
		n := a.AssertUnqualified(r.First()).Name()
		if n != restMarker && !r.Rest().IsSequence() {
			return n
		}
	}
	panic(a.ErrStr(InvalidRestArgument, s))
}

func makeRestArgProcessor(cl a.Context, an []a.Name, rn a.Name) argProcessor {
	ac := a.MakeMinimumArityChecker(len(an))

	return func(_ a.Context, args a.Sequence) (a.Context, bool) {
		if _, ok := ac(args); !ok {
			return nil, false
		}
		l := a.ChildContext(cl)
		i := args
		for _, n := range an {
			l.Put(n, i.First())
			i = i.Rest()
		}
		l.Put(rn, a.SequenceToList(i))
		return l, true
	}
}

func makeFixedArgProcessor(cl a.Context, an []a.Name) argProcessor {
	ac := a.MakeArityChecker(len(an))

	return func(_ a.Context, args a.Sequence) (a.Context, bool) {
		if _, ok := ac(args); !ok {
			return nil, false
		}
		l := a.ChildContext(cl)
		i := args
		for _, n := range an {
			l.Put(n, i.First())
			i = i.Rest()
		}
		return l, true
	}
}

func optionalMetadata(args a.Sequence) (a.Object, a.Sequence) {
	r := args
	var md a.Object
	if s, ok := r.First().(a.Str); ok {
		md = a.Properties{a.DocKey: s}
		r = r.Rest()
	} else {
		md = emptyMetadata
	}

	if m, ok := r.First().(a.MappedSequence); ok {
		md = md.Child(toProperties(m))
		r = r.Rest()
	}
	return md, r
}

func optionalName(args a.Sequence) (a.Name, a.Sequence) {
	f := args.First()
	if s, ok := f.(a.Symbol); ok {
		if s.Domain() == a.LocalDomain {
			return s.Name(), args.Rest()
		}
		panic(a.ErrStr(a.ExpectedUnqualified, s.Qualified()))
	}
	return a.DefaultFunctionName, args
}

func parseNamedFunction(args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	fn := a.AssertUnqualified(args.First()).Name()
	return parseFunctionRest(fn, args.Rest())
}

func parseFunction(args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	return parseFunctionRest(fn, r)
}

func parseFunctionRest(fn a.Name, r a.Sequence) *functionDefinition {
	md, r := optionalMetadata(r)
	sigs := parseFunctionSignatures(r)
	md = md.Child(a.Properties{
		a.NameKey: fn,
	})

	return &functionDefinition{
		name: fn,
		sigs: sigs,
		meta: md,
	}
}

func parseFunctionSignatures(r a.Sequence) functionSignatures {
	if args, ok := r.First().(a.Vector); ok {
		return functionSignatures{
			{args: args, body: r.Rest()},
		}
	}
	res := functionSignatures{}
	for i := r; i.IsSequence(); i = i.Rest() {
		l := a.AssertList(i.First())
		res = append(res, &functionSignature{
			args: a.AssertVector(l.First()),
			body: l.Rest(),
		})
	}
	return res
}

func makeFunction(c a.Context, d *functionDefinition) a.Function {
	var res a.Function

	if len(d.sigs) > 1 {
		res = makeMultiFunction(c, d.sigs)
	} else {
		res = makeSingleFunction(c, d.sigs[0])
	}
	return res.WithMetadata(d.meta).(a.Function)
}

func makeSingleFunction(c a.Context, s *functionSignature) a.Function {
	ap := makeArgProcessor(c, s.args)
	ex := a.MacroExpandAll(c, s.body).(a.Sequence)
	db := a.NewBlock(ex)

	return a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		if l, ok := ap(c, args); ok {
			return a.Eval(l, db)
		}
		panic(a.ErrStr(ExpectedArguments, s.args))
	})
}

func makeMultiFunction(c a.Context, sigs functionSignatures) a.Function {
	ls := len(sigs)
	procMap := make([]argProcessorMatch, ls)
	args := make([]string, ls)

	for i, s := range sigs {
		ex := a.MacroExpandAll(c, s.body).(a.Sequence)
		procMap[i] = argProcessorMatch{
			args: makeArgProcessor(c, s.args),
			body: a.NewBlock(ex),
		}
		args[i] = string(s.args.Str())
	}

	argPatterns := strings.Join(args, " or ")

	return a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		for _, m := range procMap {
			if l, ok := m.args(c, args); ok {
				return a.Eval(l, m.body)
			}
		}

		panic(a.ErrStr(ExpectedArguments, argPatterns))
	})
}

func (f *lambdaFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fd := parseFunction(args)
	return makeFunction(c, fd)
}

func (f *applyFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	fn := a.AssertApplicable(args.First())
	s := a.AssertSequence(args.Rest().First())
	return fn.Apply(c, s)
}

func isApplicable(v a.Value) bool {
	if _, ok := v.(a.Applicable); ok {
		return true
	}
	return false
}

func isSpecialForm(v a.Value) bool {
	if ap, ok := v.(a.Applicable); ok {
		return a.IsSpecialForm(ap)
	}
	return false
}

func init() {
	var lambda *lambdaFunction
	var apply *applyFunction

	RegisterBaseFunction("lambda", lambda)
	RegisterBaseFunction("apply", apply)

	RegisterSequencePredicate("apply?", isApplicable)
	RegisterSequencePredicate("special-form?", isSpecialForm)
}
