package scalar

import (
	"cmp"

	"github.com/ArborDB/arbordb/src/core"
)

type String string

var _ core.Expression = String("")

func (s String) String() string {
	return string(s)
}

var _ core.LogicalIdentifiable = String("")

func (s String) LogicalID() core.Identifier {
	return core.Identifier{
		Kind: "string",
		Key:  string(s),
	}
}

var _ core.PhysicalIdentifiable = String("")

func (s String) PhysicalID() core.Identifier {
	return s.LogicalID()
}

var _ core.CanonicalIdentifiable = String("")

func (s String) CanonicalID(*core.Context) (core.Identifier, error) {
	return s.LogicalID(), nil
}

var _ core.Ordered[String] = String("")

func (s String) Compare(to String) int {
	return cmp.Compare(s, to)
}
