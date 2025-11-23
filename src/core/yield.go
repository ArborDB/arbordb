package core

import "time"

func (c *Context) Yield() error {
	// empty context
	if c == nil {
		return nil
	}

	// no yield func
	if c.YieldFunc == nil {
		return nil
	}

	switch {

	case c.YieldQuota > 0:
		// quota based yield
		c.YieldQuota--
		return nil

	case c.YieldInterval > 0:
		// epoch based yield
		currentEpoch := globalEpochProvider.GetEpoch()
		elapsed := time.Duration((currentEpoch - c.lastYieldEpoch)) * epochDuration
		if elapsed >= c.YieldInterval {
			// yield
			c.lastYieldEpoch = currentEpoch
		} else {
			// not yield
			return nil
		}

	}

	// yield
	if !c.YieldFunc() {
		return Err[ErrCanceled]()
	}

	return nil
}
