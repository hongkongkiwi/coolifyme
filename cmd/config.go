package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hongkongkiwi/coolifyme/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration and profiles",
	Long: `Manage coolifyme configuration settings and profiles.
	
Profiles allow you to manage multiple Coolify instances or different API tokens
for the same instance (e.g., different users or environments).

Examples:
  # Initialize default configuration
  coolifyme config init

  # Create a new profile for production
  coolifyme config profile create production --token TOKEN --url https://coolify.prod.com/api/v1

  # List all profiles
  coolifyme config profile list

  # Switch to production profile
  coolifyme config profile use production

  # Set global preferences
  coolifyme config set --output json --log-level debug`,
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set global configuration values",
	Long: `Set global configuration values that apply across all profiles.
These settings include output format, logging level, and color preferences.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = &config.Config{
				Profile: "default",
			}
		}

		// Get flags
		outputFormat, _ := cmd.Flags().GetString("output")
		logLevel, _ := cmd.Flags().GetString("log-level")
		colorOutput, _ := cmd.Flags().GetString("color")

		// Update config
		updated := false
		if outputFormat != "" {
			validFormats := []string{"json", "yaml", "table"}
			isValid := false
			for _, format := range validFormats {
				if outputFormat == format {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid output format: %s. Valid options: %s", outputFormat, strings.Join(validFormats, ", "))
			}
			cfg.OutputFormat = outputFormat
			updated = true
			fmt.Printf("✅ Output format set to: %s\n", outputFormat)
		}

		if logLevel != "" {
			validLevels := []string{"debug", "info", "warn", "error"}
			isValid := false
			for _, level := range validLevels {
				if logLevel == level {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid log level: %s. Valid options: %s", logLevel, strings.Join(validLevels, ", "))
			}
			cfg.LogLevel = logLevel
			updated = true
			fmt.Printf("✅ Log level set to: %s\n", logLevel)
		}

		if colorOutput != "" {
			validColors := []string{"auto", "always", "never"}
			isValid := false
			for _, color := range validColors {
				if colorOutput == color {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid color setting: %s. Valid options: %s", colorOutput, strings.Join(validColors, ", "))
			}
			colorBool := colorOutput == "always"
			cfg.ColorOutput = &colorBool
			updated = true
			fmt.Printf("✅ Color output set to: %s\n", colorOutput)
		}

		if !updated {
			return fmt.Errorf("no configuration values provided")
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("📁 Configuration saved successfully")
		return nil
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings and active profile",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		fmt.Printf("📋 Current Configuration\n")
		fmt.Printf("=======================\n")
		fmt.Printf("🔧 Active Profile:  %s\n", cfg.Profile)
		fmt.Printf("🌐 Base URL:        %s\n", cfg.BaseURL)
		if cfg.APIToken != "" {
			fmt.Printf("🔑 API Token:       %s...\n", cfg.APIToken[:minInt(8, len(cfg.APIToken))])
		} else {
			fmt.Printf("🔑 API Token:       (not set)\n")
		}
		fmt.Printf("📄 Output Format:   %s\n", cfg.OutputFormat)
		fmt.Printf("📊 Log Level:       %s\n", cfg.LogLevel)
		if cfg.ColorOutput != nil {
			if *cfg.ColorOutput {
				fmt.Printf("🎨 Color Output:    enabled\n")
			} else {
				fmt.Printf("🎨 Color Output:    disabled\n")
			}
		} else {
			fmt.Printf("🎨 Color Output:    auto\n")
		}

		// Show config file location
		configDir, err := config.GetConfigDir()
		if err == nil {
			fmt.Printf("📁 Config File:     %s/config.yaml\n", configDir)
		}

		return nil
	},
}

// configInitCmd represents the config init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  "Initialize coolifyme configuration with default values and create default profile",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Check if config already exists
		if _, err := config.LoadConfig(); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				return fmt.Errorf("configuration already exists. Use --force to reinitialize")
			}
		}

		cfg := &config.Config{
			BaseURL:      "https://app.coolify.io/api/v1",
			Profile:      "default",
			OutputFormat: "table",
			LogLevel:     "info",
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		configDir, _ := config.GetConfigDir()
		fmt.Printf("✅ Configuration initialized\n")
		fmt.Printf("   📁 Config file: %s/config.yaml\n", configDir)
		fmt.Printf("   🔧 Default profile created: default\n")
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Println("   1. Set your API token: coolifyme config profile set --token YOUR_API_TOKEN")
		fmt.Println("   2. Or create a new profile: coolifyme config profile create production --token TOKEN --url URL")

		return nil
	},
}

// Profile management commands
var configProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long: `Manage configuration profiles for different Coolify instances or environments.
	
Profiles allow you to switch between different Coolify instances, API tokens,
or user accounts without having to reconfigure each time.`,
}

var configProfileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  "List all available configuration profiles",
	RunE: func(cmd *cobra.Command, _ []string) error {
		profiles, defaultProfile, err := config.ListProfiles()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output := map[string]interface{}{
				"default_profile": defaultProfile,
				"profiles":        profiles,
			}
			data, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		if len(profiles) == 0 {
			fmt.Println("❌ No profiles found. Run 'coolifyme config init' to create default profile.")
			return nil
		}

		fmt.Printf("📋 Configuration Profiles\n")
		fmt.Printf("=========================\n")

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "ACTIVE\tNAME\tBASE URL\tAPI TOKEN")
		_, _ = fmt.Fprintln(w, "------\t----\t--------\t---------")

		// Print profiles
		for _, profile := range profiles {
			active := ""
			if profile.Name == defaultProfile {
				active = StatusSuccess
			}

			tokenDisplay := "(not set)"
			if profile.APIToken != "" {
				tokenDisplay = profile.APIToken[:minInt(8, len(profile.APIToken))] + "..."
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				active, profile.Name, profile.BaseURL, tokenDisplay)
		}

		return nil
	},
}

var configProfileCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new profile",
	Long:  "Create a new configuration profile with specified name, API token, and base URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		// Validate profile name
		if err := config.ValidateProfileName(profileName); err != nil {
			return err
		}

		// Get flags
		token, _ := cmd.Flags().GetString("token")
		url, _ := cmd.Flags().GetString("url")

		if token == "" {
			return fmt.Errorf("API token is required (--token)")
		}

		if url == "" {
			url = "https://app.coolify.io/api/v1"
		}

		// Create the profile
		if err := config.CreateProfile(profileName, token, url); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}

		fmt.Printf("✅ Profile '%s' created successfully\n", profileName)
		fmt.Printf("   🌐 Base URL: %s\n", url)
		fmt.Printf("   🔑 API Token: %s...\n", token[:minInt(8, len(token))])
		fmt.Println()
		fmt.Printf("💡 To use this profile: coolifyme config profile use %s\n", profileName)

		return nil
	},
}

var configProfileUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the default profile",
	Long:  "Set the specified profile as the default active profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		if err := config.SetDefaultProfile(profileName); err != nil {
			return fmt.Errorf("failed to set default profile: %w", err)
		}

		fmt.Printf("✅ Default profile set to '%s'\n", profileName)
		return nil
	},
}

var configProfileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Long:  "Delete the specified profile (cannot delete the default profile)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("⚠️  Are you sure you want to delete profile '%s'? This action cannot be undone.\n", profileName)
			fmt.Print("Type 'yes' to confirm: ")
			var confirmation string
			if _, err := fmt.Scanln(&confirmation); err != nil || confirmation != "yes" {
				fmt.Println("❌ Deletion cancelled")
				return nil
			}
		}

		if err := config.DeleteProfile(profileName); err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		fmt.Printf("✅ Profile '%s' deleted successfully\n", profileName)
		return nil
	},
}

var configProfileSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Update current profile settings",
	Long:  "Update API token and base URL for the current profile",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Get flags
		token, _ := cmd.Flags().GetString("token")
		url, _ := cmd.Flags().GetString("url")

		// Update config
		updated := false
		if token != "" {
			cfg.APIToken = token
			updated = true
			fmt.Printf("✅ API token updated for profile '%s'\n", cfg.Profile)
		}
		if url != "" {
			cfg.BaseURL = url
			updated = true
			fmt.Printf("✅ Base URL updated to: %s\n", url)
		}

		if !updated {
			return fmt.Errorf("no configuration values provided")
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("📁 Profile configuration saved successfully")
		return nil
	},
}

func init() {
	// Add subcommands to config
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configProfileCmd)

	// Add profile subcommands
	configProfileCmd.AddCommand(configProfileListCmd)
	configProfileCmd.AddCommand(configProfileCreateCmd)
	configProfileCmd.AddCommand(configProfileUseCmd)
	configProfileCmd.AddCommand(configProfileDeleteCmd)
	configProfileCmd.AddCommand(configProfileSetCmd)

	// Flags for config set command
	configSetCmd.Flags().String("output", "", "Set default output format (json, yaml, table)")
	configSetCmd.Flags().String("log-level", "", "Set log level (debug, info, warn, error)")
	configSetCmd.Flags().String("color", "", "Set color output (auto, always, never)")

	// Flags for config show command
	configShowCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for config init command
	configInitCmd.Flags().Bool("force", false, "Force reinitialize existing configuration")

	// Flags for profile list command
	configProfileListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for profile create command
	configProfileCreateCmd.Flags().String("token", "", "API token (required)")
	configProfileCreateCmd.Flags().String("url", "", "Base URL (default: https://app.coolify.io/api/v1)")
	_ = configProfileCreateCmd.MarkFlagRequired("token")

	// Flags for profile delete command
	configProfileDeleteCmd.Flags().Bool("force", false, "Force delete without confirmation")

	// Flags for profile set command
	configProfileSetCmd.Flags().String("token", "", "Update API token")
	configProfileSetCmd.Flags().String("url", "", "Update base URL")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
