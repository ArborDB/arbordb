package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type DictSet[K interface {
	comparable
	core.Expression
}, V core.Expression] struct {
	Dict  Dict[K, V]
	Key   K
	Value V
}

var _ core.Expression = DictSet[scalar.String, scalar.Int]{}

func (ds DictSet[K, V]) String() string {
	return fmt.Sprintf("DictSet(%v, %v, %v)", ds.Dict, ds.Key, ds.Value)
}

var _ Dict[scalar.String, scalar.Int] = DictSet[scalar.String, scalar.Int]{}

func (ds DictSet[K, V]) Get(ctx *core.Context, key K) (value V, err error) {
	if key == ds.Key {
		return ds.Value, nil
	}
	return ds.Dict.Get(ctx, key)
}

func (ds DictSet[K, V]) Exists(ctx *core.Context, key K) (bool, error) {
	if key == ds.Key {
		return true, nil
	}
	return ds.Dict.Exists(ctx, key)
}

func (ds DictSet[K, V]) Size(ctx *core.Context) (int, error) {
	exists, err := ds.Dict.Exists(ctx, ds.Key)
	if err != nil {
		return 0, err
	}
	size, err := ds.Dict.Size(ctx)
	if err != nil {
		return 0, err
	}
	if exists {
		return size, nil
	}
	return size + 1, nil
}

func (ds DictSet[K, V]) IterDict(ctx *core.Context) iter.Seq2[KV[K, V], error] {
	return func(yield func(KV[K, V], error) bool) {
		yieldedSetKey := false
		for kv, err := range ds.Dict.IterDict(ctx) {
			if err != nil {
				yield(KV[K, V]{}, err)
				return
			}
			if kv.Key == ds.Key {
				if !yield(KV[K, V]{Key: ds.Key, Value: ds.Value}, nil) {
					return
				}
				yieldedSetKey = true
			} else {
				if !yield(kv, nil) {
					return
				}
			}
		}
		if !yieldedSetKey {
			if !yield(KV[K, V]{Key: ds.Key, Value: ds.Value}, nil) {
				return
			}
		}
	}
}
