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

//TODO implement Dict for Array
