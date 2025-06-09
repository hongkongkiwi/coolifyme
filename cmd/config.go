package main

import (
	"fmt"

	"github.com/hongkongkiwi/coolifyme/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Manage coolifyme configuration settings",
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  "Set configuration values like API token and base URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = &config.Config{}
		}

		// Get flags
		token, _ := cmd.Flags().GetString("token")
		url, _ := cmd.Flags().GetString("url")
		profile, _ := cmd.Flags().GetString("profile")

		// Update config
		updated := false
		if token != "" {
			cfg.APIToken = token
			updated = true
			fmt.Println("API token updated")
		}
		if url != "" {
			cfg.BaseURL = url
			updated = true
			fmt.Printf("Base URL updated to: %s\n", url)
		}
		if profile != "" {
			cfg.Profile = profile
			updated = true
			fmt.Printf("Profile updated to: %s\n", profile)
		}

		if !updated {
			return fmt.Errorf("no configuration values provided")
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("Configuration saved successfully")
		return nil
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		fmt.Printf("Configuration:\n")
		fmt.Printf("  Profile:   %s\n", cfg.Profile)
		fmt.Printf("  Base URL:  %s\n", cfg.BaseURL)
		if cfg.APIToken != "" {
			fmt.Printf("  API Token: %s...\n", cfg.APIToken[:min(8, len(cfg.APIToken))])
		} else {
			fmt.Printf("  API Token: (not set)\n")
		}

		// Show config file location
		configDir, err := config.GetConfigDir()
		if err == nil {
			fmt.Printf("  Config:    %s/config.yaml\n", configDir)
		}

		return nil
	},
}

// configInitCmd represents the config init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  "Initialize coolifyme configuration with default values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{
			BaseURL: "https://app.coolify.io/api/v1",
			Profile: "default",
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		configDir, _ := config.GetConfigDir()
		fmt.Printf("Configuration initialized at %s/config.yaml\n", configDir)
		fmt.Println("Use 'coolifyme config set --token YOUR_API_TOKEN' to set your API token")

		return nil
	},
}

func init() {
	// Add subcommands to config
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)

	// Flags for config set command
	configSetCmd.Flags().String("token", "", "Set API token")
	configSetCmd.Flags().String("url", "", "Set base URL")
	configSetCmd.Flags().String("profile", "", "Set profile name")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
