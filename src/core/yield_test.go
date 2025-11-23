package core

import (
	"testing"
)

func BenchmarkYield(b *testing.B) {
	ctx := &Context{
		YieldFunc: func() bool {
			return true
		},
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ctx.Yield(); err != nil {
			b.Fatal(err)
		}
	}
}
