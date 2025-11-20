package core

import "fmt"

type Expression interface {
	fmt.Stringer
	CanApply(transform any) bool
}
