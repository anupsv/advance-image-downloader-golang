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
)

type Config struct {
	ImageURLFile       string
	DownloadDirectory  string
	BatchSize          int
	MinWaitTime        float64
	MaxWaitTime        float64
}

func main() {
	configFile := "config.txt"

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

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			filename := filepath.Base(url)
			filepath := filepath.Join(config.DownloadDirectory, filename)

			// Skip downloading if the file already exists
			if _, err := os.Stat(filepath); err == nil {
				fmt.Printf("Skipped %s (already exists)\n", filename)
				return
			}

			semaphore <- struct{}{} // Acquire semaphore slot

			err := downloadImage(url, filepath)
			if err != nil {
				fmt.Printf("Error downloading %s: %s\n", filename, err)
			} else {
				fmt.Printf("Downloaded %s\n", filename)
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
	config := &Config{}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "ImageURLFile":
				config.ImageURLFile = value
			case "DownloadDirectory":
				config.DownloadDirectory = value
			case "BatchSize":
				config.BatchSize, err = strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
			case "MinWaitTime":
				config.MinWaitTime, err = strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, err
				}
			case "MaxWaitTime":
				config.MaxWaitTime, err = strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
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

func generateRandomWaitTime(min, max float64) time.Duration {
	waitSeconds := min + rand.Float64()*(max-min)
	waitDuration := time.Duration(waitSeconds * float64(time.Second))

	return waitDuration
}
