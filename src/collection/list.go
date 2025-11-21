package collection

import (
	"iter"

	"github.com/ArborDB/arbordb/src/core"
)

type List[T core.Expression] interface {
	Length(ctx *core.Context) (int, error)
	IsEmpty(ctx *core.Context) (bool, error)
	Iter(ctx *core.Context) iter.Seq2[T, error]
}

type RandomAccessList[T core.Expression] interface {
	List[T]
	At(ctx *core.Context, index int) (T, error)
}

type SortedList[T core.Ordered[T]] interface {
	RandomAccessList[T]
	BinarySearch(ctx *core.Context, target T) (int, error)
}
