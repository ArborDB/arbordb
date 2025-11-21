package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
)

type Dict[K interface {
	comparable
	core.Expression
}, V core.Expression] interface {
	Get(ctx *core.Context, key K) (value V, err error)
	Exists(ctx *core.Context, key K) (bool, error)
	Size(ctx *core.Context) (int, error)
	IterDict(ctx *core.Context) iter.Seq2[KV[K, V], error]
}

type KV[K core.Expression, V core.Expression] struct {
	Key   K
	Value V
}

var _ core.Expression = KV[core.Expression, core.Expression]{}

func (kv KV[K, V]) String() string {
	return fmt.Sprintf("(%v: %v)", kv.Key, kv.Value)
}
