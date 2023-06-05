package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ImageDownloader struct {
	HTTPClient  HTTPClient
	FileChecker FileChecker
}

func (d *ImageDownloader) DownloadImage(url, downloadDir string) error {
	fileName := filepath.Base(url)
	filePath := filepath.Join(downloadDir, fileName)

	// Check if the file already exists
	if d.FileChecker.IsFileExists(filePath) {
		// File already exists, skip downloading
		return nil
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Download the image
	resp, err := d.HTTPClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image, status: %s", resp.Status)
	}

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}

func batchImageURLs(imageURLs []string, batchSize int) [][]string {
	var batches [][]string
	length := len(imageURLs)

	for i := 0; i < length; i += batchSize {
		end := i + batchSize
		if end > length {
			end = length
		}
		batch := imageURLs[i:end]
		batches = append(batches, batch)
	}

	return batches
}
