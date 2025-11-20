package scalar

import (
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

func (i Int) LogicalID() core.Identifier {
	return core.Identifier{
		Kind: "int",
		Key:  strconv.Itoa(int(i)),
	}
}

func (i Int) PhysicalID() core.Identifier {
	return i.LogicalID()
}

func (i Int) CanonicalID(*core.Context) (core.Identifier, error) {
	return i.LogicalID(), nil
}
