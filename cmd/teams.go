package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// teamsCmd represents the teams command
var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage teams",
	Long:  "Manage teams in your Coolify instance. List teams, view team details, and manage team members.",
}

// teamsListCmd represents the teams list command
var teamsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all teams",
	Long:  "List all teams that you have access to.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		teams, err := client.Teams().List(context.Background())
		if err != nil {
			return fmt.Errorf("failed to list teams: %w", err)
		}

		outputJSON, _ := cmd.Flags().GetBool("json")
		if outputJSON {
			data, err := json.MarshalIndent(teams, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal teams: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Table output
		fmt.Printf("%-8s %-30s %-50s\n", "ID", "NAME", "DESCRIPTION")
		fmt.Println("------------------------------------------------------------------------------------------------------------")
		for _, team := range teams {
			id := ""
			if team.Id != nil {
				id = fmt.Sprintf("%d", *team.Id)
			}
			name := ""
			if team.Name != nil {
				name = *team.Name
			}
			description := ""
			if team.Description != nil {
				description = *team.Description
			}
			fmt.Printf("%-8s %-30s %-50s\n", id, name, description)
		}

		return nil
	},
}

// teamsGetCmd represents the teams get command
var teamsGetCmd = &cobra.Command{
	Use:   "get <team-id>",
	Short: "Get team details",
	Long:  "Get detailed information about a specific team by ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid team ID: %w", err)
		}

		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		team, err := client.Teams().Get(context.Background(), teamID)
		if err != nil {
			return fmt.Errorf("failed to get team: %w", err)
		}

		outputJSON, _ := cmd.Flags().GetBool("json")
		if outputJSON {
			data, err := json.MarshalIndent(team, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal team: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Formatted output
		fmt.Printf("Team Details:\n")
		fmt.Printf("  ID: %d\n", *team.Id)
		fmt.Printf("  Name: %s\n", *team.Name)
		if team.Description != nil {
			fmt.Printf("  Description: %s\n", *team.Description)
		}

		return nil
	},
}

// teamsGetMembersCmd represents the teams get-members command
var teamsGetMembersCmd = &cobra.Command{
	Use:   "get-members <team-id>",
	Short: "Get team members",
	Long:  "Get all members of a specific team by ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid team ID: %w", err)
		}

		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		members, err := client.Teams().GetMembers(context.Background(), teamID)
		if err != nil {
			return fmt.Errorf("failed to get team members: %w", err)
		}

		outputJSON, _ := cmd.Flags().GetBool("json")
		if outputJSON {
			data, err := json.MarshalIndent(members, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal members: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Table output
		fmt.Printf("%-8s %-30s %-40s\n", "ID", "NAME", "EMAIL")
		fmt.Println("-------------------------------------------------------------------------------")
		for _, member := range members {
			id := ""
			if member.Id != nil {
				id = fmt.Sprintf("%d", *member.Id)
			}
			name := ""
			if member.Name != nil {
				name = *member.Name
			}
			email := ""
			if member.Email != nil {
				email = *member.Email
			}
			fmt.Printf("%-8s %-30s %-40s\n", id, name, email)
		}

		return nil
	},
}

// teamsGetCurrentCmd represents the teams get-current command
var teamsGetCurrentCmd = &cobra.Command{
	Use:   "get-current",
	Short: "Get current team",
	Long:  "Get details of your current team.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		team, err := client.Teams().GetCurrent(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get current team: %w", err)
		}

		outputJSON, _ := cmd.Flags().GetBool("json")
		if outputJSON {
			data, err := json.MarshalIndent(team, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal team: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Formatted output
		fmt.Printf("Current Team:\n")
		fmt.Printf("  ID: %d\n", *team.Id)
		fmt.Printf("  Name: %s\n", *team.Name)
		if team.Description != nil {
			fmt.Printf("  Description: %s\n", *team.Description)
		}

		return nil
	},
}

// teamsGetCurrentMembersCmd represents the teams get-current-members command
var teamsGetCurrentMembersCmd = &cobra.Command{
	Use:   "get-current-members",
	Short: "Get current team members",
	Long:  "Get all members of your current team.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		members, err := client.Teams().GetCurrentMembers(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get current team members: %w", err)
		}

		outputJSON, _ := cmd.Flags().GetBool("json")
		if outputJSON {
			data, err := json.MarshalIndent(members, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal members: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Table output
		fmt.Printf("%-8s %-30s %-40s\n", "ID", "NAME", "EMAIL")
		fmt.Println("-------------------------------------------------------------------------------")
		for _, member := range members {
			id := ""
			if member.Id != nil {
				id = fmt.Sprintf("%d", *member.Id)
			}
			name := ""
			if member.Name != nil {
				name = *member.Name
			}
			email := ""
			if member.Email != nil {
				email = *member.Email
			}
			fmt.Printf("%-8s %-30s %-40s\n", id, name, email)
		}

		return nil
	},
}

func init() {
	// Add subcommands
	teamsCmd.AddCommand(teamsListCmd)
	teamsCmd.AddCommand(teamsGetCmd)
	teamsCmd.AddCommand(teamsGetMembersCmd)
	teamsCmd.AddCommand(teamsGetCurrentCmd)
	teamsCmd.AddCommand(teamsGetCurrentMembersCmd)

	// Add flags
	teamsListCmd.Flags().Bool("json", false, "Output in JSON format")
	teamsGetCmd.Flags().Bool("json", false, "Output in JSON format")
	teamsGetMembersCmd.Flags().Bool("json", false, "Output in JSON format")
	teamsGetCurrentCmd.Flags().Bool("json", false, "Output in JSON format")
	teamsGetCurrentMembersCmd.Flags().Bool("json", false, "Output in JSON format")
}
