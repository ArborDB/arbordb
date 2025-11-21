package collection

import (
	"fmt"

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
	return false
}

func (l ListInsert[T]) String() string {
	return fmt.Sprintf(`ListInsert(%v, %v, %v)`, l.List, l.Position, l.Element)
}

//TODO implement RandomAccessList for ListInsert
