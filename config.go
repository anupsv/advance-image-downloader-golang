package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ImageURLFile              string
	DownloadDirectory         string
	BatchSize                 int
	MinWaitTime               float64
	MaxWaitTime               float64
	MaxImageSizeMB            string
	ReplaceDownloadedFileSize bool
	SkipIfFileExists          bool
}

func ReadConfigFile(configFilePath string) (*Config, error) {
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

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

	return config, nil
}
