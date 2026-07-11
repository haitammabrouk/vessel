package cli

import (
	"fmt"
	"flag"
	"vessel/internal/cgroup/resources"
	"vessel/internal/sizeconverter"
)

func ParseOptions() (resources.ResouceLimits, error) {
	memoryMax := flag.String("memory", "", "memory limit")
	memorySwapMax := flag.String("memory-max", "", "memory swap max")
	cpuMax := flag.Int64("cpu", 0, "cpu max")

	flag.Parse()

	memoryMaxInBytes, err := sizeconverter.ConvertSize(*memoryMax)
	if err != nil {
		return resources.ResouceLimits{}, fmt.Errorf("convert memory max to bytes: %w", err)
	}
	memorySwapMaxInBytes, err := sizeconverter.ConvertSize(*memorySwapMax)
	if err != nil {
		return resources.ResouceLimits{}, fmt.Errorf("convert memory swap max to bytes: %w", err)
	}

	return resources.ResouceLimits{
		Memory: resources.Memory{
			Max: memoryMaxInBytes,
			SwapMax: memorySwapMaxInBytes,
		},
		Cpu: resources.Cpu{Max: *cpuMax}}, nil
}