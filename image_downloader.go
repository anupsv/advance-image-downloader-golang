package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func downloadImages(config *Config) error {
	// Check if the batch size is less than 10
	if config.BatchSize < 10 {
		return fmt.Errorf("batch size should be at least 10")
	}

	// Ensure the download directory exists
	err := ensureDownloadDirectory(config.DownloadDirectory)
	if err != nil {
		return fmt.Errorf("failed to ensure download directory: %v", err)
	}

	imageURLs, err := readImageURLsFromFile(config.ImageURLFile)
	if err != nil {
		return fmt.Errorf("failed to read image URLs from file: %v", err)
	}

	batches := batchImageURLs(imageURLs, config.BatchSize)
	for _, batch := range batches {
		for _, url := range batch {
			// Check if the file already exists
			filePath := filepath.Join(config.DownloadDirectory, filepath.Base(url))
			if isFileExists(filePath) {
				if config.SkipIfFileExists {
					// Skip downloading if the file already exists
					continue
				} else if config.ReplaceDownloadedFileSize {
					// Check if the file size has changed, replace if it has
					newSize, err := getImageFileSize(url)
					if err != nil {
						return fmt.Errorf("failed to get image file size: %v", err)
					}
					currentSize, err := getFileSize(filePath)
					if err != nil {
						return fmt.Errorf("failed to get current file size: %v", err)
					}
					if newSize == currentSize {
						// Skip downloading if the file size hasn't changed
						continue
					}
				}
			}

			// Check if the image size is exceeded
			if config.MaxImageSizeMB != "MAX" {
				maxSize, err := parseMaxImageSize(config.MaxImageSizeMB)
				if err != nil {
					return fmt.Errorf("failed to parse max image size: %v", err)
				}
				if isImageSizeExceeded(url, maxSize) {
					continue
				}
			}

			err := downloadImage(url, config.DownloadDirectory)
			if err != nil {
				return fmt.Errorf("failed to download image: %v", err)
			}
		}

		waitTime := generateRandomWaitTime(config.MinWaitTime, config.MaxWaitTime)
		time.Sleep(waitTime)
	}

	return nil
}

func batchImageURLs(imageURLs []string, batchSize int) [][]string {
	var batches [][]string

	for batchSize < len(imageURLs) {
		imageURLs, batches = imageURLs[batchSize:], append(batches, imageURLs[0:batchSize:batchSize])
	}

	if len(imageURLs) > 0 {
		batches = append(batches, imageURLs)
	}

	return batches
}

func downloadImage(url, downloadDir string) error {
	fileName := filepath.Base(url)
	filePath := filepath.Join(downloadDir, fileName)

	// Check if the file already exists
	if isFileExists(filePath) {
		// File already exists, skip downloading
		return nil
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

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

func getImageFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get image file size: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get image file size, status: %s", resp.Status)
	}

	sizeStr := resp.Header.Get("Content-Length")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse image file size: %v", err)
	}

	return size, nil
}

func parseMaxImageSize(maxSize string) (int64, error) {
	if strings.ToUpper(maxSize) == "MAX" {
		return -1, nil
	}

	size, err := strconv.ParseInt(maxSize, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse max image size: %v", err)
	}

	return size * 1024 * 1024, nil
}

func isImageSizeExceeded(url string, maxSize int64) bool {
	if maxSize == -1 {
		return false
	}

	size, err := getImageFileSize(url)
	if err != nil {
		return true
	}

	return size > maxSize
}

func generateRandomWaitTime(min, max float64) time.Duration {
	waitSeconds := min + rand.Float64()*(max-min)
	waitDuration := time.Duration(waitSeconds * float64(time.Second))

	return waitDuration
}

func isFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
