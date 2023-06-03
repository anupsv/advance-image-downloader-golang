package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ImageDownloader struct {
	Config   *Config
	URLs     []string
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewImageDownloader(config *Config, urls []string) *ImageDownloader {
	return &ImageDownloader{
		Config:   config,
		URLs:     urls,
		stopChan: make(chan struct{}),
	}
}

func (d *ImageDownloader) Start() {
	totalImages := len(d.URLs)
	imagesLeft := totalImages

	log.Printf("Starting image downloader...")
	log.Printf("Configuration:")
	log.Printf("  - Image URL File: %s", d.Config.ImageURLFile)
	log.Printf("  - Download Directory: %s", d.Config.DownloadDirectory)
	log.Printf("  - Batch Size: %d", d.Config.BatchSize)
	log.Printf("  - Min Wait Time: %.2f", d.Config.MinWaitTime)
	log.Printf("  - Max Wait Time: %.2f", d.Config.MaxWaitTime)
	log.Printf("  - Max Image Size: %s", d.Config.MaxImageSizeMB)
	log.Printf("  - Replace Downloaded File Size: %v", d.Config.ReplaceDownloadedFileSize)
	log.Printf("  - Skip If File Exists: %v", d.Config.SkipIfFileExists)
	log.Printf("Downloading %d images...", totalImages)

	var semaphore = make(chan struct{}, d.Config.BatchSize)

	go func() {
		<-d.stopChan // Wait for the stop signal

		log.Println("Interrupt signal received. Gracefully shutting down...")

		// Wait for the current batch to complete
		d.wg.Wait()
		log.Printf("Batch processed. %d images remaining...", imagesLeft)

		close(semaphore) // Stop accepting new requests
	}()

	for _, url := range d.URLs {
		if len(semaphore) == 0 {
			log.Printf("Batch processed. %d images remaining...", imagesLeft)

			// Wait for the specified wait time between batches
			waitTime := generateRandomWaitTime(d.Config.MinWaitTime, d.Config.MaxWaitTime)
			time.Sleep(waitTime)
		}

		// Skip downloading if the file already exists
		if d.Config.SkipIfFileExists && isFileExists(filepath.Join(d.Config.DownloadDirectory, filepath.Base(url))) {
			log.Printf("Skipped %s (already exists)", filepath.Base(url))
			continue
		}

		// Skip size check if max_image_size_mb is set to "MAX"
		if d.Config.MaxImageSizeMB != "MAX" && isImageSizeExceeded(url, d.Config.MaxImageSizeMB) {
			log.Printf("Skipped %s (exceeded maximum size)", filepath.Base(url))
			continue
		}

		d.wg.Add(1)
		go func(url string) {
			defer d.wg.Done()

			filename := filepath.Base(url)
			filepath := filepath.Join(d.Config.DownloadDirectory, filename)

			semaphore <- struct{}{} // Acquire semaphore slot

			if d.Config.ReplaceDownloadedFileSize {
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
}

func (d *ImageDownloader) Stop() {
	close(d.stopChan) // Send stop signal
	d.wg.Wait()       // Wait for all downloads to finish
}
