package core

import (
	"context"
	"time"
)

type Context struct {
	context.Context
	PhysicalStorage PhysicalStorage
	Cost            Cost

	YieldFunc      YieldFunc
	YieldQuota     int
	YieldInterval  time.Duration
	lastYieldEpoch uint64
}

type YieldFunc = func() bool
