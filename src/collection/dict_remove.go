package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type DictRemove[K interface {
	comparable
	core.Expression
}, V core.Expression] struct {
	Dict Dict[K, V]
	Key  K
}

var _ core.Expression = DictRemove[scalar.String, scalar.Int]{}

func (dr DictRemove[K, V]) String() string {
	return fmt.Sprintf("DictRemove(%v, %v)", dr.Dict, dr.Key)
}

var _ Dict[scalar.String, scalar.Int] = DictRemove[scalar.String, scalar.Int]{}

func (dr DictRemove[K, V]) Get(ctx *core.Context, key K) (value V, err error) {
	if key == dr.Key {
		var zero V
		return zero, fmt.Errorf("key %v not found", key)
	}
	return dr.Dict.Get(ctx, key)
}

func (dr DictRemove[K, V]) Exists(ctx *core.Context, key K) (bool, error) {
	if key == dr.Key {
		return false, nil
	}
	return dr.Dict.Exists(ctx, key)
}

func (dr DictRemove[K, V]) Size(ctx *core.Context) (int, error) {
	exists, err := dr.Dict.Exists(ctx, dr.Key)
	if err != nil {
		return 0, err
	}
	size, err := dr.Dict.Size(ctx)
	if err != nil {
		return 0, err
	}
	if exists {
		return size - 1, nil
	}
	return size, nil
}

func (dr DictRemove[K, V]) IterDict(ctx *core.Context) iter.Seq2[KV[K, V], error] {
	return func(yield func(KV[K, V], error) bool) {
		for kv, err := range dr.Dict.IterDict(ctx) {
			if err != nil {
				yield(KV[K, V]{}, err)
				return
			}
			if kv.Key == dr.Key {
				continue
			}
			if !yield(kv, nil) {
				return
			}
		}
	}
}
