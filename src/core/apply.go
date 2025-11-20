package core

func Apply[From Expression, To Expression](
	ctx *Context,
	yield func(*TransformStep, error) bool,
	transform Transform[From, To],
	from From,
	to *To,
) {
	for step, err := range transform.Apply(ctx, from, to) {
		ctx.Cost.Merge(step.Cost)
		if !yield(step, err) {
			break
		}
	}
}
