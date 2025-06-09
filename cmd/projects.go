package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

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
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		projects, err := client.Projects().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
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

func init() {
	// Add subcommands to projects
	projectsCmd.AddCommand(projectsListCmd)

	// Flags for projects list command
	projectsListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
