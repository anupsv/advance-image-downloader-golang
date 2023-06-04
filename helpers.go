package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ensureDownloadDirectory(downloadDir string) error {
	// Check if the download directory already exists
	_, err := os.Stat(downloadDir)
	if os.IsNotExist(err) {
		// Create the download directory
		err := os.MkdirAll(downloadDir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func readImageURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image URL file: %v", err)
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image URLs from file: %v", err)
	}

	imageURLs := strings.Split(string(contents), "\n")
	var filteredURLs []string
	for _, url := range imageURLs {
		url = strings.TrimSpace(url)
		if url != "" {
			filteredURLs = append(filteredURLs, url)
		}
	}

	return filteredURLs, nil
}

func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %v", err)
	}

	return fileInfo.Size(), nil
}
