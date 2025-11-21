package scalar

import (
	"cmp"
	"strconv"

	"github.com/ArborDB/arbordb/src/core"
)

type Int int

var _ core.Expression = Int(0)

func (i Int) CanApply(transform any) bool {
	switch transform.(type) {
	case core.ToLogicalID, core.ToPhysicalID, core.ToCanonicalID:
		return true
	}
	return false
}

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

var _ core.LogicalIdentifiable = Int(0)

func (i Int) LogicalID() core.Identifier {
	return core.Identifier{
		Kind: "int",
		Key:  strconv.Itoa(int(i)),
	}
}

var _ core.PhysicalIdentifiable = Int(0)

func (i Int) PhysicalID() core.Identifier {
	return i.LogicalID()
}

var _ core.CanonicalIdentifiable = Int(0)

func (i Int) CanonicalID(*core.Context) (core.Identifier, error) {
	return i.LogicalID(), nil
}

var _ core.Ordered[Int] = Int(0)

func (i Int) Compare(to Int) int {
	return cmp.Compare(i, to)
}
