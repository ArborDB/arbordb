package scalar

import (
	"strconv"

	"github.com/ArborDB/arbordb/src/core"
)

type Int int

var _ core.Expression = Int(0)

func (i Int) CanApply(transform any) bool {
	return false
}

func (i Int) String() string {
	return strconv.Itoa(int(i))
}
