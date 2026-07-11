package cli

import (
	"flag"
	"fmt"
	"os"
	"vessel/internal/cgroup/resources"
	"vessel/internal/sizeconverter"
)

func ParseOptions() (resources.ResouceLimits, error) {
	memoryMax := flag.String("memory", "", "memory limit")
	memorySwapMax := flag.String("memory-swap-max", "", "memory swap max")

	flag.CommandLine.Parse(os.Args[2:])

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
		}}, nil
}