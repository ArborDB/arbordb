package collection

import (
	"iter"

	"github.com/ArborDB/arbordb/src/core"
)

type Dict[K interface {
	comparable
	core.Expression
}, V core.Expression] interface {
	core.Expression
	Get(ctx *core.Context, key K) (value V, err error)
	Exists(ctx *core.Context, key K) (bool, error)
	Size(ctx *core.Context) (int, error)
	IterDict(ctx *core.Context) iter.Seq2[KV[K, V], error]
}
