package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "config.yaml")
	configContent := `
image_url_file: image_urls.txt
download_directory: downloads
batch_size: 10
min_wait_time: 1.0
max_wait_time: 2.0
max_image_size_mb: 5
replace_downloaded_file_size: true
skip_if_file_exists: false
`
	err := ioutil.WriteFile(tempFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Read the config file
	config, err := ReadConfigFile(tempFile)
	assert.NoError(t, err)

	// Assert the values
	assert.Equal(t, "image_urls.txt", config.ImageURLFile)
	assert.Equal(t, "downloads", config.DownloadDirectory)
	assert.Equal(t, 10, config.BatchSize)
	assert.Equal(t, 1.0, config.MinWaitTime)
	assert.Equal(t, 2.0, config.MaxWaitTime)
	assert.Equal(t, "5", config.MaxImageSizeMB)
	assert.True(t, config.ReplaceDownloadedFileSize)
	assert.False(t, config.SkipIfFileExists)
}

func TestReadConfigFile_NotFound(t *testing.T) {
	// Read a non-existent config file
	configFile := "nonexistent.yaml"
	config, err := ReadConfigFile(configFile)

	// Assert the error and config
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestReadConfigFile_InvalidContent(t *testing.T) {
	// Create a temporary config file with invalid content
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "config.yaml")
	configContent := "invalid config content"
	err := ioutil.WriteFile(tempFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Read the config file
	config, err := ReadConfigFile(tempFile)

	// Assert the error and config
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestEnsureDownloadDirectory_DirectoryDoesNotExist(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Remove the download directory
	err := os.RemoveAll(tempDir)
	assert.NoError(t, err)

	// Ensure the download directory
	err = ensureDownloadDirectory(tempDir)
	assert.NoError(t, err)

	// Check if the directory exists
	_, err = os.Stat(tempDir)
	assert.NoError(t, err)
}
