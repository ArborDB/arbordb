package core

func Apply[To Expression](
	ctx *Context,
	yield func(*TransformStep, error) bool,
	transform Transform[To],
	from Expression,
	to *To,
) {
	for step, err := range transform.Apply(ctx, from, to) {
		ctx.Cost.Merge(step.Cost)
		if !yield(step, err) {
			break
		}
	}
}
