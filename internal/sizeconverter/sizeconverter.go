package sizeconverter

import (
	"fmt"
	"strings"
	"strconv"
)

func ConvertSize(size string) (int64, error) {
	size = strings.TrimSpace(strings.ToLower(size))
	if size == "" {
		return 0, nil
	}

	multiplier := int64(1)
	
	switch {
	case strings.HasSuffix(size, "k"):
		multiplier = 1024
		size = strings.TrimSuffix(size, "k")
	case strings.HasSuffix(size, "m"):
		multiplier = 1024 * 1024
		size = strings.TrimSuffix(size, "m")
	case strings.HasSuffix(size, "g"):
		multiplier = 1024 * 1024 * 1024
		size = strings.TrimSuffix(size, "g")
	}

	sizeInBytes, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("convert to bytes: %w", err)
	}

	return sizeInBytes * multiplier, nil

}