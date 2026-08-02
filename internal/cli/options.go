package cli

import (
	"flag"
	"fmt"
	"os"
	"vessel/internal/cgroup/resources"
	"vessel/internal/sizeconverter"
	"strconv"
	"strings"
)

func ParseOptions() (resources.Limits, error) {
	memoryMaxParam := flag.String("memory", "", "maximum memory to be used by the unit")
	pidsMaxParam := flag.String("pids-max", "", "maximum tasks that can run in the unit")

	flag.CommandLine.Parse(os.Args[2:])

	memoryMax, err := sizeconverter.ConvertSize(*memoryMaxParam)
	if err != nil {
		return resources.Limits{}, fmt.Errorf("cannot convert memory max: %w", err)
	}

	var pidsMax uint64
	if strings.TrimSpace(*pidsMaxParam) != "" {
		pidsMax, err = strconv.ParseUint(*pidsMaxParam, 10, 64)
		if err != nil {
			return resources.Limits{}, fmt.Errorf("cannot convert pids max: %w", err)
		}
	}

	return resources.Limits{
		Memory: resources.Memory{
			Max: memoryMax,
		},
		Pids: resources.Pids{
			Max: pidsMax,
		},
	}, nil
}
