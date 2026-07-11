package resources

type Memory struct {
	Max     int64
	SwapMax int64
}

type Cpu struct {
	Max  int64
}

type ResouceLimits struct {
	Memory Memory
	Cpu Cpu
}