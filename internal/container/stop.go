package container

import (
	"fmt"
	"os"
	"strings"
	"vessel/internal/cgroup"
)

func StopContainer(containerId string) error {
	metadataDir, err := os.ReadDir(metadataPath)
	isFound := false
	if err != nil {
		return err
	}

	for _, entry := range metadataDir {
		if strings.HasPrefix(entry.Name(), containerId) {
			isFound = true
			if err := cgroup.StopUnitScope(entry.Name()); err != nil {
				return err
			}
		}
	}

	if !isFound {
		fmt.Printf("cannot find a container with id: %s", containerId)
	}
	return nil
}