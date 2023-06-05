package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Downloader interface {
	DownloadImage(url, downloadDir string) error
}

type URLReader interface {
	ReadImageURLsFromFile(filePath string) ([]string, error)
}

type ImageSizeChecker interface {
	IsImageSizeExceeded(url string, maxSize int64) bool
}

type FileChecker interface {
	IsFileExists(filePath string) bool
}

type FileSizeGetter interface {
	GetImageFileSize(url string) (int64, error)
}

type WaitTimeGenerator interface {
	GenerateRandomWaitTime(min, max float64) time.Duration
}

type Helper struct {
	Downloader        Downloader
	URLReader         URLReader
	ImageSizeChecker  ImageSizeChecker
	FileChecker       FileChecker
	FileSizeGetter    FileSizeGetter
	WaitTimeGenerator WaitTimeGenerator
}

func NewHelper(
	downloader Downloader,
	urlReader URLReader,
	imageSizeChecker ImageSizeChecker,
	fileChecker FileChecker,
	fileSizeGetter FileSizeGetter,
	waitTimeGenerator WaitTimeGenerator,
) *Helper {
	return &Helper{
		Downloader:        downloader,
		URLReader:         urlReader,
		ImageSizeChecker:  imageSizeChecker,
		FileChecker:       fileChecker,
		FileSizeGetter:    fileSizeGetter,
		WaitTimeGenerator: waitTimeGenerator,
	}
}

func (h *Helper) DownloadImages(config *Config) error {
	imageURLs, err := h.URLReader.ReadImageURLsFromFile(config.ImageURLFile)
	if err != nil {
		return fmt.Errorf("failed to read image URLs from file: %v", err)
	}

	err = h.ensureDownloadDirectory(config.DownloadDirectory)
	if err != nil {
		return fmt.Errorf("failed to ensure download directory: %v", err)
	}

	batches := batchImageURLs(imageURLs, config.BatchSize)
	for _, batch := range batches {
		err := h.downloadBatch(batch, config.DownloadDirectory, config.MaxImageSizeMB)
		if err != nil {
			return fmt.Errorf("failed to download image batch: %v", err)
		}

		waitTime := h.WaitTimeGenerator.GenerateRandomWaitTime(config.MinWaitTime, config.MaxWaitTime)
		time.Sleep(waitTime)
	}

	return nil
}

func (h *Helper) ensureDownloadDirectory(directory string) error {
	if !h.FileChecker.IsFileExists(directory) {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return fmt.Errorf("failed to create download directory: %v", err)
		}
	}

	return nil
}

func (h *Helper) downloadBatch(batch []string, downloadDir string, maxImageSizeMB string) error {
	for _, url := range batch {
		if h.FileChecker.IsFileExists(url) {
			continue
		}

		m, err := strconv.ParseInt(maxImageSizeMB, 10, 64)

		if err != nil {
			return fmt.Errorf("failed to parse maxImageSizeMB: %v", err)
		}

		if !h.ImageSizeChecker.IsImageSizeExceeded(url, m) {
			err := h.Downloader.DownloadImage(url, downloadDir)
			if err != nil {
				return fmt.Errorf("failed to download image: %v", err)
			}
		}
	}

	return nil
}

func (h *Helper) ReadImageURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image URL file: %v", err)
	}
	defer file.Close()

	var imageURLs []string
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read image URLs from file: %v", err)
		}
		imageURLs = append(imageURLs, string(buf[:n]))
	}

	return imageURLs, nil
}

func (h *Helper) IsFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (h *Helper) GenerateRandomWaitTime(min, max float64) time.Duration {
	waitSeconds := min + rand.Float64()*(max-min)
	return time.Duration(waitSeconds * float64(time.Second))
}

func NewImageDownloader(httpClient HTTPClient, fileChecker FileChecker) *ImageDownloader {
	return &ImageDownloader{
		HTTPClient:  httpClient,
		FileChecker: fileChecker,
	}
}

func NewDefaultFileChecker() *DefaultFileChecker {
	return &DefaultFileChecker{}
}

type DefaultFileChecker struct{}

func (f *DefaultFileChecker) IsFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func NewDefaultURLReader() *DefaultURLReader {
	return &DefaultURLReader{}
}

type DefaultURLReader struct{}

func (r *DefaultURLReader) ReadImageURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image URL file: %v", err)
	}
	defer file.Close()

	var imageURLs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		imageURLs = append(imageURLs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read image URLs from file: %v", err)
	}

	return imageURLs, nil
}

func NewDefaultImageSizeChecker() *DefaultImageSizeChecker {
	return &DefaultImageSizeChecker{}
}

type DefaultImageSizeChecker struct {
	FileSizeGetter FileSizeGetter
}

func (c *DefaultImageSizeChecker) IsImageSizeExceeded(url string, maxSize int64) bool {
	if maxSize == -1 {
		return false
	}

	size, err := c.FileSizeGetter.GetImageFileSize(url)
	if err != nil {
		return true
	}

	return size > maxSize
}

func NewDefaultFileSizeGetter() *DefaultFileSizeGetter {
	return &DefaultFileSizeGetter{}
}

type DefaultFileSizeGetter struct{}

func (f *DefaultFileSizeGetter) GetImageFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get image file size: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get image file size, status: %s", resp.Status)
	}

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse image file size: %v", err)
	}

	return size, nil
}

func NewDefaultWaitTimeGenerator() *DefaultWaitTimeGenerator {
	return &DefaultWaitTimeGenerator{}
}

type DefaultWaitTimeGenerator struct{}

func (g *DefaultWaitTimeGenerator) GenerateRandomWaitTime(min, max float64) time.Duration {
	waitSeconds := min + rand.Float64()*(max-min)
	return time.Duration(waitSeconds * float64(time.Second))
}
