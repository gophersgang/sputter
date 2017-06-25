package core

import (
	"strings"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assets"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

const prefix = "core/"

func init() {
	for _, name := range assets.AssetNames() {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		src := a.Str(assets.MustGet(name))
		e.EvalStr(b.Namespace, src)
	}
}
