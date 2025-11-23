package core

import (
	"testing"
	"time"
)

func TestEpoch(t *testing.T) {
	provider := NewEpochProvider()
	n := provider.GetEpoch()
	time.Sleep(epochDuration * 2)
	n2 := provider.GetEpoch()
	if n2 <= n {
		t.Fatal()
	}
}

func BenchmarkEpoch(b *testing.B) {
	provider := NewEpochProvider()
	b.ResetTimer()
	for b.Loop() {
		provider.GetEpoch()
	}
}

func BenchmarkEpochParallel(b *testing.B) {
	provider := NewEpochProvider()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			provider.GetEpoch()
		}
	})
}
