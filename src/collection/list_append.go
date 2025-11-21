package collection

import (
	"fmt"

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

//TODO implement List for ListAppend
