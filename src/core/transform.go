package core

type Transform[From Expression, To Expression] interface {
	EstimateCost(ctx *Context, from From) (Cost, error)
	Apply(ctx *Context, from From, to *To) error
}
