package scalar

import (
	"cmp"
	"strconv"

	"github.com/ArborDB/arbordb/src/core"
)

type Int int

var _ core.Expression = Int(0)

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

var _ core.LogicalIdentifiable = Int(0)

func (i Int) LogicalID(ctx *core.Context) (core.Identifier, error) {
	return core.Identifier{
		Kind: "int",
		Key:  strconv.Itoa(int(i)),
	}, nil
}

var _ core.PhysicalIdentifiable = Int(0)

func (i Int) PhysicalID(ctx *core.Context) (core.Identifier, error) {
	return i.LogicalID(ctx)
}

var _ core.CanonicalIdentifiable = Int(0)

func (i Int) CanonicalID(ctx *core.Context) (core.Identifier, error) {
	return i.LogicalID(ctx)
}

var _ core.Ordered[Int] = Int(0)

func (i Int) Compare(to Int) int {
	return cmp.Compare(i, to)
}
