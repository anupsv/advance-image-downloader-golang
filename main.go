package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFile := "config.yaml"

	config, err := ReadConfigFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	urls, err := ReadImageURLs(config.ImageURLFile)
	if err != nil {
		log.Fatalf("Error reading image URLs: %s", err)
	}

	downloader := NewImageDownloader(config, urls)
	downloader.Start()

	// Listen for the Ctrl+C signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // Wait for the signal

	downloader.Stop()
	log.Println("Shut down complete.")
}
