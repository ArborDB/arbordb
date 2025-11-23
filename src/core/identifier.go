package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"reflect"
	"unsafe"

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

type ToLogicalID struct{}

type LogicalIdentifiable interface {
	LogicalID(*Context) (Identifier, error)
}

var _ Transform[Expression, Identifier] = ToLogicalID{}

func (g ToLogicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToLogicalID) Apply(ctx *Context, from Expression, to *Identifier) error {
	switch from := from.(type) {

	case LogicalIdentifiable:
		id, err := from.LogicalID(ctx)
		if err != nil {
			return err
		}
		*to = id
		return nil

	}

	id, err := structuralHash(from)
	if err != nil {
		return err
	}
	*to = id

	return nil
}

type ToPhysicalID struct{}

type PhysicalIdentifiable interface {
	PhysicalID(*Context) (Identifier, error)
}

var _ Transform[Expression, Identifier] = ToPhysicalID{}

func (g ToPhysicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToPhysicalID) Apply(ctx *Context, from Expression, to *Identifier) error {
	switch from := from.(type) {

	case PhysicalIdentifiable:
		id, err := from.PhysicalID(ctx)
		if err != nil {
			return err
		}
		*to = id
		return nil

	}

	id, err := structuralHash(from)
	if err != nil {
		return err
	}
	*to = id

	return nil
}

type ToCanonicalID struct{}

type CanonicalIdentifiable interface {
	CanonicalID(*Context) (Identifier, error)
}

var _ Transform[Expression, Identifier] = ToCanonicalID{}

func (g ToCanonicalID) EstimateCost(ctx *Context, from Expression) (Cost, error) {
	return Cost{}, nil
}

func (g ToCanonicalID) Apply(ctx *Context, from Expression, to *Identifier) error {
	state := sha256.New()
	visited := make(map[unsafe.Pointer]struct{})

	var hashRecursive func(h hash.Hash, expr Expression) error
	hashRecursive = func(h hash.Hash, expr Expression) error {
		if expr == nil {
			h.Write([]byte{CanonicalTagNil})
			return nil
		}

		val := reflect.ValueOf(expr)
		if val.Kind() == reflect.Pointer && !val.IsNil() {
			ptr := val.UnsafePointer()
			if _, ok := visited[ptr]; ok {
				return nil // cycle detected
			}
			visited[ptr] = struct{}{}
		}

		switch e := expr.(type) {
		case CanonicalIdentifiable:
			id, err := e.CanonicalID(ctx)
			if err != nil {
				return err
			}
			// Hash the identifier string representation, e.g., "int:42"
			h.Write([]byte{byte(CanonicalTagString)})
			h.Write([]byte(id.String()))
			return nil

		case CanonicalList:
			h.Write([]byte{byte(CanonicalTagList)})
			for item, err := range e.IterCanonical(ctx) {
				if err != nil {
					return err
				}
				if err := hashRecursive(h, item); err != nil {
					return err
				}
			}
			h.Write([]byte{byte(CanonicalTagListEnd)})
			return nil

		default:
			return fmt.Errorf("type %T does not support canonical identification", expr)
		}
	}

	if err := hashRecursive(state, from); err != nil {
		return err
	}

	*to = Identifier{
		Kind: "canonical-sha256",
		Key:  hex.EncodeToString(state.Sum(nil)),
	}

	return nil
}

func structuralHash(expr Expression) (id Identifier, err error) {
	state := sha256.New()
	if err = dshash.Hash(state, expr); err != nil {
		return
	}
	return Identifier{
		Kind: "dshash-sha256",
		Key:  hex.EncodeToString(state.Sum(nil)),
	}, nil
}
