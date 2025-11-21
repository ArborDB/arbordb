package core

import "iter"

// CanonicalList represents an expression that has a list-like canonical form.
type CanonicalList interface {
	Expression
	// IterCanonical iterates through the logical elements of the expression
	// in a deterministic order.
	IterCanonical(ctx *Context) iter.Seq2[Expression, error]
}

const (
	CanonicalTagInvalid = 0
	CanonicalTagNil     = 10
	CanonicalTagString  = 20
	CanonicalTagList    = 30
	CanonicalTagListEnd = 40
)
