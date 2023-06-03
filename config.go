package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	ImageURLFile               string  `mapstructure:"image_url_file"`
	DownloadDirectory          string  `mapstructure:"download_directory"`
	BatchSize                  int     `mapstructure:"batch_size"`
	MinWaitTime                float64 `mapstructure:"min_wait_time"`
	MaxWaitTime                float64 `mapstructure:"max_wait_time"`
	MaxImageSizeMB             string  `mapstructure:"max_image_size_mb"`
	ReplaceDownloadedFileSize  bool    `mapstructure:"replace_downloaded_file_size"`
	SkipIfFileExists           bool    `mapstructure:"skip_if_file_exists"`
}

func ReadConfigFile(file string) (*Config, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
