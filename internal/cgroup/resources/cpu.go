package resources

import ()

type CpuMax struct {
	Quota  int64
	Period int64
}

type Cpu struct {
	Max CpuMax
}