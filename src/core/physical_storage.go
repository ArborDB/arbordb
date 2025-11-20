package core

type PhysicalStorage interface {
	Set(ctx *Context, expr Expression) (id Identifier, err error)
	Get(ctx *Context, id Identifier, target any) (err error)
}
