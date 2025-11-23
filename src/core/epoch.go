package core

import (
	"runtime"
	"sync/atomic"
	"time"

	"golang.org/x/sys/cpu"
)

const (
	epochDuration = time.Millisecond * 10
)

var globalEpochProvider = NewEpochProvider()

type epochShard struct {
	_ cpu.CacheLinePad
	N atomic.Uint64
}

type EpochProvider struct {
	shards    []epochShard
	maxShards int
}

func NewEpochProvider() *EpochProvider {
	maxShards := runtime.GOMAXPROCS(0)
	shards := make([]epochShard, maxShards)
	ret := &EpochProvider{
		shards:    shards,
		maxShards: maxShards,
	}

	exit := make(chan bool)
	go func() {
		ticker := time.NewTicker(epochDuration)
		for {
			select {
			case <-ticker.C:
				for i := range maxShards {
					shards[i].N.Add(1)
				}
			case <-exit:
				return
			}
		}
	}()
	runtime.AddCleanup(ret, func(exit chan bool) {
		close(exit)
	}, exit)

	return ret
}

func (e *EpochProvider) GetEpoch() uint64 {
	pid := runtime_procPin() % e.maxShards
	defer runtime_procUnpin()
	return e.shards[pid%e.maxShards].N.Load()
}
