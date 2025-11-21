package collection

import (
	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
	"iter"
)

type Array[T core.Expression] []T

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
