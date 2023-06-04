package main

import (
	"os"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	// Create a temporary config file
	configFile, err := os.CreateTemp("/tmp", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(configFile.Name())

	// Test case: Valid config file
	configContent := `
image_url_file: image_urls.txt
download_directory: downloads
batch_size: 10
min_wait_time: 0.8
max_wait_time: 3.0
max_image_size_mb: 10
replace_downloaded_file_size: true
skip_if_file_exists: true
`
	_, err = configFile.WriteString(configContent)
	if err != nil {
		t.Fatalf("Failed to write to the config file: %v", err)
	}

	// Read the config file
	config, err := ReadConfigFile(configFile.Name())
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Verify the values of the config
	if config.ImageURLFile != "image_urls.txt" {
		t.Errorf("Expected ImageURLFile to be 'image_urls.txt', but got '%s'", config.ImageURLFile)
	}

	if config.DownloadDirectory != "downloads" {
		t.Errorf("Expected DownloadDirectory to be 'downloads', but got '%s'", config.DownloadDirectory)
	}

	if config.BatchSize != 10 {
		t.Errorf("Expected BatchSize to be 10, but got %d", config.BatchSize)
	}

	if config.MinWaitTime != 0.8 {
		t.Errorf("Expected MinWaitTime to be 0.8, but got %.2f", config.MinWaitTime)
	}

	if config.MaxWaitTime != 3.0 {
		t.Errorf("Expected MaxWaitTime to be 3.0, but got %.2f", config.MaxWaitTime)
	}

	if config.MaxImageSizeMB != "10" {
		t.Errorf("Expected MaxImageSizeMB to be '10', but got '%s'", config.MaxImageSizeMB)
	}

	if !config.ReplaceDownloadedFileSize {
		t.Errorf("Expected ReplaceDownloadedFileSize to be true, but got false")
	}

	if !config.SkipIfFileExists {
		t.Errorf("Expected SkipIfFileExists to be true, but got false")
	}
}
