// Package config provides configuration management for the coolifyme CLI tool.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	APIToken string `mapstructure:"api_token"`
	BaseURL  string `mapstructure:"base_url"`
	Profile  string `mapstructure:"profile"`
	// Output format preferences
	OutputFormat string `mapstructure:"output_format"` // json, yaml, table
	ColorOutput  *bool  `mapstructure:"color_output"`
	LogLevel     string `mapstructure:"log_level"` // debug, info, warn, error
}

// Profile represents a configuration profile
type Profile struct {
	Name     string `yaml:"name" mapstructure:"name"`
	APIToken string `yaml:"api_token" mapstructure:"api_token"`
	BaseURL  string `yaml:"base_url" mapstructure:"base_url"`
}

// File represents the entire configuration file structure
type File struct {
	DefaultProfile string             `yaml:"default_profile" mapstructure:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles" mapstructure:"profiles"`
	GlobalSettings struct {
		OutputFormat string `yaml:"output_format,omitempty" mapstructure:"output_format"`
		ColorOutput  *bool  `yaml:"color_output,omitempty" mapstructure:"color_output"`
		LogLevel     string `yaml:"log_level,omitempty" mapstructure:"log_level"`
	} `yaml:"global_settings,omitempty" mapstructure:"global_settings"`
}

const (
	// DefaultProfileName represents the default profile name
	DefaultProfileName = "default"
)

var defaultConfig = Config{
	BaseURL:      "https://app.coolify.io/api/v1",
	Profile:      DefaultProfileName,
	OutputFormat: "table",
	LogLevel:     "info",
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// Create a new viper instance for this operation to avoid conflicts
	v := viper.New()

	// Set default values
	v.SetDefault("base_url", defaultConfig.BaseURL)
	v.SetDefault("profile", defaultConfig.Profile)
	v.SetDefault("output_format", defaultConfig.OutputFormat)
	v.SetDefault("log_level", defaultConfig.LogLevel)

	// Environment variable bindings with different prefixes for flexibility
	v.SetEnvPrefix("COOLIFY")
	v.AutomaticEnv()

	// Also support COOLIFYME prefix for backward compatibility
	// Note: viper.BindEnv errors are typically only configuration issues and safe to ignore
	_ = v.BindEnv("api_token", "COOLIFYME_API_TOKEN", "COOLIFY_API_TOKEN")
	_ = v.BindEnv("base_url", "COOLIFYME_BASE_URL", "COOLIFY_BASE_URL", "COOLIFY_URL")
	_ = v.BindEnv("profile", "COOLIFYME_PROFILE", "COOLIFY_PROFILE")
	_ = v.BindEnv("log_level", "COOLIFYME_LOG_LEVEL", "COOLIFY_LOG_LEVEL")

	// Get the active profile name from environment or default
	profileName := v.GetString("profile")

	// Try to load the config file to get the default profile
	configFile, configFileErr := loadConfigFile()
	if configFileErr == nil {
		// If no profile is specified, use the default profile from config file
		if profileName == "" || profileName == DefaultProfileName {
			if configFile.DefaultProfile != "" {
				profileName = configFile.DefaultProfile
			} else {
				profileName = DefaultProfileName
			}
		}
	} else {
		// No config file exists, use defaults
		if profileName == "" {
			profileName = DefaultProfileName
		}
	}

	// Initialize config with defaults
	config := &Config{
		Profile:      profileName,
		OutputFormat: v.GetString("output_format"),
		LogLevel:     v.GetString("log_level"),
		BaseURL:      defaultConfig.BaseURL, // Set default first
	}

	// Check if color output is explicitly set
	if v.IsSet("color_output") {
		colorOutput := v.GetBool("color_output")
		config.ColorOutput = &colorOutput
	}

	// Try to load from profile-specific configuration
	if configFileErr == nil {
		if profileConfig, err := LoadProfile(profileName); err == nil {
			config.APIToken = profileConfig.APIToken
			config.BaseURL = profileConfig.BaseURL
		}

		// Load global settings from config file
		config.OutputFormat = configFile.GlobalSettings.OutputFormat
		if config.OutputFormat == "" {
			config.OutputFormat = v.GetString("output_format")
		}
		config.LogLevel = configFile.GlobalSettings.LogLevel
		if config.LogLevel == "" {
			config.LogLevel = v.GetString("log_level")
		}
		if configFile.GlobalSettings.ColorOutput != nil {
			config.ColorOutput = configFile.GlobalSettings.ColorOutput
		}
	}

	// Command-line flags and environment variables override profile settings
	// Only override if explicitly set via environment variables
	if v.IsSet("api_token") && os.Getenv("COOLIFYME_API_TOKEN") != "" || os.Getenv("COOLIFY_API_TOKEN") != "" {
		if token := v.GetString("api_token"); token != "" {
			config.APIToken = token
		}
	}
	if v.IsSet("base_url") && (os.Getenv("COOLIFYME_BASE_URL") != "" || os.Getenv("COOLIFY_BASE_URL") != "" || os.Getenv("COOLIFY_URL") != "") {
		if url := v.GetString("base_url"); url != "" {
			config.BaseURL = url
		}
	}

	return config, nil
}

// LoadProfile loads a specific profile configuration
func LoadProfile(profileName string) (*Profile, error) {
	configFile, err := loadConfigFile()
	if err != nil {
		return nil, err
	}

	profile, exists := configFile.Profiles[profileName]
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", profileName)
	}

	return &profile, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	configFile, err := loadConfigFile()
	if err != nil {
		// Create new config file if it doesn't exist
		configFile = &File{
			DefaultProfile: config.Profile,
			Profiles:       make(map[string]Profile),
		}
	}

	// If no profile specified, use "default"
	profileName := config.Profile
	if profileName == "" {
		profileName = DefaultProfileName
		config.Profile = profileName
	}

	// Update or create the profile
	profile := Profile{
		Name:     profileName,
		APIToken: config.APIToken,
		BaseURL:  config.BaseURL,
	}

	if configFile.Profiles == nil {
		configFile.Profiles = make(map[string]Profile)
	}
	configFile.Profiles[profileName] = profile

	// Update global settings
	configFile.GlobalSettings.OutputFormat = config.OutputFormat
	configFile.GlobalSettings.ColorOutput = config.ColorOutput
	configFile.GlobalSettings.LogLevel = config.LogLevel

	// Set as default profile if it's the only one or if we're saving the default profile
	if len(configFile.Profiles) == 1 || configFile.DefaultProfile == "" || profileName == DefaultProfileName {
		configFile.DefaultProfile = profileName
	}

	return saveConfigFile(configFile)
}

// CreateProfile creates a new profile
func CreateProfile(name, apiToken, baseURL string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	configFile, err := loadConfigFile()
	if err != nil {
		configFile = &File{
			DefaultProfile: name,
			Profiles:       make(map[string]Profile),
		}
	}

	if configFile.Profiles == nil {
		configFile.Profiles = make(map[string]Profile)
	}

	// Check if profile already exists
	if _, exists := configFile.Profiles[name]; exists {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	// Create the profile
	profile := Profile{
		Name:     name,
		APIToken: apiToken,
		BaseURL:  baseURL,
	}

	if profile.BaseURL == "" {
		profile.BaseURL = defaultConfig.BaseURL
	}

	configFile.Profiles[name] = profile

	// Set as default if it's the first profile
	if len(configFile.Profiles) == 1 {
		configFile.DefaultProfile = name
	}

	return saveConfigFile(configFile)
}

// DeleteProfile deletes a profile
func DeleteProfile(name string) error {
	if name == "default" {
		return fmt.Errorf("cannot delete the default profile")
	}

	configFile, err := loadConfigFile()
	if err != nil {
		return fmt.Errorf("no configuration file found")
	}

	if _, exists := configFile.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	delete(configFile.Profiles, name)

	// If this was the default profile, switch to 'default' or first available
	if configFile.DefaultProfile == name {
		if _, exists := configFile.Profiles["default"]; exists {
			configFile.DefaultProfile = "default"
		} else if len(configFile.Profiles) > 0 {
			// Pick the first available profile
			for profileName := range configFile.Profiles {
				configFile.DefaultProfile = profileName
				break
			}
		} else {
			configFile.DefaultProfile = ""
		}
	}

	return saveConfigFile(configFile)
}

// ListProfiles returns all available profiles
func ListProfiles() ([]Profile, string, error) {
	configFile, err := loadConfigFile()
	if err != nil {
		return nil, "", fmt.Errorf("no configuration file found")
	}

	var profiles []Profile
	for _, profile := range configFile.Profiles {
		profiles = append(profiles, profile)
	}

	return profiles, configFile.DefaultProfile, nil
}

// SetDefaultProfile sets the default profile
func SetDefaultProfile(name string) error {
	configFile, err := loadConfigFile()
	if err != nil {
		return fmt.Errorf("no configuration file found")
	}

	if _, exists := configFile.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	configFile.DefaultProfile = name
	return saveConfigFile(configFile)
}

// loadConfigFile loads the configuration file structure
func loadConfigFile() (*File, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configFile File
	if err := v.Unmarshal(&configFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &configFile, nil
}

// saveConfigFile saves the configuration file structure
func saveConfigFile(configFile *File) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create a new viper instance for writing
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set all the values
	v.Set("default_profile", configFile.DefaultProfile)
	v.Set("profiles", configFile.Profiles)
	if configFile.GlobalSettings.OutputFormat != "" {
		v.Set("global_settings.output_format", configFile.GlobalSettings.OutputFormat)
	}
	if configFile.GlobalSettings.ColorOutput != nil {
		v.Set("global_settings.color_output", *configFile.GlobalSettings.ColorOutput)
	}
	if configFile.GlobalSettings.LogLevel != "" {
		v.Set("global_settings.log_level", configFile.GlobalSettings.LogLevel)
	}

	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigFilePath returns the path to the configuration file
func getConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".config", "coolifyme"), nil
}

// ValidateProfileName validates a profile name
func ValidateProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if strings.Contains(name, " ") {
		return fmt.Errorf("profile name cannot contain spaces")
	}
	if strings.ContainsAny(name, "/:*?\"<>|") {
		return fmt.Errorf("profile name contains invalid characters")
	}
	return nil
}
