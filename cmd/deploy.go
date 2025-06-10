package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	clientpkg "github.com/hongkongkiwi/coolifyme/pkg/client"
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
	cmd.AddCommand(deployWatchCmd())
	cmd.AddCommand(deployLogsCmd())
	cmd.AddCommand(deployMultipleCmd())

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
		RunE: func(_ *cobra.Command, args []string) error {
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
			options := &clientpkg.DeployApplicationOptions{
				Force:  force,
				Branch: branch,
			}
			if pr > 0 {
				options.PR = &pr
			}

			deployResponse, err := client.Deployments().DeployApplicationWithOptions(ctx, applicationUUID, options)
			if err != nil {
				return fmt.Errorf("failed to deploy application: %w", err)
			}

			if deployResponse != nil && len(deployResponse.Deployments) > 0 {
				fmt.Printf("✅ Application deployment triggered successfully for %s\n", applicationUUID)
				for _, deployment := range deployResponse.Deployments {
					fmt.Printf("   📦 Deployment UUID: %s\n", deployment.DeploymentUUID)
					fmt.Printf("   🎯 Resource UUID:   %s\n", deployment.ResourceUUID)
					if deployment.Message != "" {
						fmt.Printf("   📝 Message:         %s\n", deployment.Message)
					}
				}
			} else {
				fmt.Printf("✅ Application deployment triggered successfully for %s\n", applicationUUID)
			}

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
		RunE: func(_ *cobra.Command, args []string) error {
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

			skip, _ := cmd.Flags().GetInt("skip")
			take, _ := cmd.Flags().GetInt("take")

			var deployments []coolify.Application
			if skip > 0 || take > 0 {
				deployments, err = client.Deployments().ListWithPagination(ctx, appUUID, skip, take)
			} else {
				deployments, err = client.Deployments().List(ctx, appUUID)
			}
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
			defer func() {
				_ = w.Flush()
			}()

			// Print header
			_, _ = fmt.Fprintln(w, "UUID\tNAME\tSTATUS\tBRANCH\tDOMAINS")
			_, _ = fmt.Fprintln(w, "----\t----\t------\t------\t-------")

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

				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					uuid, name, status, branch, domains)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	cmd.Flags().Int("skip", 0, "Number of records to skip (pagination)")
	cmd.Flags().Int("take", 10, "Number of records to take (pagination)")

	return cmd
}

func deployListAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-all",
		Aliases: []string{"all"},
		Short:   "List all running deployments",
		Long:    "List all currently running deployments across all applications",
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
				fmt.Println("No running deployments found")
				return nil
			}

			// Create a tabwriter for nicely formatted output
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer func() {
				_ = w.Flush()
			}()

			// Print header
			_, _ = fmt.Fprintln(w, "ID\tAPP NAME\tSTATUS\tCREATED\tSERVER")
			_, _ = fmt.Fprintln(w, "--\t--------\t------\t-------\t------")

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

				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					id, appName, status, created, server)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	cmd.Flags().BoolP("logs", "l", false, "Show deployment logs")

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

			// Display deployment details in a readable format using correct ApplicationDeploymentQueue fields
			fmt.Printf("Deployment Details:\n")
			fmt.Printf("==================\n")
			if deployment.Id != nil {
				fmt.Printf("ID:                 %d\n", *deployment.Id)
			}
			if deployment.DeploymentUuid != nil {
				fmt.Printf("Deployment UUID:    %s\n", *deployment.DeploymentUuid)
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
			if deployment.ServerName != nil {
				fmt.Printf("Server:             %s\n", *deployment.ServerName)
			}
			if deployment.DeploymentUrl != nil && *deployment.DeploymentUrl != "" {
				fmt.Printf("Deployment URL:     %s\n", *deployment.DeploymentUrl)
			}
			if deployment.ForceRebuild != nil {
				fmt.Printf("Force Rebuild:      %t\n", *deployment.ForceRebuild)
			}
			if deployment.IsWebhook != nil {
				fmt.Printf("Triggered by Webhook: %t\n", *deployment.IsWebhook)
			}
			if deployment.IsApi != nil {
				fmt.Printf("Triggered by API:   %t\n", *deployment.IsApi)
			}
			if deployment.PullRequestId != nil && *deployment.PullRequestId > 0 {
				fmt.Printf("Pull Request ID:    %d\n", *deployment.PullRequestId)
			}

			// Show logs if available and not empty
			showLogs, _ := cmd.Flags().GetBool("logs")
			if showLogs && deployment.Logs != nil && *deployment.Logs != "" {
				fmt.Printf("\nDeployment Logs:\n")
				fmt.Printf("===============\n")
				fmt.Printf("%s\n", *deployment.Logs)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	return cmd
}

func deployWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [deployment-uuid]",
		Short: "Watch deployment logs",
		Long:  "Watch the logs for a specific deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			deploymentUUID := args[0]
			ctx := context.Background()

			fmt.Printf("Watching deployment logs for %s\n", deploymentUUID)

			err = client.Deployments().Watch(ctx, deploymentUUID)
			if err != nil {
				return fmt.Errorf("failed to watch deployment logs: %w", err)
			}

			return nil
		},
	}

	return cmd
}

func deployLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [deployment-uuid]",
		Short: "Get deployment logs",
		Long:  "Get the logs for a specific deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			deploymentUUID := args[0]
			ctx := context.Background()

			deployment, err := client.Deployments().GetByUUID(ctx, deploymentUUID)
			if err != nil {
				return fmt.Errorf("failed to get deployment: %w", err)
			}

			logs := ""
			if deployment.Logs != nil {
				logs = *deployment.Logs
			}

			jsonOutput, _ := cmd.Flags().GetBool("json")
			if jsonOutput {
				output, err := json.MarshalIndent(logs, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(output))
				return nil
			}

			fmt.Println(logs)
			return nil
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	return cmd
}

func deployMultipleCmd() *cobra.Command {
	var force bool
	var branch string

	cmd := &cobra.Command{
		Use:   "multiple [uuid1] [uuid2]...",
		Short: "Deploy multiple applications or services",
		Long:  "Trigger deployments for multiple applications or services",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			ctx := context.Background()

			fmt.Printf("🚀 Starting deployments for %d applications/services\n", len(args))
			if branch != "" {
				fmt.Printf("   Branch: %s\n", branch)
			}
			if force {
				fmt.Printf("   Force deployment: enabled\n")
			}

			// Use the multiple deployment method which supports comma-separated UUIDs
			options := &clientpkg.DeployApplicationOptions{
				Force:  force,
				Branch: branch,
			}

			deployResponse, err := client.Deployments().DeployMultiple(ctx, args, options)
			if err != nil {
				return fmt.Errorf("failed to deploy multiple applications: %w", err)
			}

			if deployResponse != nil && len(deployResponse.Deployments) > 0 {
				fmt.Printf("✅ Deployments triggered successfully for %d applications/services\n", len(args))
				for i, deployment := range deployResponse.Deployments {
					fmt.Printf("   %d. 📦 Deployment UUID: %s\n", i+1, deployment.DeploymentUUID)
					fmt.Printf("      🎯 Resource UUID:   %s\n", deployment.ResourceUUID)
					if deployment.Message != "" {
						fmt.Printf("      �� Message:         %s\n", deployment.Message)
					}
				}
			} else {
				fmt.Printf("✅ Deployments triggered successfully for %d applications/services\n", len(args))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deployment even if one is already running")
	cmd.Flags().StringVarP(&branch, "branch", "b", "", "Deploy from specific branch/tag")

	return cmd
}
