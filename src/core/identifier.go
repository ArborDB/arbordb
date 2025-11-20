package core

import "fmt"

type Identifier struct {
	Kind string
	Key  string
}

var _ Expression = Identifier{}

func (i Identifier) String() string {
	return i.Kind + ":" + i.Key
}

func (i Identifier) CanApply(transform any) bool {
	return false
}

func (i Identifier) LogicalID() Identifier {
	return i
}

func (i Identifier) PhysicalID() Identifier {
	return i
}

func (i Identifier) CanonicalID(*Context) (Identifier, error) {
	return i, nil
}

type ToLogicalID struct{}

type LogicalIdentifiable interface {
	LogicalID() Identifier
}

var _ Transform[Identifier] = ToLogicalID{}

func (g ToLogicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToLogicalID) Apply(ctx *Context, from Expression, to *Identifier) TransformSteps {
	return func(yield func(*TransformStep, error) bool) {
		switch from := from.(type) {
		case LogicalIdentifiable:
			*to = from.LogicalID()
			return
		}
		yield(nil, fmt.Errorf("logical id not supported: %T", from))
	}
}

type ToPhysicalID struct{}

type PhysicalIdentifiable interface {
	PhysicalID() Identifier
}

var _ Transform[Identifier] = ToPhysicalID{}

func (g ToPhysicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToPhysicalID) Apply(ctx *Context, from Expression, to *Identifier) TransformSteps {
	return func(yield func(*TransformStep, error) bool) {
		switch from := from.(type) {
		case PhysicalIdentifiable:
			*to = from.PhysicalID()
			return
		}
		yield(nil, fmt.Errorf("physical id not supported: %T", from))
	}
}

type ToCanonicalID struct{}

type CanonicalIdentifiable interface {
	CanonicalID(*Context) (Identifier, error)
}

var _ Transform[Identifier] = ToCanonicalID{}

func (g ToCanonicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToCanonicalID) Apply(ctx *Context, from Expression, to *Identifier) TransformSteps {
	return func(yield func(*TransformStep, error) bool) {
		switch from := from.(type) {
		case CanonicalIdentifiable:
			id, err := from.CanonicalID(ctx)
			if err != nil {
				yield(nil, err)
				return
			}
			*to = id
			return
		}
		yield(nil, fmt.Errorf("canonical id not supported: %T", from))
	}
}
