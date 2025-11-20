package core

import "iter"

type Transform[To Expression] interface {
	EstimateCost(ctx *Context, from Expression) (Cost, error)
	Apply(ctx *Context, from Expression, to *To) (steps TransformSteps)
}

type TransformSteps = iter.Seq2[*TransformStep, error]

// TransformStep is for information exchanging between transform process and scheduler
type TransformStep struct {
	// resource cost during the last step
	Cost Cost
}
