package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:     "servers",
	Aliases: []string{"server", "srv"},
	Short:   "Manage servers",
	Long:    "Manage Coolify servers - list, create, update, and delete servers",
}

// serversListCmd represents the servers list command
var serversListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List servers",
	Long:    "List all servers in your Coolify instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		servers, err := client.Servers().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list servers: %w", err)
		}

		if len(servers) == 0 {
			fmt.Println("No servers found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer w.Flush()

		// Print header
		fmt.Fprintln(w, "UUID\tNAME\tIP\tSTATUS\tDESCRIPTION")
		fmt.Fprintln(w, "----\t----\t--\t------\t-----------")

		// Print servers
		for _, server := range servers {
			uuid := ""
			name := ""
			ip := ""
			status := ""
			description := ""

			if server.Uuid != nil {
				uuid = *server.Uuid
			}
			if server.Name != nil {
				name = *server.Name
			}
			if server.Ip != nil {
				ip = *server.Ip
			}
			if server.ValidationLogs != nil {
				status = "validated"
			} else {
				status = "unknown"
			}
			if server.Description != nil {
				description = *server.Description
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", uuid, name, ip, status, description)
		}

		return nil
	},
}

func init() {
	// Add subcommands to servers
	serversCmd.AddCommand(serversListCmd)

	// Flags for servers list command
	serversListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
