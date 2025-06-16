package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/hongkongkiwi/coolifyme/internal/config"
	"github.com/hongkongkiwi/coolifyme/internal/logger"
	"github.com/hongkongkiwi/coolifyme/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	apiToken     string
	baseURL      string
	profile      string
	outputFormat string
	colorOutput  string // "auto", "always", "never"
	verbose      bool
	debug        bool
	quiet        bool

	// Version information - set by build process
	Version = "dev"
	// GitCommit is the git commit hash used to build this version
	GitCommit = "unknown"
	// BuildDate is the date this version was built
	BuildDate = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "coolifyme",
	Short: "A powerful CLI for the Coolify API",
	Long: `coolifyme is a command-line interface for the Coolify API.
It provides easy access to deploy and manage your applications,
services, and infrastructure through Coolify.

Examples:
  # Initialize configuration
  coolifyme config init

  # Set up a profile for your Coolify instance
  coolifyme config profile create production --token YOUR_TOKEN --url https://coolify.yourdomain.com/api/v1

  # Switch between profiles
  coolifyme config profile use production

  # List applications (with debug logging)
  coolifyme --debug applications list

  # Deploy an application
  coolifyme deploy application app-uuid-here

Created by Andy Savage <andy@savage.hk>
Source: https://github.com/hongkongkiwi/coolifyme`,
	Version: getVersionString(),
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		setupLogging()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add subcommands
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(applicationsCmd)
	rootCmd.AddCommand(deployCmd())
	rootCmd.AddCommand(databasesCmd)
	rootCmd.AddCommand(servicesCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(serversCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(teamsCmd)
	rootCmd.AddCommand(privateKeysCmd)
	rootCmd.AddCommand(resourcesCmd)
	rootCmd.AddCommand(deploymentsCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(initInteractiveCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.AddCommand(aliasCmd)

	// Add alias commands at root level for convenience
	rootCmd.AddCommand(deployAppCmd)
	rootCmd.AddCommand(quickStatusCmd)
	rootCmd.AddCommand(quickHealthCmd)
	rootCmd.AddCommand(lsAppsCmd)
	rootCmd.AddCommand(lsServersCmd)
	rootCmd.AddCommand(lsServicesCmd)

	// Add bulk operation commands
	applicationsCmd.AddCommand(appsStartAllCmd)
	applicationsCmd.AddCommand(appsStopAllCmd)
	applicationsCmd.AddCommand(appsRestartAllCmd)
	applicationsCmd.AddCommand(appCreateWizardCmd)
	servicesCmd.AddCommand(servicesDeployAllCmd)
	serversCmd.AddCommand(serverAddWizardCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/coolifyme/config.yaml)")
	rootCmd.PersistentFlags().StringP("server", "s", "", "Coolify server URL")
	rootCmd.PersistentFlags().StringP("token", "t", "", "API token")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "configuration profile to use")
	rootCmd.PersistentFlags().StringP("output", "o", "", "output format (json, yaml, table)")
	rootCmd.PersistentFlags().String("color", "auto", "colorize output (auto, always, never)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output (shows API calls)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet output (errors only)")

	// Bind flags to viper
	_ = viper.BindPFlag("server_url", rootCmd.PersistentFlags().Lookup("server"))
	_ = viper.BindPFlag("api_token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	_ = viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("color_output", rootCmd.PersistentFlags().Lookup("color"))
}

// setupLogging configures the logging system based on flags and config
func setupLogging() {
	var logLevel slog.Level

	// Determine log level based on flags
	if debug {
		logLevel = slog.LevelDebug
	} else if verbose {
		logLevel = slog.LevelInfo
	} else if quiet {
		logLevel = slog.LevelError
	} else {
		// Try to get from config
		cfg, err := config.LoadConfig()
		if err == nil {
			switch cfg.LogLevel {
			case "debug":
				logLevel = slog.LevelDebug
			case "info":
				logLevel = slog.LevelInfo
			case "warn":
				logLevel = slog.LevelWarn
			case "error":
				logLevel = slog.LevelError
			default:
				logLevel = slog.LevelInfo
			}
		} else {
			logLevel = slog.LevelInfo
		}
	}

	logger.SetLevel(logLevel)

	// Set JSON output if explicitly requested or if outputting JSON
	if outputFormat == "json" {
		logger.SetJSONOutput()
	}

	// Configure color output based on setting
	shouldUseColor := shouldEnableColor()
	logger.SetColorOutput(shouldUseColor)

	logger.Debug("Logging initialized",
		"level", logLevel.String(),
		"color", shouldUseColor,
	)
}

// shouldEnableColor determines if color output should be enabled
func shouldEnableColor() bool {
	switch colorOutput {
	case "always":
		return true
	case "never":
		return false
	case "auto":
		fallthrough
	default:
		// Auto-detect based on TTY
		return logger.IsTerminal()
	}
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
		logger.Debug("Configuration file loaded", "file", viper.ConfigFileUsed())
	}

	// Store global flag values for use in other functions
	outputFormat = viper.GetString("output_format")
	colorOutput = viper.GetString("color_output")
	profile = viper.GetString("profile")
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

	logger.Debug("Creating client",
		"baseURL", cfg.BaseURL,
		"profile", cfg.Profile,
		"hasToken", cfg.APIToken != "",
	)

	return client.New(cfg)
}

// Enhanced version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print detailed version information including build details",
	Run: func(cmd *cobra.Command, _ []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if jsonOutput {
			versionInfo := map[string]string{
				"version":   Version,
				"gitCommit": GitCommit,
				"buildDate": BuildDate,
			}
			fmt.Println(mustMarshalJSON(versionInfo))
			return
		}

		fmt.Printf("coolifyme %s\n", getVersionString())
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build date: %s\n", BuildDate)
		fmt.Println()
		fmt.Println("Built with ❤️ for the Coolify community")
		fmt.Println("Source: https://github.com/hongkongkiwi/coolifyme")
	},
}

// Shell completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(coolifyme completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ coolifyme completion bash > /etc/bash_completion.d/coolifyme
  # macOS:
  $ coolifyme completion bash > /usr/local/etc/bash_completion.d/coolifyme

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ coolifyme completion zsh > "${fpath[1]}/_coolifyme"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ coolifyme completion fish | source

  # To load completions for each session, execute once:
  $ coolifyme completion fish > ~/.config/fish/completions/coolifyme.fish

PowerShell:
  PS> coolifyme completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> coolifyme completion powershell > coolifyme.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func getVersionString() string {
	if Version == "dev" {
		return fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate)
	}
	return Version
}

func mustMarshalJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal JSON: %s"}`, err.Error())
	}
	return string(data)
}

func init() {
	// Add version command flags
	versionCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
