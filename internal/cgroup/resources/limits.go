package resources

type Memory struct {
	Max     uint64
	SwapMax uint64
}

type Limits struct {
	Memory Memory
}