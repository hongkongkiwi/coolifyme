package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// deploymentsCmd represents the deployments command
var deploymentsCmd = &cobra.Command{
	Use:     "deployments",
	Aliases: []string{"deployment", "deploy-status"},
	Short:   "Manage deployments",
	Long:    "Manage Coolify deployments - list current deployments and get deployment details",
}

// deploymentsListCmd represents the deployments list command
var deploymentsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List deployments",
	Long:    "List all currently running deployments in your Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		deployments, err := client.Deployments().ListAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to list deployments: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(deployments, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(deployments) == 0 {
			fmt.Println("No active deployments found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "DEPLOYMENT_UUID\tAPPLICATION\tSTATUS\tCOMMIT")
		_, _ = fmt.Fprintln(w, "---------------\t-----------\t------\t------")

		// Print deployments
		for _, deployment := range deployments {
			deploymentUUID := ""
			appName := ""
			status := ""
			commit := ""

			if deployment.DeploymentUuid != nil {
				deploymentUUID = *deployment.DeploymentUuid
			}
			if deployment.ApplicationName != nil {
				appName = *deployment.ApplicationName
			}
			if deployment.Status != nil {
				status = *deployment.Status
			}
			if deployment.Commit != nil {
				commit = *deployment.Commit
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", deploymentUUID, appName, status, commit)
		}

		return nil
	},
}

// deploymentsGetCmd represents the deployments get command
var deploymentsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get deployment details",
	Long:  "Get detailed information about a specific deployment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		deploymentUUID := args[0]

		deployment, err := client.Deployments().GetByUUID(ctx, deploymentUUID)
		if err != nil {
			return fmt.Errorf("failed to get deployment: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(deployment, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display deployment details in a readable format
		fmt.Printf("Deployment Details:\n")
		fmt.Printf("==================\n")
		if deployment.DeploymentUuid != nil {
			fmt.Printf("UUID:         %s\n", *deployment.DeploymentUuid)
		}
		if deployment.ApplicationName != nil {
			fmt.Printf("Application:  %s\n", *deployment.ApplicationName)
		}
		if deployment.Status != nil {
			fmt.Printf("Status:       %s\n", *deployment.Status)
		}
		if deployment.Commit != nil {
			fmt.Printf("Commit:       %s\n", *deployment.Commit)
		}
		if deployment.CommitMessage != nil {
			fmt.Printf("Message:      %s\n", *deployment.CommitMessage)
		}
		if deployment.ServerName != nil {
			fmt.Printf("Server:       %s\n", *deployment.ServerName)
		}

		return nil
	},
}

// deploymentsListByAppCmd represents the deployments list-by-app command
var deploymentsListByAppCmd = &cobra.Command{
	Use:   "list-by-app <app-uuid>",
	Short: "List deployments by application",
	Long:  "List deployments for a specific application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		appUUID := args[0]

		skip, _ := cmd.Flags().GetInt("skip")
		take, _ := cmd.Flags().GetInt("take")

		deployments, err := client.Deployments().ListWithPagination(ctx, appUUID, skip, take)
		if err != nil {
			return fmt.Errorf("failed to list deployments for application: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(deployments, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(deployments) == 0 {
			fmt.Printf("No deployments found for application %s\n", appUUID)
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tSTATUS")
		_, _ = fmt.Fprintln(w, "----\t----\t------")

		// Print applications (note: the API actually returns Application objects for this endpoint)
		for _, app := range deployments {
			uuid := ""
			name := ""
			status := ""

			if app.Uuid != nil {
				uuid = *app.Uuid
			}
			if app.Name != nil {
				name = *app.Name
			}
			if app.Status != nil {
				status = *app.Status
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", uuid, name, status)
		}

		return nil
	},
}

func init() {
	// Add subcommands to deployments
	deploymentsCmd.AddCommand(deploymentsListCmd)
	deploymentsCmd.AddCommand(deploymentsGetCmd)
	deploymentsCmd.AddCommand(deploymentsListByAppCmd)

	// Flags for list command
	deploymentsListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for get command
	deploymentsGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for list-by-app command
	deploymentsListByAppCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	deploymentsListByAppCmd.Flags().Int("skip", 0, "Number of records to skip (default: 0)")
	deploymentsListByAppCmd.Flags().Int("take", 10, "Number of records to take (default: 10)")
}
