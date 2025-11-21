package collection

import (
	"fmt"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type ListRemovePosition[T core.Expression] struct {
	List     RandomAccessList[T]
	Position int
}

var _ core.Expression = ListRemovePosition[scalar.Int]{}

func (l ListRemovePosition[T]) CanApply(transform any) bool {
	return false
}

func (l ListRemovePosition[T]) String() string {
	return fmt.Sprintf(`ListRemovePosition(%v, %v)`, l.List, l.Position)
}

//TODO implement RandomAccessList for ListRemovePosition

type ListRemoveElement[T core.Ordered[T]] struct {
	List    SortedList[T]
	Element T
}

var _ core.Expression = ListRemoveElement[scalar.Int]{}

func (l ListRemoveElement[T]) CanApply(transform any) bool {
	return false
}

func (l ListRemoveElement[T]) String() string {
	return fmt.Sprintf(`ListRemoveElement(%v, %v)`, l.List, l.Element)
}

//TODO implement SortedList for ListRemoveElement
