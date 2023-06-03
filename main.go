package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
		log.Fatalf("Error reading config file: %s", err)
	}

	urls, err := readImageURLs(config.ImageURLFile)
	if err != nil {
		log.Fatalf("Error reading image URLs: %s", err)
	}

	totalImages := len(urls)
	imagesLeft := totalImages

	log.Printf("Starting image downloader...")
	log.Printf("Configuration:")
	log.Printf("  - Image URL File: %s", config.ImageURLFile)
	log.Printf("  - Download Directory: %s", config.DownloadDirectory)
	log.Printf("  - Batch Size: %d", config.BatchSize)
	log.Printf("  - Min Wait Time: %.2f", config.MinWaitTime)
	log.Printf("  - Max Wait Time: %.2f", config.MaxWaitTime)
	log.Printf("  - Max Image Size: %s", config.MaxImageSizeMB)
	log.Printf("  - Replace Downloaded File Size: %v", config.ReplaceDownloadedFileSize)
	log.Printf("  - Skip If File Exists: %v", config.SkipIfFileExists)
	log.Printf("Downloading %d images...", totalImages)

	// Create a channel to receive the Ctrl+C signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.BatchSize)

	// Create a set of already downloaded images
	alreadyDownloaded := make(map[string]struct{})

	// Populate the set with filenames in the download directory
	files, err := os.ReadDir(config.DownloadDirectory)
	if err != nil {
		log.Fatalf("Error reading download directory: %s", err)
	}
	for _, file := range files {
		if !file.IsDir() {
			alreadyDownloaded[file.Name()] = struct{}{}
		}
	}

	go func() {
		<-stop // Wait for the Ctrl+C signal

		log.Println("Interrupt signal received. Gracefully shutting down...")

		// Wait for the current batch to complete
		wg.Wait()
		log.Printf("Batch processed. %d images remaining...", imagesLeft)

		// Stop accepting new requests by closing the semaphore
		close(semaphore)

		// Wait for all remaining downloads to finish
		wg.Wait()

		log.Println("Shut down complete.")
		os.Exit(0)
	}()

	for _, url := range urls {
		if len(semaphore) == 0 {
			log.Printf("Batch processed. %d images remaining...", imagesLeft)

			// Wait for the specified wait time between batches
			waitTime := generateRandomWaitTime(config.MinWaitTime, config.MaxWaitTime)
			time.Sleep(waitTime)
		}

		// Skip downloading if the file already exists
		if config.SkipIfFileExists && isFileExists(filepath.Join(config.DownloadDirectory, filepath.Base(url))) {
			log.Printf("Skipped %s (already exists)", filepath.Base(url))
			continue
		}

		// Skip size check if max_image_size_mb is set to "MAX"
		if config.MaxImageSizeMB != "MAX" && isImageSizeExceeded(url, config.MaxImageSizeMB) {
			log.Printf("Skipped %s (exceeded maximum size)", filepath.Base(url))
			continue
		}

		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			filename := filepath.Base(url)
			filepath := filepath.Join(config.DownloadDirectory, filename)

			semaphore <- struct{}{} // Acquire semaphore slot

			if config.ReplaceDownloadedFileSize {
				if err := replaceDownloadedFile(url, filepath); err != nil {
					log.Printf("[Goroutine %d] Error replacing %s: %s", getGoroutineID(), filename, err)
				} else {
					log.Printf("[Goroutine %d] Replaced %s", getGoroutineID(), filename)
				}
			} else {
				if err := downloadImage(url, filepath); err != nil {
					log.Printf("[Goroutine %d] Error downloading %s: %s", getGoroutineID(), filename, err)
				} else {
					log.Printf("[Goroutine %d] Downloaded %s", getGoroutineID(), filename)
				}
			}

			<-semaphore // Release semaphore slot
			imagesLeft--
		}(url)
	}

	wg.Wait() // Wait for the remaining downloads to complete
	log.Println("All images downloaded!")
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

// getGoroutineID returns the ID of the current goroutine.
func getGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.Atoi(idField)
	return id
}
