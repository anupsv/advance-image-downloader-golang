package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ImageURLFile               string  `mapstructure:"image_url_file"`
	DownloadDirectory          string  `mapstructure:"download_directory"`
	BatchSize                  int     `mapstructure:"batch_size"`
	MinWaitTime                float64 `mapstructure:"min_wait_time"`
	MaxWaitTime                float64 `mapstructure:"max_wait_time"`
	MaxImageSizeMB             string  `mapstructure:"max_image_size_mb"`
	ReplaceDownloadedFileSize  bool    `mapstructure:"replace_downloaded_file_size"`
	SkipIfFileExists           bool    `mapstructure:"skip_if_file_exists"`
}

func main() {
	configFile := "config.yaml"

	config, err := readConfigFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	urls, err := readImageURLs(config.ImageURLFile)
	if err != nil {
		fmt.Printf("Error reading image URLs: %s\n", err)
		return
	}

	totalImages := len(urls)
	imagesLeft := totalImages

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.BatchSize)

	// Create a set of already downloaded images
	alreadyDownloaded := make(map[string]struct{})

	// Populate the set with filenames in the download directory
	files, err := os.ReadDir(config.DownloadDirectory)
	if err != nil {
		fmt.Printf("Error reading download directory: %s\n", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			alreadyDownloaded[file.Name()] = struct{}{}
		}
	}

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			filename := filepath.Base(url)
			filepath := filepath.Join(config.DownloadDirectory, filename)

			// Skip downloading if the file already exists
			if config.SkipIfFileExists && isFileExists(filepath) {
				fmt.Printf("Skipped %s (already exists)\n", filename)
				return
			}

			// Skip size check if max_image_size_mb is set to "MAX"
			if config.MaxImageSizeMB != "MAX" && isImageSizeExceeded(url, config.MaxImageSizeMB) {
				fmt.Printf("Skipped %s (exceeded maximum size)\n", filename)
				return
			}

			semaphore <- struct{}{} // Acquire semaphore slot

			if config.ReplaceDownloadedFileSize {
				if err := replaceDownloadedFile(url, filepath); err != nil {
					fmt.Printf("Error replacing %s: %s\n", filename, err)
				} else {
					fmt.Printf("Replaced %s\n", filename)
				}
			} else {
				if err := downloadImage(url, filepath); err != nil {
					fmt.Printf("Error downloading %s: %s\n", filename, err)
				} else {
					fmt.Printf("Downloaded %s\n", filename)
				}
			}

			<-semaphore // Release semaphore slot
		}(url)

		if len(semaphore) == config.BatchSize {
			wg.Wait() // Wait for the current batch to complete
			imagesLeft -= config.BatchSize
			fmt.Printf("Batch processed. %d images remaining...\n", imagesLeft)

			waitTime := generateRandomWaitTime(config.MinWaitTime, config.MaxWaitTime)
			time.Sleep(waitTime)
		}
	}

	wg.Wait() // Wait for the remaining downloads to complete
	fmt.Println("All images downloaded!")
}

func readImageURLs(file string) ([]string, error) {
	var urls []string

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func readConfigFile(file string) (*Config, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func downloadImage(url, filepath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

func replaceDownloadedFile(url, filepath string) error {
	tempFilepath := filepath + ".temp"

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	tempFile, err := os.Create(tempFilepath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		return err
	}

	// Check if the new file size differs from the existing file size
	existingFileInfo, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	tempFileInfo, err := os.Stat(tempFilepath)
	if err != nil {
		return err
	}

	if existingFileInfo.Size() == tempFileInfo.Size() {
		os.Remove(tempFilepath) // Delete temporary file if sizes match
		return nil
	}

	// Remove the existing file
	err = os.Remove(filepath)
	if err != nil {
		return err
	}

	// Rename the temporary file to the original filename
	err = os.Rename(tempFilepath, filepath)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomWaitTime(min, max float64) time.Duration {
	waitSeconds := min + rand.Float64()*(max-min)
	waitDuration := time.Duration(waitSeconds * float64(time.Second))

	return waitDuration
}

func isImageSizeExceeded(url string, maxSizeMB string) bool {
	if maxSizeMB == "MAX" {
		return false
	}

	maxSize, err := strconv.ParseInt(maxSizeMB, 10, 64)
	if err != nil {
		return false
	}

	response, err := http.Head(url)
	if err != nil {
		return true
	}
	defer response.Body.Close()

	contentLength := response.ContentLength
	if contentLength == -1 {
		return true
	}

	maxSizeBytes := maxSize * 1024 * 1024
	if contentLength > maxSizeBytes {
		return true
	}

	return false
}

func isFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}
