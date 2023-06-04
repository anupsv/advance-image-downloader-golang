package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDownloadDirectory_DirectoryExists(t *testing.T) {
	// Create a temporary download directory
	tempDir := t.TempDir()

	// Ensure the download directory
	err := ensureDownloadDirectory(tempDir)
	assert.NoError(t, err)
}

func TestReadImageURLsFromFile(t *testing.T) {
	// Create a temporary file with image URLs
	tempFile := createTempFile(t, []byte("https://example.com/image1.jpg\nhttps://example.com/image2.jpg\nhttps://example.com/image3.jpg"))

	// Read the image URLs from the file
	imageURLs, err := readImageURLsFromFile(tempFile)
	assert.NoError(t, err)

	// Assert the image URLs
	assert.Equal(t, []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg", "https://example.com/image3.jpg"}, imageURLs)
}

func TestReadImageURLsFromFile_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tempFile := createTempFile(t, []byte{})

	// Read the image URLs from the file
	imageURLs, err := readImageURLsFromFile(tempFile)
	assert.NoError(t, err)

	// Assert the image URLs
	assert.Empty(t, imageURLs)
}

func TestReadImageURLsFromFile_NonexistentFile(t *testing.T) {
	// Read from a non-existent file
	imageURLs, err := readImageURLsFromFile("nonexistent.txt")

	// Assert the error and image URLs
	assert.Error(t, err)
	assert.Nil(t, imageURLs)
}

func TestBatchImageURLs(t *testing.T) {
	// Create a list of image URLs
	imageURLs := []string{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
		"https://example.com/image3.jpg",
		"https://example.com/image4.jpg",
		"https://example.com/image5.jpg",
		"https://example.com/image6.jpg",
		"https://example.com/image7.jpg",
		"https://example.com/image8.jpg",
		"https://example.com/image9.jpg",
		"https://example.com/image10.jpg",
		"https://example.com/image11.jpg",
	}

	// Batch the image URLs
	batches := batchImageURLs(imageURLs, 5)

	// Assert the batches
	assert.Equal(t, [][]string{
		{"https://example.com/image1.jpg", "https://example.com/image2.jpg", "https://example.com/image3.jpg", "https://example.com/image4.jpg", "https://example.com/image5.jpg"},
		{"https://example.com/image6.jpg", "https://example.com/image7.jpg", "https://example.com/image8.jpg", "https://example.com/image9.jpg", "https://example.com/image10.jpg"},
		{"https://example.com/image11.jpg"},
	}, batches)
}

func TestBatchImageURLs_EmptyURLs(t *testing.T) {
	// Create an empty list of image URLs
	imageURLs := []string{}

	// Batch the image URLs
	batches := batchImageURLs(imageURLs, 5)

	// Assert the batches
	assert.Empty(t, batches)
}

func TestDownloadImage_FileExists(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a dummy file in the directory
	dummyFile := filepath.Join(tempDir, "dummy.txt")
	err := ioutil.WriteFile(dummyFile, []byte("dummy"), 0644)
	assert.NoError(t, err)

	// Download an image
	err = downloadImage("https://via.placeholder.com/150", tempDir)
	assert.NoError(t, err)

	// Check if the image file exists
	imagePath := filepath.Join(tempDir, "150")
	_, err = os.Stat(imagePath)
	assert.NoError(t, err)

	// Remove the image file
	err = os.Remove(imagePath)
	assert.NoError(t, err)

	// Remove the dummy file
	err = os.Remove(dummyFile)
	assert.NoError(t, err)
}

func TestGetImageFileSize(t *testing.T) {
	// Get the size of a remote image
	size, err := getImageFileSize("https://via.placeholder.com/150")
	assert.NoError(t, err)
	assert.NotZero(t, size)
}

func TestGetImageFileSize_InvalidURL(t *testing.T) {
	// Get the size of an invalid image URL
	size, err := getImageFileSize("https://example.com/invalid.jpg")
	assert.Error(t, err)
	assert.Zero(t, size)
}

func TestParseMaxImageSize(t *testing.T) {
	// Parse max image size in MB
	size, err := parseMaxImageSize("5")
	assert.NoError(t, err)
	assert.Equal(t, int64(5*1024*1024), size)
}

func TestParseMaxImageSize_Max(t *testing.T) {
	// Parse max image size with "MAX"
	size, err := parseMaxImageSize("MAX")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), size)
}

func TestParseMaxImageSize_InvalidSize(t *testing.T) {
	// Parse invalid max image size
	size, err := parseMaxImageSize("invalid")
	assert.Error(t, err)
	assert.Zero(t, size)
}

func TestIsImageSizeExceeded(t *testing.T) {
	// Check if image size is exceeded
	size := int64(5 * 1024 * 1024)
	exceeded := isImageSizeExceeded("https://via.placeholder.com/150", size)
	assert.False(t, exceeded)
}

func TestIsImageSizeExceeded_InvalidURL(t *testing.T) {
	// Check if image size is exceeded with an invalid URL
	size := int64(5 * 1024 * 1024)
	exceeded := isImageSizeExceeded("https://example.com/invalid.jpg", size)
	assert.True(t, exceeded)
}

func TestGenerateRandomWaitTime(t *testing.T) {
	// Generate random wait time
	waitTime := generateRandomWaitTime(0.8, 3.0)
	assert.GreaterOrEqual(t, waitTime.Seconds(), 0.8)
	assert.LessOrEqual(t, waitTime.Seconds(), 3.0)
}

func createTempFile(t *testing.T, content []byte) string {
	tempFile, err := ioutil.TempFile("", "testfile")
	assert.NoError(t, err)

	_, err = tempFile.Write(content)
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	return tempFile.Name()
}
