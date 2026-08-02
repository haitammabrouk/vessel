package container

import (
	"os"
	"fmt"
	"vessel/internal/id"
)

const metadataPath = "/var/lib/vessel/containers"

func ListContainers() error {
	metadataDir, err := os.ReadDir(metadataPath)
	if err != nil {
		return err
	}
	for _, entry := range metadataDir {
		fmt.Println(id.GetShortHandId(entry.Name()))
	}
	return nil
}