package builtins

import a "github.com/kode4food/sputter/api"

type assocFunction struct{ BaseBuiltIn }

func (*assocFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToAssociative(args)
}

func isAssociative(v a.Value) bool {
	if _, ok := v.(a.Associative); ok {
		return true
	}
	return false
}

func isMapped(v a.Value) bool {
	if _, ok := v.(a.MappedSequence); ok {
		return true
	}
	return false
}

func init() {
	var assoc *assocFunction

	RegisterBuiltIn("assoc", assoc)
	RegisterSequencePredicate("assoc?", isAssociative)
	RegisterSequencePredicate("mapped?", isMapped)
}
