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

type dagNode struct {
	Name     string
	Children []Expression
}

var _ CanonicalList = &dagNode{}

func (d *dagNode) String() string { return d.Name }

func (d *dagNode) IterCanonical(ctx *Context) iter.Seq2[Expression, error] {
	return func(yield func(Expression, error) bool) {
		for _, child := range d.Children {
			if !yield(child, nil) {
				return
			}
		}
	}
}

func TestCanonicalDAG(t *testing.T) {
	leaf := &dagNode{Name: "leaf"}
	// DAG: root -> [leaf, leaf] (shared pointer)
	dagRoot := &dagNode{
		Name:     "root",
		Children: []Expression{leaf, leaf},
	}

	// Tree: root -> [leaf1, leaf2] (distinct pointers, same value)
	treeRoot := &dagNode{
		Name: "root",
		Children: []Expression{
			&dagNode{Name: "leaf"},
			&dagNode{Name: "leaf"},
		},
	}

	var transformer ToCanonicalID
	var idDAG Identifier
	if err := transformer.Apply(nil, dagRoot, &idDAG); err != nil {
		t.Fatal(err)
	}

	var idTree Identifier
	if err := transformer.Apply(nil, treeRoot, &idTree); err != nil {
		t.Fatal(err)
	}

	if idDAG != idTree {
		t.Fatalf("DAG and Tree should have same Canonical ID: %v vs %v", idDAG, idTree)
	}
}
