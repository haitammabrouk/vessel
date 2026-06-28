package resources

import ()

type Memory struct {
	min     int64
	max     int64
	swapMax int64
}

func (m *Memory) setMax(max int64) {
	m.max = max
}

func (m *Memory) setMin(min int64) {
	m.min = min
}

func (m *Memory) setSwapMax(swapMax int64) {
	m.swapMax = swapMax
}