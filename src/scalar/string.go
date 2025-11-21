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

func (s String) LogicalID(ctx *core.Context) (core.Identifier, error) {
	return core.Identifier{
		Kind: "string",
		Key:  string(s),
	}, nil
}

var _ core.PhysicalIdentifiable = String("")

func (s String) PhysicalID(ctx *core.Context) (core.Identifier, error) {
	return s.LogicalID(ctx)
}

var _ core.CanonicalIdentifiable = String("")

func (s String) CanonicalID(ctx *core.Context) (core.Identifier, error) {
	return s.LogicalID(ctx)
}

var _ core.Ordered[String] = String("")

func (s String) Compare(to String) int {
	return cmp.Compare(s, to)
}
