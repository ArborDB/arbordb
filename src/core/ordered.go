package core

type Ordered[T Expression] interface {
	Expression
	Compare(to T) int
}
