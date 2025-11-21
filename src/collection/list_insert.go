package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type ListInsert[T core.Expression] struct {
	List     RandomAccessList[T]
	Position int
	Element  T
}

var _ core.Expression = ListInsert[scalar.Int]{}

func (l ListInsert[T]) CanApply(transform any) bool {
	if _, ok := transform.(ListToArray[T]); ok {
		return true
	}
	return false
}

func (l ListInsert[T]) String() string {
	return fmt.Sprintf(`ListInsert(%v, %v, %v)`, l.List, l.Position, l.Element)
}

var _ RandomAccessList[scalar.Int] = ListInsert[scalar.Int]{}

func (l ListInsert[T]) Length(ctx *core.Context) (int, error) {
	length, err := l.List.Length(ctx)
	if err != nil {
		return 0, err
	}
	return length + 1, nil
}

func (l ListInsert[T]) IsEmpty(ctx *core.Context) (bool, error) {
	return false, nil
}

func (l ListInsert[T]) Iter(ctx *core.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		length, err := l.List.Length(ctx)
		if err != nil {
			var zero T
			yield(zero, err)
			return
		}

		i := 0
		for v, err := range l.List.Iter(ctx) {
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if i == l.Position {
				if !yield(l.Element, nil) {
					return
				}
			}
			if !yield(v, nil) {
				return
			}
			i++
		}

		if l.Position == length {
			yield(l.Element, nil)
		}
	}
}

func (l ListInsert[T]) At(ctx *core.Context, index int) (T, error) {
	if index < l.Position {
		return l.List.At(ctx, index)
	}
	if index == l.Position {
		return l.Element, nil
	}
	return l.List.At(ctx, index-1)
}
