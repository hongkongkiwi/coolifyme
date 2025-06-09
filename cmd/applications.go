// Package main provides the coolifyme CLI application for managing Coolify deployments.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/spf13/cobra"
)

// applicationsCmd represents the applications command
var applicationsCmd = &cobra.Command{
	Use:     "applications",
	Aliases: []string{"apps", "app"},
	Short:   "Manage applications",
	Long:    "Manage Coolify applications - list, create, update, and delete applications",
}

// applicationsListCmd represents the applications list command
var applicationsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List applications",
	Long:    "List all applications in your Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		applications, err := client.Applications().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list applications: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(applications, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(applications) == 0 {
			fmt.Println("No applications found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tSTATUS\tGIT REPOSITORY\tDOMAINS")
		_, _ = fmt.Fprintln(w, "----\t----\t------\t--------------\t-------")

		// Print applications
		for _, app := range applications {
			uuid := ""
			name := ""
			status := ""
			gitRepo := ""
			domains := ""

			if app.Uuid != nil {
				uuid = *app.Uuid
			}
			if app.Name != nil {
				name = *app.Name
			}
			if app.Status != nil {
				status = *app.Status
			}
			if app.GitRepository != nil {
				gitRepo = *app.GitRepository
			}
			if app.Fqdn != nil {
				domains = *app.Fqdn
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				uuid, name, status, gitRepo, domains)
		}

		return nil
	},
}

// applicationsGetCmd represents the applications get command
var applicationsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get application details",
	Long:  "Get detailed information about a specific application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		applicationUUID := args[0]

		// Get application details by fetching all applications and filtering
		// This is a workaround since there's no direct get-by-uuid endpoint for applications
		applications, err := client.Applications().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list applications: %w", err)
		}

		var foundApp *coolify.Application

		for i := range applications {
			if applications[i].Uuid != nil && *applications[i].Uuid == applicationUUID {
				foundApp = &applications[i]
				break
			}
		}

		if foundApp == nil {
			return fmt.Errorf("application with UUID %s not found", applicationUUID)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(foundApp, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display application details in a readable format
		fmt.Printf("Application Details:\n")
		fmt.Printf("==================\n")
		if foundApp.Uuid != nil {
			fmt.Printf("UUID:           %s\n", *foundApp.Uuid)
		}
		if foundApp.Name != nil {
			fmt.Printf("Name:           %s\n", *foundApp.Name)
		}
		if foundApp.Status != nil {
			fmt.Printf("Status:         %s\n", *foundApp.Status)
		}
		if foundApp.GitRepository != nil {
			fmt.Printf("Repository:     %s\n", *foundApp.GitRepository)
		}
		if foundApp.GitBranch != nil {
			fmt.Printf("Branch:         %s\n", *foundApp.GitBranch)
		}
		if foundApp.BuildPack != nil {
			fmt.Printf("Build Pack:     %s\n", *foundApp.BuildPack)
		}
		if foundApp.Fqdn != nil {
			fmt.Printf("Domains:        %s\n", *foundApp.Fqdn)
		}

		return nil
	},
}

// applicationsCreateCmd represents the applications create command
var applicationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new application",
	Long:  "Create a new application from a Git repository",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Get flag values
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		buildPack, _ := cmd.Flags().GetString("build-pack")
		project, _ := cmd.Flags().GetString("project")
		server, _ := cmd.Flags().GetString("server")
		environment, _ := cmd.Flags().GetString("environment")

		// Validate required fields
		if repo == "" {
			return fmt.Errorf("repository URL is required (--repo)")
		}
		if branch == "" {
			branch = "main" // default branch
		}
		if buildPack == "" {
			buildPack = "nixpacks" // default build pack
		}
		if project == "" {
			return fmt.Errorf("project UUID is required (--project)")
		}
		if server == "" {
			return fmt.Errorf("server UUID is required (--server)")
		}
		if environment == "" {
			return fmt.Errorf("environment name is required (--environment)")
		}

		fmt.Printf("Creating application...\n")
		fmt.Printf("Repository:   %s\n", repo)
		fmt.Printf("Branch:       %s\n", branch)
		fmt.Printf("Build Pack:   %s\n", buildPack)
		fmt.Printf("Project:      %s\n", project)
		fmt.Printf("Server:       %s\n", server)
		fmt.Printf("Environment:  %s\n", environment)

		// This is a placeholder - the actual implementation would depend on
		// the complete API client implementation
		return fmt.Errorf("application creation is not fully implemented yet - API client needs to be extended")
	},
}

// applicationsDeleteCmd represents the applications delete command
var applicationsDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete an application",
	Long:  "Delete an application by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		deleteVolumes, _ := cmd.Flags().GetBool("delete-volumes")
		deleteConfigs, _ := cmd.Flags().GetBool("delete-configurations")

		options := &coolify.DeleteApplicationByUuidParams{
			DeleteVolumes:        &deleteVolumes,
			DeleteConfigurations: &deleteConfigs,
		}

		err = client.Applications().Delete(context.Background(), args[0], options)
		if err != nil {
			return fmt.Errorf("failed to delete application: %w", err)
		}

		fmt.Printf("Application %s deleted successfully\n", args[0])
		return nil
	},
}

// applicationsUpdateCmd represents the applications update command
var applicationsUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update an application",
	Long:  "Update an application by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// For now, this is a stub - you'd need to implement fields to update
		req := coolify.UpdateApplicationByUuidJSONRequestBody{}

		uuid, err := client.Applications().Update(context.Background(), args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update application: %w", err)
		}

		fmt.Printf("Application updated: %s\n", uuid)
		return nil
	},
}

// applicationsStartCmd represents the applications start command
var applicationsStartCmd = &cobra.Command{
	Use:   "start <uuid>",
	Short: "Start an application",
	Long:  "Start an application by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		force, _ := cmd.Flags().GetBool("force")
		options := &coolify.StartApplicationByUuidParams{
			Force: &force,
		}

		err = client.Applications().Start(context.Background(), args[0], options)
		if err != nil {
			return fmt.Errorf("failed to start application: %w", err)
		}

		fmt.Printf("Application %s started successfully\n", args[0])
		return nil
	},
}

// applicationsStopCmd represents the applications stop command
var applicationsStopCmd = &cobra.Command{
	Use:   "stop <uuid>",
	Short: "Stop an application",
	Long:  "Stop an application by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		err = client.Applications().Stop(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("failed to stop application: %w", err)
		}

		fmt.Printf("Application %s stopped successfully\n", args[0])
		return nil
	},
}

// applicationsRestartCmd represents the applications restart command
var applicationsRestartCmd = &cobra.Command{
	Use:   "restart <uuid>",
	Short: "Restart an application",
	Long:  "Restart an application by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		err = client.Applications().Restart(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("failed to restart application: %w", err)
		}

		fmt.Printf("Application %s restarted successfully\n", args[0])
		return nil
	},
}

// applicationsLogsCmd represents the applications logs command
var applicationsLogsCmd = &cobra.Command{
	Use:   "logs <uuid>",
	Short: "Get application logs",
	Long:  "Get logs for an application by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		lines, _ := cmd.Flags().GetInt("lines")

		params := &coolify.GetApplicationLogsByUuidParams{}
		if lines > 0 {
			lines32 := int32(lines)
			params.Lines = &lines32
		}

		logs, err := client.Applications().GetLogs(context.Background(), args[0], params)
		if err != nil {
			return fmt.Errorf("failed to get application logs: %w", err)
		}

		fmt.Print(logs)
		return nil
	},
}

// applicationsEnvCmd represents the applications env command
var applicationsEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage application environment variables",
	Long:  "Manage environment variables for applications",
}

func init() {
	// Add subcommands to applications
	applicationsCmd.AddCommand(applicationsListCmd)
	applicationsCmd.AddCommand(applicationsGetCmd)
	applicationsCmd.AddCommand(applicationsCreateCmd)
	applicationsCmd.AddCommand(applicationsDeleteCmd)
	applicationsCmd.AddCommand(applicationsUpdateCmd)
	applicationsCmd.AddCommand(applicationsStartCmd)
	applicationsCmd.AddCommand(applicationsStopCmd)
	applicationsCmd.AddCommand(applicationsRestartCmd)
	applicationsCmd.AddCommand(applicationsLogsCmd)
	applicationsCmd.AddCommand(applicationsEnvCmd)

	// Flags for applications list command
	applicationsListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for applications get command
	applicationsGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for applications create command
	applicationsCreateCmd.Flags().String("repo", "", "Git repository URL (required)")
	applicationsCreateCmd.Flags().String("branch", "main", "Git branch")
	applicationsCreateCmd.Flags().String("build-pack", "nixpacks", "Build pack (nixpacks, static, dockerfile, dockercompose)")
	applicationsCreateCmd.Flags().String("project", "", "Project UUID (required)")
	applicationsCreateCmd.Flags().String("server", "", "Server UUID (required)")
	applicationsCreateCmd.Flags().String("environment", "", "Environment name (required)")

	// Delete command flags
	applicationsDeleteCmd.Flags().Bool("force", false, "Force delete")
	applicationsDeleteCmd.Flags().Bool("delete-volumes", false, "Delete volumes")
	applicationsDeleteCmd.Flags().Bool("delete-configurations", false, "Delete configurations")

	// Start command flags
	applicationsStartCmd.Flags().Bool("force", false, "Force start")

	// Logs command flags
	applicationsLogsCmd.Flags().Int("lines", 0, "Number of lines to retrieve")
	applicationsLogsCmd.Flags().Int("since", 0, "Show logs since N seconds ago")

	// Add env subcommands
	applicationsEnvCmd.AddCommand(applicationsEnvListCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvCreateCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvUpdateCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvDeleteCmd)
}

// applicationsEnvListCmd represents the applications env list command
var applicationsEnvListCmd = &cobra.Command{
	Use:   "list <app-uuid>",
	Short: "List environment variables",
	Long:  "List environment variables for an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		envs, err := client.Applications().ListEnvs(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("failed to list environment variables: %w", err)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "json" {
			jsonOutput, err := json.MarshalIndent(envs, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonOutput))
			return nil
		}

		if len(envs) == 0 {
			fmt.Println("No environment variables found")
			return nil
		}

		fmt.Printf("%-36s %-20s %-50s\n", "UUID", "KEY", "VALUE")
		fmt.Println(strings.Repeat("-", 106))
		for _, env := range envs {
			uuid := ""
			key := ""
			value := ""
			if env.Uuid != nil {
				uuid = *env.Uuid
			}
			if env.Key != nil {
				key = *env.Key
			}
			if env.Value != nil {
				value = *env.Value
			}
			fmt.Printf("%-36s %-20s %-50s\n", uuid, key, value)
		}
		return nil
	},
}

// applicationsEnvCreateCmd represents the applications env create command
var applicationsEnvCreateCmd = &cobra.Command{
	Use:   "create <app-uuid> <key> <value>",
	Short: "Create environment variable",
	Long:  "Create a new environment variable for an application",
	Args:  cobra.ExactArgs(3),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		key := args[1]
		value := args[2]
		req := coolify.CreateEnvByApplicationUuidJSONRequestBody{
			Key:   &key,
			Value: &value,
		}

		uuid, err := client.Applications().CreateEnv(context.Background(), args[0], req)
		if err != nil {
			return fmt.Errorf("failed to create environment variable: %w", err)
		}

		fmt.Printf("Environment variable created: %s\n", uuid)
		return nil
	},
}

// applicationsEnvUpdateCmd represents the applications env update command
var applicationsEnvUpdateCmd = &cobra.Command{
	Use:   "update <app-uuid> <key> <value>",
	Short: "Update environment variable",
	Long:  "Update an environment variable for an application",
	Args:  cobra.ExactArgs(3),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := coolify.UpdateEnvByApplicationUuidJSONRequestBody{
			Key:   args[1],
			Value: args[2],
		}

		message, err := client.Applications().UpdateEnv(context.Background(), args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update environment variable: %w", err)
		}

		fmt.Printf("Environment variable updated: %s\n", message)
		return nil
	},
}

// applicationsEnvDeleteCmd represents the applications env delete command
var applicationsEnvDeleteCmd = &cobra.Command{
	Use:   "delete <app-uuid> <env-uuid>",
	Short: "Delete environment variable",
	Long:  "Delete an environment variable for an application",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		message, err := client.Applications().DeleteEnv(context.Background(), args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to delete environment variable: %w", err)
		}

		fmt.Printf("Environment variable deleted: %s\n", message)
		return nil
	},
}
