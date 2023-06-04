package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadImage(t *testing.T) {
	// Create a temporary directory for test downloads
	downloadDir, err := ioutil.TempDir("/tmp", "test-*-downloads")
	if err != nil {
		t.Fatalf("Failed to create temporary download directory: %v", err)
	}
	defer os.RemoveAll(downloadDir)

	// Set up a mock server to serve the image file
	mockServer := createMockServer()
	defer mockServer.Close()

	// Create a mock image URL
	imageURL := mockServer.URL + "/image.jpg"

	// Download the image
	err = downloadImage(imageURL, downloadDir)
	if err != nil {
		t.Fatalf("Failed to download image: %v", err)
	}

	// Verify that the downloaded file exists
	downloadedFilePath := filepath.Join(downloadDir, "image.jpg")
	_, err = os.Stat(downloadedFilePath)
	if os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %s", downloadedFilePath)
	}
}

func createMockServer() *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a sample image file
		imageFile := filepath.Join("/tmp", "sample_image.jpg")
		_, err := os.OpenFile(imageFile, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		http.ServeFile(w, r, imageFile)
	}))
	return mockServer
}

func createMockImageURLFile(mockURL string) string {
	imageURLFile := filepath.Join("testdata", "image_urls.txt")

	// Create a temporary copy of the mock image URL file
	tempFile, err := ioutil.TempFile("", "mock-image-urls")
	if err != nil {
		log.Fatalf("Failed to create temporary image URL file: %v", err)
	}
	defer tempFile.Close()

	// Read the contents of the mock image URL file
	contents, err := ioutil.ReadFile(imageURLFile)
	if err != nil {
		log.Fatalf("Failed to read mock image URL file: %v", err)
	}

	// Replace the image URLs in the contents with the mock URL
	modifiedContents := strings.Replace(string(contents), "IMAGE_URL", mockURL, -1)

	// Write the modified contents to the temporary file
	_, err = tempFile.WriteString(modifiedContents)
	if err != nil {
		log.Fatalf("Failed to write to temporary image URL file: %v", err)
	}

	return tempFile.Name()
}
