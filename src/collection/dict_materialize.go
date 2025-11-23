package collection

import (
	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type DictToMap[K interface {
	comparable
	core.Expression
}, V core.Expression] struct {
}

var _ core.Transform[Dict[scalar.String, scalar.Int], Map[scalar.String, scalar.Int]] = DictToMap[scalar.String, scalar.Int]{}

func (d DictToMap[K, V]) Apply(ctx *core.Context, from Dict[K, V], to *Map[K, V]) error {
	size, err := from.Size(ctx)
	if err != nil {
		return err
	}
	newMap := make(Map[K, V], size)
	for kv, err := range from.IterDict(ctx) {
		if err != nil {
			return err
		}
		newMap[kv.Key] = kv.Value
	}
	*to = newMap
	return nil
}

func (d DictToMap[K, V]) EstimateCost(ctx *core.Context, from Dict[K, V]) (core.Cost, error) {
	return core.Cost{}, nil
}
