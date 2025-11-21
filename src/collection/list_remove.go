package collection

import (
	"fmt"
	"iter"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type ListRemovePosition[T core.Expression] struct {
	List     RandomAccessList[T]
	Position int
}

var _ core.Expression = ListRemovePosition[scalar.Int]{}

func (l ListRemovePosition[T]) String() string {
	return fmt.Sprintf(`ListRemovePosition(%v, %v)`, l.List, l.Position)
}

var _ RandomAccessList[scalar.Int] = ListRemovePosition[scalar.Int]{}

func (l ListRemovePosition[T]) Length(ctx *core.Context) (int, error) {
	length, err := l.List.Length(ctx)
	if err != nil {
		return 0, err
	}
	return length - 1, nil
}

func (l ListRemovePosition[T]) IsEmpty(ctx *core.Context) (bool, error) {
	length, err := l.List.Length(ctx)
	if err != nil {
		return false, err
	}
	return length <= 1, nil
}

func (l ListRemovePosition[T]) Iter(ctx *core.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		i := 0
		for v, err := range l.List.Iter(ctx) {
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if i == l.Position {
				i++
				continue
			}
			if !yield(v, nil) {
				return
			}
			i++
		}
	}
}

func (l ListRemovePosition[T]) At(ctx *core.Context, index int) (T, error) {
	if index < l.Position {
		return l.List.At(ctx, index)
	}
	return l.List.At(ctx, index+1)
}

type ListRemoveElement[T core.Ordered[T]] struct {
	List    SortedList[T]
	Element T
}

var _ core.Expression = ListRemoveElement[scalar.Int]{}

func (l ListRemoveElement[T]) String() string {
	return fmt.Sprintf(`ListRemoveElement(%v, %v)`, l.List, l.Element)
}

var _ SortedList[scalar.Int] = ListRemoveElement[scalar.Int]{}

func (l ListRemoveElement[T]) Length(ctx *core.Context) (int, error) {
	length, err := l.List.Length(ctx)
	if err != nil {
		return 0, err
	}
	pos, err := l.List.BinarySearch(ctx, l.Element)
	if err != nil {
		return 0, err
	}
	if pos < 0 {
		return length, nil
	}
	return length - 1, nil
}

func (l ListRemoveElement[T]) IsEmpty(ctx *core.Context) (bool, error) {
	length, err := l.Length(ctx)
	if err != nil {
		return false, err
	}
	return length == 0, nil
}

func (l ListRemoveElement[T]) Iter(ctx *core.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		pos, err := l.List.BinarySearch(ctx, l.Element)
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
			if pos >= 0 && i == pos {
				i++
				continue
			}
			if !yield(v, nil) {
				return
			}
			i++
		}
	}
}

func (l ListRemoveElement[T]) At(ctx *core.Context, index int) (T, error) {
	pos, err := l.List.BinarySearch(ctx, l.Element)
	if err != nil {
		var zero T
		return zero, err
	}
	if pos < 0 { // not found, act as identity
		return l.List.At(ctx, index)
	}
	if index < pos {
		return l.List.At(ctx, index)
	}
	return l.List.At(ctx, index+1)
}

func (l ListRemoveElement[T]) BinarySearch(ctx *core.Context, target T) (int, error) {
	removedPos, err := l.List.BinarySearch(ctx, l.Element)
	if err != nil {
		return 0, err
	}

	if removedPos < 0 {
		// not in list, so this is a no-op
		return l.List.BinarySearch(ctx, target)
	}

	foundPos, err := l.List.BinarySearch(ctx, target)
	if err != nil {
		return 0, err
	}

	if foundPos >= 0 { // found in original list
		if foundPos == removedPos {
			return -removedPos - 1, nil
		}
		if foundPos > removedPos {
			return foundPos - 1, nil
		}
		return foundPos, nil
	}

	// not found in original list
	insertionPoint := -foundPos - 1
	if insertionPoint > removedPos {
		return -(insertionPoint - 1) - 1, nil
	}
	return foundPos, nil
}
