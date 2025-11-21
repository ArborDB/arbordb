package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type ListAppend[T core.Expression] struct {
	List    List[T]
	Element T
}

var _ core.Expression = ListAppend[scalar.Int]{}

func (l ListAppend[T]) CanApply(transform any) bool {
	return false
}

func (l ListAppend[T]) String() string {
	return fmt.Sprintf(`ListAppend(%v, %v)`, l.List, l.Element)
}

var _ List[scalar.Int] = ListAppend[scalar.Int]{}

func (l ListAppend[T]) Length(ctx *core.Context) (int, error) {
	length, err := l.List.Length(ctx)
	if err != nil {
		return 0, err
	}
	return length + 1, nil
}

func (l ListAppend[T]) IsEmpty(ctx *core.Context) (bool, error) {
	return false, nil
}

func (l ListAppend[T]) Iter(ctx *core.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for v, err := range l.List.Iter(ctx) {
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if !yield(v, nil) {
				return
			}
		}
		yield(l.Element, nil)
	}
}
