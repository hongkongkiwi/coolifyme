package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project", "proj"},
	Short:   "Manage projects",
	Long:    "Manage Coolify projects - list, create, update, and delete projects",
}

// projectsListCmd represents the projects list command
var projectsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects",
	Long:    "List all projects in your Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		projects, err := client.Projects().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(projects, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(projects) == 0 {
			fmt.Println("No projects found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tDESCRIPTION")
		_, _ = fmt.Fprintln(w, "----\t----\t-----------")

		// Print projects
		for _, project := range projects {
			uuid := ""
			name := ""
			description := ""

			if project.Uuid != nil {
				uuid = *project.Uuid
			}
			if project.Name != nil {
				name = *project.Name
			}
			if project.Description != nil {
				description = *project.Description
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", uuid, name, description)
		}

		return nil
	},
}

// projectsGetCmd represents the projects get command
var projectsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get project details",
	Long:  "Get detailed information about a specific project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		projectUUID := args[0]

		project, err := client.Projects().Get(ctx, projectUUID)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(project, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display project details in a readable format
		fmt.Printf("Project Details:\n")
		fmt.Printf("===============\n")
		if project.Uuid != nil {
			fmt.Printf("UUID:         %s\n", *project.Uuid)
		}
		if project.Name != nil {
			fmt.Printf("Name:         %s\n", *project.Name)
		}
		if project.Description != nil {
			fmt.Printf("Description:  %s\n", *project.Description)
		}

		return nil
	},
}

// projectsCreateCmd represents the projects create command
var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create project",
	Long:  "Create a new project",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("project name is required")
		}

		req := coolify.CreateProjectJSONRequestBody{
			Name:        &name,
			Description: &description,
		}

		result, err := client.Projects().Create(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		fmt.Printf("✅ Project created successfully\n")
		fmt.Printf("   UUID: %s\n", result)

		return nil
	},
}

// projectsUpdateCmd represents the projects update command
var projectsUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update project",
	Long:  "Update an existing project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		projectUUID := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		req := coolify.UpdateProjectByUuidJSONRequestBody{}
		if name != "" {
			req.Name = &name
		}
		if description != "" {
			req.Description = &description
		}

		result, err := client.Projects().Update(context.Background(), projectUUID, req)
		if err != nil {
			return fmt.Errorf("failed to update project: %w", err)
		}

		fmt.Printf("✅ Project updated successfully\n")
		if result.Uuid != nil {
			fmt.Printf("   UUID: %s\n", *result.Uuid)
		}

		return nil
	},
}

// projectsDeleteCmd represents the projects delete command
var projectsDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete project",
	Long:  "Delete a project by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		projectUUID := args[0]

		err = client.Projects().Delete(ctx, projectUUID)
		if err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}

		fmt.Printf("✅ Project %s deleted successfully\n", projectUUID)
		return nil
	},
}

// projectsGetEnvironmentCmd represents the projects get-environment command
var projectsGetEnvironmentCmd = &cobra.Command{
	Use:   "get-environment <project-uuid> <environment-name-or-uuid>",
	Short: "Get environment from project",
	Long:  "Get environment details from a specific project",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		projectUUID := args[0]
		environmentNameOrUUID := args[1]

		environment, err := client.Projects().GetEnvironment(ctx, projectUUID, environmentNameOrUUID)
		if err != nil {
			return fmt.Errorf("failed to get environment: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(environment, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display environment details in a readable format
		fmt.Printf("Environment Details:\n")
		fmt.Printf("===================\n")
		if environment.Id != nil {
			fmt.Printf("ID:           %d\n", *environment.Id)
		}
		if environment.Name != nil {
			fmt.Printf("Name:         %s\n", *environment.Name)
		}
		if environment.Description != nil {
			fmt.Printf("Description:  %s\n", *environment.Description)
		}
		if environment.ProjectId != nil {
			fmt.Printf("Project ID:   %d\n", *environment.ProjectId)
		}

		return nil
	},
}

func init() {
	// Add subcommands to projects
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsUpdateCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)
	projectsCmd.AddCommand(projectsGetEnvironmentCmd)

	// Flags for list command
	projectsListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for get command
	projectsGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for create command
	projectsCreateCmd.Flags().StringP("name", "n", "", "Name of the project (required)")
	projectsCreateCmd.Flags().StringP("description", "d", "", "Description of the project")
	_ = projectsCreateCmd.MarkFlagRequired("name")

	// Flags for update command
	projectsUpdateCmd.Flags().StringP("name", "n", "", "Name of the project")
	projectsUpdateCmd.Flags().StringP("description", "d", "", "Description of the project")

	// Flags for get-environment command
	projectsGetEnvironmentCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
