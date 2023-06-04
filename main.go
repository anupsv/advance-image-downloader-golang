package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load the configuration
	err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print the current configuration
	printConfig()

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Start the image downloader
	go func() {
		err := startImageDownloader()
		if err != nil {
			log.Fatalf("Image downloader failed: %v", err)
		}
	}()

	// Wait for the termination signal
	<-signalCh
	log.Println("Received termination signal. Shutting down...")
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	viper.SetDefault("batch_size", 2)
	viper.SetDefault("min_wait_time", 0.8)
	viper.SetDefault("max_wait_time", 3.0)
	viper.SetDefault("max_image_size_mb", "MAX")
	viper.SetDefault("replace_downloaded_file_size", false)
	viper.SetDefault("skip_if_file_exists", true)

	return nil
}

func printConfig() {
	log.Println("Current Configuration:")
	log.Println("======================")
	log.Printf("Batch Size: %d", viper.GetInt("batch_size"))
	log.Printf("Min Wait Time: %.2f", viper.GetFloat64("min_wait_time"))
	log.Printf("Max Wait Time: %.2f", viper.GetFloat64("max_wait_time"))
	log.Printf("Max Image Size: %s", viper.GetString("max_image_size_mb"))
	log.Printf("Replace Downloaded File Size: %v", viper.GetBool("replace_downloaded_file_size"))
	log.Printf("Skip If File Exists: %v", viper.GetBool("skip_if_file_exists"))
	log.Println("======================")
}

func startImageDownloader() error {
	config := &Config{
		ImageURLFile:              viper.GetString("image_url_file"),
		DownloadDirectory:         viper.GetString("download_directory"),
		BatchSize:                 viper.GetInt("batch_size"),
		MinWaitTime:               viper.GetFloat64("min_wait_time"),
		MaxWaitTime:               viper.GetFloat64("max_wait_time"),
		MaxImageSizeMB:            viper.GetString("max_image_size_mb"),
		ReplaceDownloadedFileSize: viper.GetBool("replace_downloaded_file_size"),
		SkipIfFileExists:          viper.GetBool("skip_if_file_exists"),
	}

	err := downloadImages(config)
	if err != nil {
		return fmt.Errorf("failed to download images: %v", err)
	}

	return nil
}
