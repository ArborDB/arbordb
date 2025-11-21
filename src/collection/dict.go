package collection

import (
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
	Iter(ctx *core.Context) iter.Seq2[KV[K, V], error]
}

type KV[K any, V any] struct {
	Key   K
	Value V
}
