package dshash

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func BenchmarkMap(b *testing.B) {
	// Small map
	smallMap := make(map[string]int)
	for i := range 100 {
		smallMap[string('a'+byte(i))] = i
	}

	// Large map with complex keys
	type ComplexKey struct {
		ID   int
		Name string
	}
	largeMap := make(map[ComplexKey]string)
	for i := range 1000 {
		largeMap[ComplexKey{
			ID:   i,
			Name: fmt.Sprintf("item-%d", i),
		}] = fmt.Sprintf("value-%v", i)
	}

	b.Run("SmallMap", func(b *testing.B) {
		for b.Loop() {
			state := sha256.New()
			_ = Hash(state, smallMap)
			_ = state.Sum(nil)
		}
	})

	b.Run("LargeMap_ComplexKeys", func(b *testing.B) {
		for b.Loop() {
			state := sha256.New()
			_ = Hash(state, largeMap)
			_ = state.Sum(nil)
		}
	})
}
