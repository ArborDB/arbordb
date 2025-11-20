package scalar

import "github.com/ArborDB/arbordb/src/core"

type String string

var _ core.Expression = String("")

func (s String) CanApply(transform any) bool {
	return false
}

func (s String) String() string {
	return string(s)
}
