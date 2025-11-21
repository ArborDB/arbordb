package core

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/ArborDB/arbordb/src/dshash"
)

type Identifier struct {
	Kind string
	Key  string
}

var _ Expression = Identifier{}

func (i Identifier) String() string {
	return i.Kind + ":" + i.Key
}

func (i Identifier) LogicalID(*Context) (Identifier, error) {
	return i, nil
}

func (i Identifier) PhysicalID(*Context) (Identifier, error) {
	return i, nil
}

func (i Identifier) CanonicalID(*Context) (Identifier, error) {
	return i, nil
}

type ToLogicalID struct{}

type LogicalIdentifiable interface {
	LogicalID(*Context) (Identifier, error)
}

var _ Transform[Expression, Identifier] = ToLogicalID{}

func (g ToLogicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToLogicalID) Apply(ctx *Context, from Expression, to *Identifier) TransformSteps {
	return func(yield func(*TransformStep, error) bool) {
		switch from := from.(type) {

		case LogicalIdentifiable:
			id, err := from.LogicalID(ctx)
			if err != nil {
				yield(nil, err)
				return
			}
			*to = id
			return

		}

		id, err := hash(from)
		if err != nil {
			yield(nil, err)
			return
		}
		*to = id

	}
}

type ToPhysicalID struct{}

type PhysicalIdentifiable interface {
	PhysicalID(*Context) (Identifier, error)
}

var _ Transform[Expression, Identifier] = ToPhysicalID{}

func (g ToPhysicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToPhysicalID) Apply(ctx *Context, from Expression, to *Identifier) TransformSteps {
	return func(yield func(*TransformStep, error) bool) {
		switch from := from.(type) {

		case PhysicalIdentifiable:
			id, err := from.PhysicalID(ctx)
			if err != nil {
				yield(nil, err)
				return
			}
			*to = id
			return

		}

		id, err := hash(from)
		if err != nil {
			yield(nil, err)
			return
		}
		*to = id

	}
}

type ToCanonicalID struct{}

type CanonicalIdentifiable interface {
	CanonicalID(*Context) (Identifier, error)
}

var _ Transform[Expression, Identifier] = ToCanonicalID{}

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

		//TODO resolve all identifiers recursively in `from`

		id, err := hash(from)
		if err != nil {
			yield(nil, err)
			return
		}
		*to = id

	}
}

func hash(expr Expression) (id Identifier, err error) {
	state := sha256.New()
	if err = dshash.Hash(state, expr); err != nil {
		return
	}
	return Identifier{
		Kind: "dshash-sha256",
		Key:  hex.EncodeToString(state.Sum(nil)),
	}, nil
}
