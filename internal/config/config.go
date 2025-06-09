// Package config provides configuration management for the coolifyme CLI tool.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	APIToken string `mapstructure:"api_token"`
	BaseURL  string `mapstructure:"base_url"`
	Profile  string `mapstructure:"profile"`
}

var defaultConfig = Config{
	BaseURL: "https://app.coolify.io/api/v1",
	Profile: "default",
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("base_url", defaultConfig.BaseURL)
	viper.SetDefault("profile", defaultConfig.Profile)

	// Set config file name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add config paths
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "coolifyme")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Environment variable bindings
	viper.SetEnvPrefix("COOLIFY")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults
	}

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "coolifyme")
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")

	viper.Set("api_token", config.APIToken)
	viper.Set("base_url", config.BaseURL)
	viper.Set("profile", config.Profile)

	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".config", "coolifyme"), nil
}
