package core

import "context"

type Context struct {
	context.Context
	PhysicalStorage PhysicalStorage
	Cost            Cost
}
