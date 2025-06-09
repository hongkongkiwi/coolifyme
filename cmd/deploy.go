package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func deployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"deployment", "deployments"},
		Short:   "Deploy applications and services",
		Long:    "Trigger deployments for applications and services in Coolify, and manage deployment history",
	}

	cmd.AddCommand(deployApplicationCmd())
	cmd.AddCommand(deployServiceCmd())
	cmd.AddCommand(deployListCmd())
	cmd.AddCommand(deployListAllCmd())
	cmd.AddCommand(deployGetCmd())

	return cmd
}

func deployApplicationCmd() *cobra.Command {
	var force bool
	var branch string
	var pr int

	cmd := &cobra.Command{
		Use:   "application [uuid]",
		Short: "Deploy an application",
		Long:  "Trigger a deployment for the specified application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			applicationUUID := args[0]
			ctx := context.Background()

			fmt.Printf("🚀 Starting application deployment for %s\n", applicationUUID)
			if branch != "" {
				fmt.Printf("   Branch: %s\n", branch)
			}
			if pr > 0 {
				fmt.Printf("   Pull Request: #%d\n", pr)
			}
			if force {
				fmt.Printf("   Force deployment: enabled\n")
			}

			if branch != "" && pr > 0 {
				return fmt.Errorf("cannot specify both branch and PR - they are mutually exclusive")
			}

			// Use the enhanced client method that supports PR deployments
			options := &client.DeployApplicationOptions{
				Force:  force,
				Branch: branch,
			}
			if pr > 0 {
				options.PR = &pr
			}

			err = client.Deployments().DeployApplicationWithOptions(ctx, applicationUUID, options)
			if err != nil {
				return fmt.Errorf("failed to deploy application: %w", err)
			}

			fmt.Printf("✅ Application deployment triggered successfully for %s\n", applicationUUID)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deployment even if one is already running")
	cmd.Flags().StringVarP(&branch, "branch", "b", "", "Deploy from specific branch/tag")
	cmd.Flags().IntVarP(&pr, "pr", "p", 0, "Deploy specific Pull Request (cannot be used with --branch)")

	return cmd
}

func deployServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service [uuid]",
		Short: "Deploy a service",
		Long:  "Trigger a deployment for the specified service (services are deployed by starting them)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			serviceUUID := args[0]
			ctx := context.Background()

			fmt.Printf("🚀 Starting service deployment for %s\n", serviceUUID)

			// Use the deployment client's method
			err = client.Deployments().DeployService(ctx, serviceUUID)
			if err != nil {
				return fmt.Errorf("failed to deploy service: %w", err)
			}

			fmt.Printf("✅ Service deployment triggered successfully for %s\n", serviceUUID)

			return nil
		},
	}

	return cmd
}

func deployListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [app-uuid]",
		Short: "List deployments for an application",
		Long:  "List deployment history for a specific application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			appUUID := args[0]
			ctx := context.Background()

			deployments, err := client.Deployments().List(ctx, appUUID)
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
				fmt.Printf("No deployments found for application %s\n", appUUID)
				return nil
			}

			// Create a tabwriter for nicely formatted output
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer w.Flush()

			// Print header
			fmt.Fprintln(w, "UUID\tNAME\tSTATUS\tBRANCH\tDOMAINS")
			fmt.Fprintln(w, "----\t----\t------\t------\t-------")

			// Print deployments (Note: this returns Application objects, not ApplicationDeploymentQueue)
			for _, deployment := range deployments {
				uuid := ""
				name := ""
				status := ""
				branch := ""
				domains := ""

				if deployment.Uuid != nil {
					uuid = *deployment.Uuid
				}
				if deployment.Name != nil {
					name = *deployment.Name
				}
				if deployment.Status != nil {
					status = *deployment.Status
				}
				if deployment.GitBranch != nil {
					branch = *deployment.GitBranch
				}
				if deployment.Fqdn != nil {
					domains = *deployment.Fqdn
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					uuid, name, status, branch, domains)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	return cmd
}

func deployListAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-all",
		Aliases: []string{"all"},
		Short:   "List all running deployments",
		Long:    "List all currently running deployments across all applications",
		RunE: func(cmd *cobra.Command, args []string) error {
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
				fmt.Println("No running deployments found")
				return nil
			}

			// Create a tabwriter for nicely formatted output
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer w.Flush()

			// Print header
			fmt.Fprintln(w, "ID\tAPP NAME\tSTATUS\tCREATED\tSERVER")
			fmt.Fprintln(w, "--\t--------\t------\t-------\t------")

			// Print deployments - using correct ApplicationDeploymentQueue fields
			for _, deployment := range deployments {
				id := ""
				appName := ""
				status := ""
				created := ""
				server := ""

				if deployment.Id != nil {
					id = fmt.Sprintf("%d", *deployment.Id)
				}
				if deployment.ApplicationName != nil {
					appName = *deployment.ApplicationName
				}
				if deployment.Status != nil {
					status = *deployment.Status
				}
				if deployment.CreatedAt != nil {
					created = *deployment.CreatedAt
				}
				if deployment.ServerName != nil {
					server = *deployment.ServerName
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					id, appName, status, created, server)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	return cmd
}

func deployGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [deployment-uuid]",
		Short: "Get deployment details",
		Long:  "Get detailed information about a specific deployment by UUID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			deploymentUUID := args[0]
			ctx := context.Background()

			deployment, err := client.Deployments().GetByUuid(ctx, deploymentUUID)
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

			// Display deployment details in a readable format using correct ApplicationDeploymentQueue fields
			fmt.Printf("Deployment Details:\n")
			fmt.Printf("==================\n")
			if deployment.Id != nil {
				fmt.Printf("ID:                 %d\n", *deployment.Id)
			}
			if deployment.ApplicationId != nil {
				fmt.Printf("Application ID:     %s\n", *deployment.ApplicationId)
			}
			if deployment.ApplicationName != nil {
				fmt.Printf("Application Name:   %s\n", *deployment.ApplicationName)
			}
			if deployment.Status != nil {
				fmt.Printf("Status:             %s\n", *deployment.Status)
			}
			if deployment.CreatedAt != nil {
				fmt.Printf("Created At:         %s\n", *deployment.CreatedAt)
			}
			if deployment.UpdatedAt != nil {
				fmt.Printf("Updated At:         %s\n", *deployment.UpdatedAt)
			}
			if deployment.Commit != nil {
				fmt.Printf("Commit:             %s\n", *deployment.Commit)
			}
			if deployment.CommitMessage != nil {
				fmt.Printf("Commit Message:     %s\n", *deployment.CommitMessage)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	return cmd
}
