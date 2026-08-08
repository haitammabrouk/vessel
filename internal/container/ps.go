package container

import (
	"os"
	"fmt"
	"vessel/internal/id"
	"vessel/internal/metadata"
)

func ListContainers() error {
	metadataDir, err := os.ReadDir(metadata.MetadataPath)
	if err != nil {
		return err
	}
	for _, entry := range metadataDir {
		fmt.Println(id.GetShortHandId(entry.Name()))
	}
	return nil
}