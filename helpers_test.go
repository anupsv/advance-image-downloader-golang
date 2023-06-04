package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestIsFileExists(t *testing.T) {
	// Create a temporary file
	file, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Check if the file exists
	exists := isFileExists(file.Name())

	// Verify that the file exists
	if !exists {
		t.Errorf("Expected file to exist, but it does not")
	}
}

func TestIsImageSizeExceeded(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Length header to 10MB
		w.Header().Set("Content-Length", strconv.Itoa(10*1024*1024))
	}))
	defer server.Close()

	// Check if the image size is exceeded
	exceeded := isImageSizeExceeded(server.URL, "5")

	// Verify that the image size is exceeded
	if !exceeded {
		t.Errorf("Expected image size to be exceeded, but it is not")
	}
}

func TestIsImageSizeExceeded_MaxSize(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Length header to 5MB
		w.Header().Set("Content-Length", strconv.Itoa(5*1024*1024))
	}))
	defer server.Close()

	// Check if the image size is exceeded (using "MAX" size)
	exceeded := isImageSizeExceeded(server.URL, "MAX")

	// Verify that the image size is not exceeded
	if exceeded {
		t.Errorf("Expected image size not to be exceeded, but it is")
	}
}

func TestGenerateRandomWaitTime(t *testing.T) {
	// Generate a random wait time
	waitTime := generateRandomWaitTime(1.0, 3.0)

	// Verify that the wait time is within the specified range
	minWait := time.Duration(1.0 * float64(time.Second))
	maxWait := time.Duration(3.0 * float64(time.Second))

	if waitTime < minWait || waitTime > maxWait {
		t.Errorf("Expected wait time to be between %s and %s, but got %s", minWait, maxWait, waitTime)
	}
}
