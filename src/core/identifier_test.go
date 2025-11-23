package core

import (
	"iter"
	"testing"
)

type recursiveList struct {
	Self *recursiveList
}

var _ CanonicalList = &recursiveList{}

func (r *recursiveList) String() string {
	return "RecursiveList"
}

func (r *recursiveList) IterCanonical(ctx *Context) iter.Seq2[Expression, error] {
	return func(yield func(Expression, error) bool) {
		if r.Self != nil {
			if !yield(r.Self, nil) {
				return
			}
		}
	}
}

func TestCanonicalIDCollision(t *testing.T) {
	// Case 1: Recursive list (points to self)
	// Structure: List -> [Self]
	recursive := &recursiveList{}
	recursive.Self = recursive

	// Case 2: Empty list
	// Structure: List -> []
	empty := &recursiveList{}

	var transformer ToCanonicalID
	var idRecursive Identifier
	if err := transformer.Apply(nil, recursive, &idRecursive); err != nil {
		t.Fatal(err)
	}

	var idEmpty Identifier
	if err := transformer.Apply(nil, empty, &idEmpty); err != nil {
		t.Fatal(err)
	}

	if idRecursive == idEmpty {
		t.Fatalf("Collision detected: Recursive list and empty list have the same Canonical ID: %v", idRecursive)
	}
}
