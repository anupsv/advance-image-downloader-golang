package main

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ReadImageURLs(file string) ([]string, error) {
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
