package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	configDir := filepath.Join(home, "/.config/flyctl")
	configFile := filepath.Join(configDir, "config.yaml")

	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Set default config values
	viper.SetDefault("adults", "1")
	viper.SetDefault("currency", "CAD")
	viper.SetDefault("amadeus_api_key", "please fill in")
	viper.SetDefault("amadeus_api_secret", "please fill in")
	viper.SetDefault("rapid_google_api_key", "please fill in")

	var fileLookupError viper.ConfigFileNotFoundError
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			if err := os.MkdirAll(configDir, 0o755); err != nil {
				log.Printf("Failed for some reason, %s", err.Error())
				return fmt.Errorf("create config dir %q: %w", configDir, err)
			}
			if err := viper.SafeWriteConfigAs(configFile); err != nil {
				log.Printf("Failed for some reason, %s", err.Error())
				return fmt.Errorf("write default config to %q: %w", configFile, err)
			}
			return nil
		} else {
			return fmt.Errorf("failed to read config from %q: %w", configFile, err)

		}
	}
	return nil
}
