package resources

import ()

type CpuMax struct {
	quota int64
	period int64
}

type Cpu struct {
	max CpuMax
}