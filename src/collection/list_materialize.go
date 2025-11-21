package collection

import (
	"slices"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type ListToArray[E core.Expression] struct {
}

var _ core.Transform[List[scalar.Int], Array[scalar.Int]] = ListToArray[scalar.Int]{}

func (l ListToArray[E]) Apply(ctx *core.Context, from List[E], to *Array[E]) (steps core.TransformSteps) {
	return func(yield func(*core.TransformStep, error) bool) {
		length, err := from.Length(ctx)
		if err != nil {
			yield(nil, err)
			return
		}
		*to = (*to)[:0]
		if cap(*to) < length {
			*to = slices.Grow(*to, length-cap(*to))
		}
		for elem, err := range from.Iter(ctx) {
			if err != nil {
				yield(nil, err)
				return
			}
			*to = append(*to, elem)
		}
	}
}

func (l ListToArray[E]) EstimateCost(ctx *core.Context, from List[E]) (core.Cost, error) {
	return core.Cost{}, nil
}
