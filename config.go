package main

import (
	"fmt"
)

type Config struct {
	ImageURLFile              string
	DownloadDirectory         string
	BatchSize                 int
	MinWaitTime               float64
	MaxWaitTime               float64
	MaxImageSizeMB            string
	ReplaceDownloadedFileSize bool
	SkipIfFileExists          bool
}

func parseMaxImageSize(size string) (int64, error) {
	if size == "MAX" {
		return -1, nil
	}

	return parseSize(size)
}

func parseSize(size string) (int64, error) {
	// Parse size with unit suffix (e.g., 10KB, 1MB, 1GB)
	var value int64
	var unit string
	_, err := fmt.Sscanf(size, "%d%s", &value, &unit)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size: %w", err)
	}

	// Convert size to bytes
	switch unit {
	case "B":
		return value, nil
	case "KB":
		return value * 1024, nil
	case "MB":
		return value * 1024 * 1024, nil
	case "GB":
		return value * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("invalid size unit: %s", unit)
	}
}
