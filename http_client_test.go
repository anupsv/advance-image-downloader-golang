package main

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"os"
	"testing"
)

type MockStandardHTTPClient struct {
	client HTTPClient
}

func (m *MockStandardHTTPClient) Get(url string) (*http.Response, error) {
	return m.client.Get(url)
}

func TestStandardHTTPClient_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nil, // Replace with the desired response body if needed
	}

	mockHTTPClient := NewMockHTTPClient(ctrl)
	mockHTTPClient.EXPECT().Get(gomock.Any()).Return(mockResponse, nil)

	client := &MockStandardHTTPClient{
		client: mockHTTPClient,
	}

	resp, err := client.Get("https://example.com")
	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestMockHTTPClient_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nil, // Replace with the desired response body if needed
	}

	mockHTTPClient := NewMockHTTPClient(ctrl)
	mockHTTPClient.EXPECT().Get(gomock.Any()).Return(mockResponse, nil)

	resp, err := mockHTTPClient.Get("https://example.com")
	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	exitCode := m.Run()

	// Clean up any resources if needed

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
