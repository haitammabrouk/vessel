package container

import (
	"fmt"
	"os"
	"strings"
	"vessel/internal/cgroup"
	"vessel/internal/id"
	"vessel/internal/metadata"
)

func StopContainer(containerId string) error {
	if strings.TrimSpace(containerId) == "" {
		return fmt.Errorf("no container id is provided")
	}
	metadataDir, err := os.ReadDir(metadata.MetadataPath)
	if err != nil {
		return err
	}
	
	matchedIds := make([]string, 0)
	for _, entry := range metadataDir {
		if strings.HasPrefix(entry.Name(), containerId) {
			matchedIds = append(matchedIds, entry.Name())
		}
	}

	shortedIds := make([]string, 0)
	for _, entry := range matchedIds {
		shortedIds = append(shortedIds, id.GetShortHandId(entry))
	}

	if len(matchedIds) == 1 {
		if err := cgroup.StopUnitScope(matchedIds[0]); err != nil {
			return err
		}
	} else if len(matchedIds) > 1 {
		return fmt.Errorf("container id %s is ambiguous, matched multiple containers:\n%s\nprovide more characters to disambiguate",
    			containerId, strings.Join(shortedIds, "\n"))
	} else {
		return fmt.Errorf("cannot find a container with id: %s", containerId)
	}

	return nil
}