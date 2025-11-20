package scalar

import "github.com/ArborDB/arbordb/src/core"

type String string

var _ core.Expression = String("")

func (s String) CanApply(transform any) bool {
	switch transform.(type) {
	case core.ToLogicalID, core.ToPhysicalID, core.ToCanonicalID:
		return true
	}
	return false
}

func (s String) String() string {
	return string(s)
}

func (s String) LogicalID() core.Identifier {
	return core.Identifier{
		Kind: "string",
		Key:  string(s),
	}
}

func (s String) PhysicalID() core.Identifier {
	return s.LogicalID()
}

func (s String) CanonicalID(*core.Context) (core.Identifier, error) {
	return s.LogicalID(), nil
}
