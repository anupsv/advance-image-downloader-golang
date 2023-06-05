package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := ioutil.TempFile("/tmp", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some configuration data to the temporary file
	configData := []byte(`
batch_size: 10
min_wait_time: 0.5
max_wait_time: 2.5
max_image_size_mb: "MAX"
replace_downloaded_file_size: true
skip_if_file_exists: false
`)
	err = ioutil.WriteFile(tempFile.Name(), configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write config data to temporary file: %v", err)
	}

	// Load the configuration
	err = loadConfig(tempFile.Name())
	if err != nil {
		t.Errorf("Failed to load configuration: %v", err)
	}

	// Check the loaded configuration values
	expectedBatchSize := 10
	if viper.GetInt("batch_size") != expectedBatchSize {
		t.Errorf("Expected batch size to be %d, but got %d", expectedBatchSize, viper.GetInt("batch_size"))
	}

	expectedMinWaitTime := 0.5
	if viper.GetFloat64("min_wait_time") != expectedMinWaitTime {
		t.Errorf("Expected min wait time to be %.2f, but got %.2f", expectedMinWaitTime, viper.GetFloat64("min_wait_time"))
	}

	expectedMaxWaitTime := 2.5
	if viper.GetFloat64("max_wait_time") != expectedMaxWaitTime {
		t.Errorf("Expected max wait time to be %.2f, but got %.2f", expectedMaxWaitTime, viper.GetFloat64("max_wait_time"))
	}

	expectedMaxImageSize := "MAX"
	if viper.GetString("max_image_size_mb") != expectedMaxImageSize {
		t.Errorf("Expected max image size to be %s, but got %s", expectedMaxImageSize, viper.GetString("max_image_size_mb"))
	}

	expectedReplaceFileSize := true
	if viper.GetBool("replace_downloaded_file_size") != expectedReplaceFileSize {
		t.Errorf("Expected replace downloaded file size to be %v, but got %v", expectedReplaceFileSize, viper.GetBool("replace_downloaded_file_size"))
	}

	expectedSkipIfExists := false
	if viper.GetBool("skip_if_file_exists") != expectedSkipIfExists {
		t.Errorf("Expected skip if file exists to be %v, but got %v", expectedSkipIfExists, viper.GetBool("skip_if_file_exists"))
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {

	// Load the configuration
	err := loadConfig("/tmp/nonexistent.yaml")
	if err == nil {
		t.Errorf("Expected error when loading invalid config file, but got nil")
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear any previously set configuration values
	viper.Reset()

	// Create a temporary config file
	tempFile, err := ioutil.TempFile("/tmp", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some configuration data to the temporary file
	configData := []byte(`
batch_size: 2
min_wait_time: 0.8
max_wait_time: 3.0
max_image_size_mb: "MAX"
replace_downloaded_file_size: false
skip_if_file_exists: true
`)
	err = ioutil.WriteFile(tempFile.Name(), configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write config data to temporary file: %v", err)
	}

	// Load the configuration
	err = loadConfig(tempFile.Name())
	if err != nil {
		t.Errorf("Failed to load configuration: %v", err)
	}

	// Check the default configuration values
	defaultBatchSize := 2
	if viper.GetInt("batch_size") != defaultBatchSize {
		t.Errorf("Expected default batch size to be %d, but got %d", defaultBatchSize, viper.GetInt("batch_size"))
	}

	defaultMinWaitTime := 0.8
	if viper.GetFloat64("min_wait_time") != defaultMinWaitTime {
		t.Errorf("Expected default min wait time to be %.2f, but got %.2f", defaultMinWaitTime, viper.GetFloat64("min_wait_time"))
	}

	defaultMaxWaitTime := 3.0
	if viper.GetFloat64("max_wait_time") != defaultMaxWaitTime {
		t.Errorf("Expected default max wait time to be %.2f, but got %.2f", defaultMaxWaitTime, viper.GetFloat64("max_wait_time"))
	}

	defaultMaxImageSize := "MAX"
	if viper.GetString("max_image_size_mb") != defaultMaxImageSize {
		t.Errorf("Expected default max image size to be %s, but got %s", defaultMaxImageSize, viper.GetString("max_image_size_mb"))
	}

	defaultReplaceFileSize := false
	if viper.GetBool("replace_downloaded_file_size") != defaultReplaceFileSize {
		t.Errorf("Expected default replace downloaded file size to be %v, but got %v", defaultReplaceFileSize, viper.GetBool("replace_downloaded_file_size"))
	}

	defaultSkipIfExists := true
	if viper.GetBool("skip_if_file_exists") != defaultSkipIfExists {
		t.Errorf("Expected default skip if file exists to be %v, but got %v", defaultSkipIfExists, viper.GetBool("skip_if_file_exists"))
	}
}
