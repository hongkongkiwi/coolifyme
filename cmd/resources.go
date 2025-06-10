package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// resourcesCmd represents the resources command
var resourcesCmd = &cobra.Command{
	Use:     "resources",
	Aliases: []string{"resource", "res"},
	Short:   "Manage resources",
	Long:    "Manage Coolify resources - list all resources across your instance",
}

// resourcesListCmd represents the resources list command
var resourcesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List resources",
	Long:    "List all resources in your Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		result, err := client.Resources().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list resources: %w", err)
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

		// The resources API currently returns a simple string
		fmt.Printf("Resources:\n%s\n", result)
		return nil
	},
}

func init() {
	// Add subcommands to resources
	resourcesCmd.AddCommand(resourcesListCmd)

	// Flags for list command
	resourcesListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
