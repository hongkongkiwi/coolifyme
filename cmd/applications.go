// Package main provides the coolifyme CLI application for managing Coolify deployments.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

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

		// Get application details directly using the UUID endpoint
		foundApp, err := client.Applications().Get(ctx, applicationUUID)
		if err != nil {
			return fmt.Errorf("failed to get application: %w", err)
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

		startResponse, err := client.Applications().Start(context.Background(), args[0], options)
		if err != nil {
			return fmt.Errorf("failed to start application: %w", err)
		}

		if startResponse != nil {
			fmt.Printf("✅ Application %s started successfully\n", args[0])
			if startResponse.DeploymentUUID != "" {
				fmt.Printf("   📦 Deployment UUID: %s\n", startResponse.DeploymentUUID)
			}
			if startResponse.Message != "" {
				fmt.Printf("   💬 Message: %s\n", startResponse.Message)
			}
		} else {
			fmt.Printf("Application %s started successfully\n", args[0])
		}
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

		restartResponse, err := client.Applications().Restart(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("failed to restart application: %w", err)
		}

		if restartResponse != nil {
			fmt.Printf("✅ Application %s restarted successfully\n", args[0])
			if restartResponse.DeploymentUUID != "" {
				fmt.Printf("   📦 Deployment UUID: %s\n", restartResponse.DeploymentUUID)
			}
			if restartResponse.Message != "" {
				fmt.Printf("   💬 Message: %s\n", restartResponse.Message)
			}
		} else {
			fmt.Printf("Application %s restarted successfully\n", args[0])
		}
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
	applicationsEnvCmd.AddCommand(applicationsEnvUpdateBulkCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvDeleteCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvExportCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvImportCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvSyncCmd)
	applicationsEnvCmd.AddCommand(applicationsEnvCleanupCmd)

	// Flags for bulk environment variable update command
	applicationsEnvUpdateBulkCmd.Flags().StringP("env-data", "d", "", "JSON string containing environment variables")
	applicationsEnvUpdateBulkCmd.Flags().StringP("env-file", "f", "", "File containing environment variables in JSON format")

	// Flags for .env file management commands
	applicationsEnvExportCmd.Flags().StringP("file", "f", ".env", "Output .env file path")
	applicationsEnvExportCmd.Flags().Bool("overwrite", false, "Overwrite existing file")
	applicationsEnvImportCmd.Flags().StringP("file", "f", ".env", "Input .env file path")
	applicationsEnvImportCmd.Flags().Bool("dry-run", false, "Show what would be imported without making changes")
	applicationsEnvSyncCmd.Flags().StringP("file", "f", ".env", ".env file to sync")
	applicationsEnvSyncCmd.Flags().Bool("dry-run", false, "Show what would be changed without making changes")
	applicationsEnvCleanupCmd.Flags().StringP("file", "f", ".env", ".env file to clean up")
	applicationsEnvCleanupCmd.Flags().Bool("dry-run", false, "Show what would be removed without making changes")
	applicationsEnvCleanupCmd.Flags().Bool("backup", true, "Create backup before cleaning up")
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

// applicationsEnvUpdateBulkCmd represents the applications env update-bulk command
var applicationsEnvUpdateBulkCmd = &cobra.Command{
	Use:   "update-bulk <app-uuid>",
	Short: "Bulk update environment variables",
	Long:  "Update multiple environment variables for an application from a file or JSON string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		envDataFlag, _ := cmd.Flags().GetString("env-data")
		envFile, _ := cmd.Flags().GetString("env-file")

		if envDataFlag == "" && envFile == "" {
			return fmt.Errorf("either --env-data or --env-file is required")
		}

		var envVarsList []interface{}
		if envFile != "" {
			// Read environment variables from file
			content, err := safeReadFile(envFile)
			if err != nil {
				return fmt.Errorf("failed to read env file: %w", err)
			}
			if err := json.Unmarshal(content, &envVarsList); err != nil {
				return fmt.Errorf("failed to parse env file JSON: %w", err)
			}
		} else {
			// Parse environment variables from JSON string
			if err := json.Unmarshal([]byte(envDataFlag), &envVarsList); err != nil {
				return fmt.Errorf("failed to parse env data JSON: %w", err)
			}
		}

		// Convert to the expected structure for applications
		var envStructs []struct {
			IsBuildTime *bool   `json:"is_build_time,omitempty"`
			IsLiteral   *bool   `json:"is_literal,omitempty"`
			IsMultiline *bool   `json:"is_multiline,omitempty"`
			IsPreview   *bool   `json:"is_preview,omitempty"`
			IsShownOnce *bool   `json:"is_shown_once,omitempty"`
			Key         *string `json:"key,omitempty"`
			Value       *string `json:"value,omitempty"`
		}

		// Parse each environment variable
		for _, item := range envVarsList {
			itemData, _ := json.Marshal(item)
			var envVar struct {
				IsBuildTime *bool   `json:"is_build_time,omitempty"`
				IsLiteral   *bool   `json:"is_literal,omitempty"`
				IsMultiline *bool   `json:"is_multiline,omitempty"`
				IsPreview   *bool   `json:"is_preview,omitempty"`
				IsShownOnce *bool   `json:"is_shown_once,omitempty"`
				Key         *string `json:"key,omitempty"`
				Value       *string `json:"value,omitempty"`
			}
			if err := json.Unmarshal(itemData, &envVar); err == nil {
				envStructs = append(envStructs, envVar)
			}
		}

		// Create request body
		req := coolify.UpdateEnvsByApplicationUuidJSONRequestBody{
			Data: envStructs,
		}

		ctx := context.Background()
		appUUID := args[0]

		message, err := client.Applications().UpdateEnvs(ctx, appUUID, req)
		if err != nil {
			return fmt.Errorf("failed to bulk update environment variables: %w", err)
		}

		fmt.Printf("✅ Environment variables updated successfully\n")
		fmt.Printf("   💬 Message: %s\n", message)
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

// applicationsEnvExportCmd represents the applications env export command
var applicationsEnvExportCmd = &cobra.Command{
	Use:   "export <app-uuid>",
	Short: "Export environment variables to .env file",
	Long:  "Export all environment variables from an application to a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		appUUID := args[0]
		filename, _ := cmd.Flags().GetString("file")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		// Check if file exists and overwrite flag
		if _, err := os.Stat(filename); err == nil && !overwrite {
			return fmt.Errorf("file %s already exists, use --overwrite to replace it", filename)
		}

		// Get environment variables
		envs, err := client.Applications().ListEnvs(context.Background(), appUUID)
		if err != nil {
			return fmt.Errorf("failed to list environment variables: %w", err)
		}

		// Create .env content
		var envContent strings.Builder
		envContent.WriteString("# Environment variables exported from Coolify\n")
		envContent.WriteString(fmt.Sprintf("# Application UUID: %s\n", appUUID))
		envContent.WriteString(fmt.Sprintf("# Exported at: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

		for _, env := range envs {
			if env.Key != nil && env.Value != nil {
				key := *env.Key
				value := *env.Value

				// Handle multiline values by quoting them
				if strings.Contains(value, "\n") {
					value = fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\\\""))
				}

				envContent.WriteString(fmt.Sprintf("%s=%s\n", key, value))
			}
		}

		// Write to file
		if err := os.WriteFile(filename, []byte(envContent.String()), 0o600); err != nil {
			return fmt.Errorf("failed to write .env file: %w", err)
		}

		fmt.Printf("✅ Environment variables exported to %s\n", filename)
		fmt.Printf("   📝 Exported %d variables\n", len(envs))
		return nil
	},
}

// applicationsEnvImportCmd represents the applications env import command
var applicationsEnvImportCmd = &cobra.Command{
	Use:   "import <app-uuid>",
	Short: "Import environment variables from .env file",
	Long:  "Import environment variables from a .env file to an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		appUUID := args[0]
		filename, _ := cmd.Flags().GetString("file")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Read .env file
		content, err := safeReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read .env file: %w", err)
		}

		// Parse .env file
		envVars := parseEnvFile(string(content))
		if len(envVars) == 0 {
			fmt.Println("No environment variables found in .env file")
			return nil
		}

		if dryRun {
			fmt.Printf("🔍 Dry run: Would import %d environment variables:\n", len(envVars))
			for key, value := range envVars {
				fmt.Printf("   %s=%s\n", key, value)
			}
			return nil
		}

		// Convert to bulk update format
		var envStructs []struct {
			IsBuildTime *bool   `json:"is_build_time,omitempty"`
			IsLiteral   *bool   `json:"is_literal,omitempty"`
			IsMultiline *bool   `json:"is_multiline,omitempty"`
			IsPreview   *bool   `json:"is_preview,omitempty"`
			IsShownOnce *bool   `json:"is_shown_once,omitempty"`
			Key         *string `json:"key,omitempty"`
			Value       *string `json:"value,omitempty"`
		}

		for key, value := range envVars {
			k := key
			v := value
			envStructs = append(envStructs, struct {
				IsBuildTime *bool   `json:"is_build_time,omitempty"`
				IsLiteral   *bool   `json:"is_literal,omitempty"`
				IsMultiline *bool   `json:"is_multiline,omitempty"`
				IsPreview   *bool   `json:"is_preview,omitempty"`
				IsShownOnce *bool   `json:"is_shown_once,omitempty"`
				Key         *string `json:"key,omitempty"`
				Value       *string `json:"value,omitempty"`
			}{
				Key:   &k,
				Value: &v,
			})
		}

		// Create request body
		req := coolify.UpdateEnvsByApplicationUuidJSONRequestBody{
			Data: envStructs,
		}

		message, err := client.Applications().UpdateEnvs(context.Background(), appUUID, req)
		if err != nil {
			return fmt.Errorf("failed to import environment variables: %w", err)
		}

		fmt.Printf("✅ Environment variables imported from %s\n", filename)
		fmt.Printf("   📝 Imported %d variables\n", len(envVars))
		fmt.Printf("   💬 Message: %s\n", message)
		return nil
	},
}

// applicationsEnvSyncCmd represents the applications env sync command
var applicationsEnvSyncCmd = &cobra.Command{
	Use:   "sync <app-uuid>",
	Short: "Sync .env file with application environment variables",
	Long:  "Synchronize a .env file with the application's environment variables (bidirectional sync)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		appUUID := args[0]
		filename, _ := cmd.Flags().GetString("file")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Get current environment variables from application
		appEnvs, err := client.Applications().ListEnvs(context.Background(), appUUID)
		if err != nil {
			return fmt.Errorf("failed to list environment variables: %w", err)
		}

		appEnvMap := make(map[string]string)
		for _, env := range appEnvs {
			if env.Key != nil && env.Value != nil {
				appEnvMap[*env.Key] = *env.Value
			}
		}

		// Read .env file (create if doesn't exist)
		var fileEnvMap map[string]string
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fileEnvMap = make(map[string]string)
			fmt.Printf("📄 .env file %s doesn't exist, will create it\n", filename)
		} else {
			content, err := safeReadFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read .env file: %w", err)
			}
			fileEnvMap = parseEnvFile(string(content))
		}

		// Compare and plan changes
		toAddToApp := make(map[string]string)
		toAddToFile := make(map[string]string)
		toUpdateInApp := make(map[string]string)
		toUpdateInFile := make(map[string]string)

		// Check file vars against app vars
		for key, value := range fileEnvMap {
			if appValue, exists := appEnvMap[key]; exists {
				if appValue != value {
					toUpdateInApp[key] = value
				}
			} else {
				toAddToApp[key] = value
			}
		}

		// Check app vars against file vars
		for key, value := range appEnvMap {
			if fileValue, exists := fileEnvMap[key]; exists {
				if fileValue != value {
					toUpdateInFile[key] = value
				}
			} else {
				toAddToFile[key] = value
			}
		}

		if dryRun {
			fmt.Printf("🔍 Sync analysis for %s:\n", filename)
			fmt.Printf("   📤 Would add to application: %d variables\n", len(toAddToApp))
			fmt.Printf("   📥 Would add to .env file: %d variables\n", len(toAddToFile))
			fmt.Printf("   🔄 Would update in application: %d variables\n", len(toUpdateInApp))
			fmt.Printf("   🔄 Would update in .env file: %d variables\n", len(toUpdateInFile))
			return nil
		}

		// Perform sync operations
		hasChanges := false

		// Update application if needed
		if len(toAddToApp) > 0 || len(toUpdateInApp) > 0 {
			// Merge all app changes
			allAppChanges := make(map[string]string)
			for k, v := range toAddToApp {
				allAppChanges[k] = v
			}
			for k, v := range toUpdateInApp {
				allAppChanges[k] = v
			}

			var envStructs []struct {
				IsBuildTime *bool   `json:"is_build_time,omitempty"`
				IsLiteral   *bool   `json:"is_literal,omitempty"`
				IsMultiline *bool   `json:"is_multiline,omitempty"`
				IsPreview   *bool   `json:"is_preview,omitempty"`
				IsShownOnce *bool   `json:"is_shown_once,omitempty"`
				Key         *string `json:"key,omitempty"`
				Value       *string `json:"value,omitempty"`
			}

			for key, value := range allAppChanges {
				k := key
				v := value
				envStructs = append(envStructs, struct {
					IsBuildTime *bool   `json:"is_build_time,omitempty"`
					IsLiteral   *bool   `json:"is_literal,omitempty"`
					IsMultiline *bool   `json:"is_multiline,omitempty"`
					IsPreview   *bool   `json:"is_preview,omitempty"`
					IsShownOnce *bool   `json:"is_shown_once,omitempty"`
					Key         *string `json:"key,omitempty"`
					Value       *string `json:"value,omitempty"`
				}{
					Key:   &k,
					Value: &v,
				})
			}

			req := coolify.UpdateEnvsByApplicationUuidJSONRequestBody{
				Data: envStructs,
			}

			_, err := client.Applications().UpdateEnvs(context.Background(), appUUID, req)
			if err != nil {
				return fmt.Errorf("failed to update application environment variables: %w", err)
			}
			hasChanges = true
		}

		// Update .env file if needed
		if len(toAddToFile) > 0 || len(toUpdateInFile) > 0 {
			// Merge all changes into file map
			for k, v := range toAddToFile {
				fileEnvMap[k] = v
			}
			for k, v := range toUpdateInFile {
				fileEnvMap[k] = v
			}

			// Generate new .env content
			var envContent strings.Builder
			envContent.WriteString("# Environment variables synced with Coolify\n")
			envContent.WriteString(fmt.Sprintf("# Application UUID: %s\n", appUUID))
			envContent.WriteString(fmt.Sprintf("# Last synced: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

			for key, value := range fileEnvMap {
				// Handle multiline values by quoting them
				if strings.Contains(value, "\n") {
					value = fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\\\""))
				}
				envContent.WriteString(fmt.Sprintf("%s=%s\n", key, value))
			}

			if err := os.WriteFile(filename, []byte(envContent.String()), 0o600); err != nil {
				return fmt.Errorf("failed to write .env file: %w", err)
			}
			hasChanges = true
		}

		if hasChanges {
			fmt.Printf("✅ Environment variables synchronized\n")
			fmt.Printf("   📤 Added/updated in application: %d variables\n", len(toAddToApp)+len(toUpdateInApp))
			fmt.Printf("   📥 Added/updated in .env file: %d variables\n", len(toAddToFile)+len(toUpdateInFile))
		} else {
			fmt.Printf("✅ Environment variables are already synchronized\n")
		}

		return nil
	},
}

// applicationsEnvCleanupCmd represents the applications env cleanup command
var applicationsEnvCleanupCmd = &cobra.Command{
	Use:   "cleanup <app-uuid>",
	Short: "Clean up .env file",
	Long:  "Remove environment variables from .env file that don't exist in the application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		appUUID := args[0]
		filename, _ := cmd.Flags().GetString("file")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		backup, _ := cmd.Flags().GetBool("backup")

		// Check if .env file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf(".env file %s does not exist", filename)
		}

		// Get current environment variables from application
		appEnvs, err := client.Applications().ListEnvs(context.Background(), appUUID)
		if err != nil {
			return fmt.Errorf("failed to list environment variables: %w", err)
		}

		appEnvKeys := make(map[string]bool)
		for _, env := range appEnvs {
			if env.Key != nil {
				appEnvKeys[*env.Key] = true
			}
		}

		// Read .env file
		content, err := safeReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read .env file: %w", err)
		}

		fileEnvMap := parseEnvFile(string(content))

		// Find variables to remove
		toRemove := make([]string, 0)
		for key := range fileEnvMap {
			if !appEnvKeys[key] {
				toRemove = append(toRemove, key)
			}
		}

		if len(toRemove) == 0 {
			fmt.Printf("✅ .env file is already clean - no variables to remove\n")
			return nil
		}

		if dryRun {
			fmt.Printf("🔍 Cleanup analysis for %s:\n", filename)
			fmt.Printf("   🗑️  Would remove %d variables not in application:\n", len(toRemove))
			for _, key := range toRemove {
				fmt.Printf("      - %s\n", key)
			}
			return nil
		}

		// Create backup if requested
		if backup {
			backupFilename := filename + ".backup." + time.Now().Format("20060102-150405")
			if err := os.WriteFile(backupFilename, content, 0o600); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
			fmt.Printf("📄 Backup created: %s\n", backupFilename)
		}

		// Remove variables from map
		for _, key := range toRemove {
			delete(fileEnvMap, key)
		}

		// Generate new .env content
		var envContent strings.Builder
		envContent.WriteString("# Environment variables cleaned up\n")
		envContent.WriteString(fmt.Sprintf("# Application UUID: %s\n", appUUID))
		envContent.WriteString(fmt.Sprintf("# Cleaned up: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

		for key, value := range fileEnvMap {
			// Handle multiline values by quoting them
			if strings.Contains(value, "\n") {
				value = fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\\\""))
			}
			envContent.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		}

		if err := os.WriteFile(filename, []byte(envContent.String()), 0o600); err != nil {
			return fmt.Errorf("failed to write cleaned .env file: %w", err)
		}

		fmt.Printf("✅ .env file cleaned up\n")
		fmt.Printf("   🗑️  Removed %d variables\n", len(toRemove))
		fmt.Printf("   📝 Remaining %d variables\n", len(fileEnvMap))

		return nil
	},
}

// parseEnvFile parses a .env file content and returns a map of key-value pairs
func parseEnvFile(content string) map[string]string {
	envMap := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Find the first = sign
		eqIndex := strings.Index(line, "=")
		if eqIndex == -1 {
			continue
		}

		key := strings.TrimSpace(line[:eqIndex])
		value := strings.TrimSpace(line[eqIndex+1:])

		// Remove quotes if present
		if len(value) >= 2 &&
			((strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'"))) {
			value = value[1 : len(value)-1]
			// Unescape quotes
			value = strings.ReplaceAll(value, "\\\"", "\"")
			value = strings.ReplaceAll(value, "\\'", "'")
		}

		if key != "" {
			envMap[key] = value
		}
	}

	return envMap
}
