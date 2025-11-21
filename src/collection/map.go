package collection

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type Map[K interface {
	comparable
	core.Expression
}, V core.Expression] map[K]V

var _ core.Expression = Map[scalar.String, scalar.Int]{}

func (m Map[K, V]) String() string {
	return fmt.Sprintf("Map(%v)", map[K]V(m))
}

var _ Dict[scalar.String, scalar.Int] = Map[scalar.String, scalar.Int]{}

func (m Map[K, V]) Get(ctx *core.Context, key K) (value V, err error) {
	v, ok := m[key]
	if !ok {
		var zero V
		return zero, fmt.Errorf("key %v not found", key)
	}
	return v, nil
}

func (m Map[K, V]) Exists(ctx *core.Context, key K) (bool, error) {
	_, ok := m[key]
	return ok, nil
}

func (m Map[K, V]) Size(ctx *core.Context) (int, error) {
	return len(m), nil
}

func (m Map[K, V]) IterDict(ctx *core.Context) iter.Seq2[KV[K, V], error] {
	return func(yield func(KV[K, V], error) bool) {
		for k, v := range m {
			if !yield(KV[K, V]{
				Key:   k,
				Value: v,
			}, nil) {
				break
			}
		}
	}
}

var _ core.CanonicalList = Map[scalar.Int, scalar.Int]{}

func (m Map[K, V]) IterCanonical(ctx *core.Context) iter.Seq2[core.Expression, error] {
	return func(yield func(core.Expression, error) bool) {
		keys := make([]K, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		slices.SortFunc(keys, func(a, b K) int {
			return cmp.Compare(a.String(), b.String())
		})
		for _, k := range keys {
			v := m[k]
			if !yield(KV[K, V]{Key: k, Value: v}, nil) {
				return
			}
		}
	}
}
