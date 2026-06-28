package resources

import ()

type Memory struct {
	Min     int64
	Max     int64
	swapMax int64
}

func (m *Memory) setMax(Max int64) {
	m.Max = Max
}

func (m *Memory) setMin(Min int64) {
	m.Min = Min
}

func (m *Memory) setSwapMax(swapMax int64) {
	m.swapMax = swapMax
}