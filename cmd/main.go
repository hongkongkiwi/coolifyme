package main

import (
	"fmt"
	"os"

	"github.com/hongkongkiwi/coolifyme/internal/config"
	"github.com/hongkongkiwi/coolifyme/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	apiToken string
	baseURL  string
	profile  string

	// Version information - set by build process
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "coolifyme",
	Short: "A powerful CLI for the Coolify API",
	Long: `coolifyme is a command-line interface for the Coolify API.
It provides easy access to deploy and manage your applications,
services, and infrastructure through Coolify.

Created by Andy Savage <andy@savage.hk>
Source: https://github.com/hongkongkiwi/coolifyme`,
	Version: getVersionString(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set up logging level based on debug flag
		if viper.GetBool("debug") {
			// Enable debug logging if needed
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add subcommands
	rootCmd.AddCommand(applicationsCmd)
	rootCmd.AddCommand(deployCmd())
	rootCmd.AddCommand(databasesCmd)
	rootCmd.AddCommand(servicesCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(serversCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(teamsCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/coolifyme/config.yaml)")
	rootCmd.PersistentFlags().StringP("server", "s", "", "Coolify server URL")
	rootCmd.PersistentFlags().StringP("token", "t", "", "API token")

	// Bind flags to viper
	viper.BindPFlag("server_url", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("api_token", rootCmd.PersistentFlags().Lookup("token"))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		// Add a config path
		viper.AddConfigPath(home + "/.config/coolifyme")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variables
	viper.SetEnvPrefix("COOLIFYME")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		// Config file loaded successfully
	}
}

// Helper function to create a client from configuration
func createClient() (*client.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with command line flags if provided
	if apiToken != "" {
		cfg.APIToken = apiToken
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if profile != "" {
		cfg.Profile = profile
	}

	return client.New(cfg)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the version number of coolifyme",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("coolifyme v1.0.0")
		fmt.Println("Built with ❤️ for the Coolify community")
	},
}

func getVersionString() string {
	if Version == "dev" {
		return fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate)
	}
	return Version
}
