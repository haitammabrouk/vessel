package cli

import (
	"flag"
	"fmt"
	"os"
	"vessel/internal/cgroup/resources"
	"vessel/internal/sizeconverter"
)

func ParseOptions() (resources.Limits, error) {
	memoryMax := flag.String("memory", "", "memory limit")

	flag.CommandLine.Parse(os.Args[2:])

	memoryMaxInBytes, err := sizeconverter.ConvertSize(*memoryMax)
	if err != nil {
		return resources.Limits{}, fmt.Errorf("convert memory max to bytes: %w", err)
	}

	return resources.Limits{
		Memory: resources.Memory{
			Max: memoryMaxInBytes,
		}}, nil
}