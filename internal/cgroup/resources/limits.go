package resources

type Memory struct {
	Max     uint64
}

type Pids struct {
	Max uint64
}

type Limits struct {
	Memory Memory
	Pids Pids
}
