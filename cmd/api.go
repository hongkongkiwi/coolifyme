package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Manage API settings",
	Long:  "Manage Coolify API settings - check version, enable/disable API access, and health status",
}

// apiVersionCmd represents the api version command
var apiVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get API version",
	Long:  "Get the current Coolify API version information",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		version, err := client.System().Version(ctx)
		if err != nil {
			return fmt.Errorf("failed to get API version: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(version, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		fmt.Printf("📋 Coolify API Version Information\n")
		fmt.Printf("==================================\n")
		fmt.Printf("Version: %s\n", version)
		return nil
	},
}

// apiEnableCmd represents the api enable command
var apiEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable API access",
	Long:  "Enable API access for the current Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		result, err := client.System().EnableAPI(ctx)
		if err != nil {
			return fmt.Errorf("failed to enable API: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		fmt.Printf("✅ API access enabled successfully\n")
		fmt.Printf("   📝 Response: %s\n", result)
		return nil
	},
}

// apiDisableCmd represents the api disable command
var apiDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable API access",
	Long:  "Disable API access for the current Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		result, err := client.System().DisableAPI(ctx)
		if err != nil {
			return fmt.Errorf("failed to disable API: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		fmt.Printf("✅ API access disabled successfully\n")
		fmt.Printf("   📝 Response: %s\n", result)
		return nil
	},
}

// apiHealthcheckCmd represents the api healthcheck command
var apiHealthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check API health",
	Long:  "Check the health status of the Coolify API",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		health, err := client.System().Healthcheck(ctx)
		if err != nil {
			return fmt.Errorf("failed to check API health: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(health, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		fmt.Printf("🩺 Coolify API Health Status\n")
		fmt.Printf("===========================\n")
		fmt.Printf("Status: %s\n", health)
		return nil
	},
}

func init() {
	// Add subcommands to api
	apiCmd.AddCommand(apiVersionCmd)
	apiCmd.AddCommand(apiEnableCmd)
	apiCmd.AddCommand(apiDisableCmd)
	apiCmd.AddCommand(apiHealthcheckCmd)

	// Flags for all commands
	apiVersionCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	apiEnableCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	apiDisableCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	apiHealthcheckCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
