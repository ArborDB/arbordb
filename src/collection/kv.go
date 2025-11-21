package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
)

type KV[K core.Expression, V core.Expression] struct {
	Key   K
	Value V
}

var _ core.Expression = KV[core.Expression, core.Expression]{}

func (kv KV[K, V]) String() string {
	return fmt.Sprintf("(%v: %v)", kv.Key, kv.Value)
}

var _ core.CanonicalList = KV[core.Expression, core.Expression]{}

func (kv KV[K, V]) IterCanonical(ctx *core.Context) iter.Seq2[core.Expression, error] {
	return func(yield func(core.Expression, error) bool) {
		if !yield(kv.Key, nil) {
			return
		}
		yield(kv.Value, nil)
	}
}
