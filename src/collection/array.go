package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type Array[T core.Expression] []T

var _ core.Expression = Array[scalar.Int]{}

func (a Array[T]) String() string {
	return fmt.Sprintf(`Array(%v)`, []T(a))
}

var _ List[scalar.Int] = Array[scalar.Int]{}

func (a Array[T]) IsEmpty(ctx *core.Context) (bool, error) {
	return len(a) == 0, nil
}

func (a Array[T]) Iter(ctx *core.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for _, elem := range a {
			if !yield(elem, nil) {
				break
			}
		}
	}
}

func (a Array[T]) Length(ctx *core.Context) (int, error) {
	return len(a), nil
}

var _ RandomAccessList[scalar.Int] = Array[scalar.Int]{}

func (a Array[T]) At(ctx *core.Context, index int) (T, error) {
	return a[index], nil
}

var _ core.CanonicalList = Array[scalar.Int]{}

func (a Array[T]) IterCanonical(ctx *core.Context) iter.Seq2[core.Expression, error] {
	return func(yield func(core.Expression, error) bool) {
		for _, elem := range a {
			if !yield(elem, nil) {
				return
			}
		}
	}
}

var _ Dict[scalar.Int, scalar.Int] = Array[scalar.Int]{}

func (a Array[T]) Get(ctx *core.Context, key scalar.Int) (T, error) {
	i := int(key)
	if i < 0 || i >= len(a) {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds for array of length %d", i, len(a))
	}
	return a[i], nil
}

func (a Array[T]) Exists(ctx *core.Context, key scalar.Int) (bool, error) {
	i := int(key)
	return i >= 0 && i < len(a), nil
}

func (a Array[T]) Size(ctx *core.Context) (int, error) {
	return len(a), nil
}

func (a Array[T]) IterDict(ctx *core.Context) iter.Seq2[KV[scalar.Int, T], error] {
	return func(yield func(KV[scalar.Int, T], error) bool) {
		for i, v := range a {
			if !yield(KV[scalar.Int, T]{Key: scalar.Int(i), Value: v}, nil) {
				return
			}
		}
	}
}
