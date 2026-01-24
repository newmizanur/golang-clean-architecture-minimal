package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json and environment variables
// Config file is optional; env vars can override or replace file-based config
// Example: APP_DATABASE_HOST will override database.host from config.json
func NewViper() *viper.Viper {
	config := viper.New()

	// Load from config file
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")

	// Allow environment variable overrides
	config.AutomaticEnv()
	config.SetEnvPrefix("APP")

	err := config.ReadInConfig()
	if err != nil {
		// Config file is optional - can run with env vars only
		fmt.Printf("Warning: config file not found (%v). Using environment variables only\n", err)
	}

	return config
}
