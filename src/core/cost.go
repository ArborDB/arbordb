package core

type Cost struct {
	CPU        int
	IO         int
	PeakMemory int
}

func (c *Cost) Merge(c2 Cost) {
	c.CPU += c2.CPU
	c.IO += c2.IO
	c.PeakMemory = max(c.PeakMemory, c2.PeakMemory)
}
