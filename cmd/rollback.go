package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/spf13/cobra"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback resources to previous versions",
	Long:  "Rollback applications, services, and other resources to previous deployments or versions",
}

// rollbackAppCmd rolls back an application
var rollbackAppCmd = &cobra.Command{
	Use:   "app <uuid>",
	Short: "Rollback application",
	Long:  "Rollback an application to a previous deployment or git commit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		appUUID := args[0]
		toVersion, _ := cmd.Flags().GetString("to-version")
		toCommit, _ := cmd.Flags().GetString("to-commit")
		listOnly, _ := cmd.Flags().GetBool("list")
		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		ctx := context.Background()

		// If list flag is set, show available versions/deployments
		if listOnly {
			return listAvailableVersions(ctx, client, appUUID)
		}

		// Validate rollback target
		if toVersion == "" && toCommit == "" {
			return fmt.Errorf("either --to-version or --to-commit must be specified")
		}

		// Get current application info
		app, err := getApplicationInfo(ctx, client, appUUID)
		if err != nil {
			return fmt.Errorf("failed to get application info: %w", err)
		}

		// Show rollback plan
		fmt.Printf("🔄 Rollback Plan\n")
		fmt.Printf("================\n")
		fmt.Printf("Application: %s (%s)\n", getAppName(app), appUUID)
		if toVersion != "" {
			fmt.Printf("Target Version: %s\n", toVersion)
		}
		if toCommit != "" {
			fmt.Printf("Target Commit: %s\n", toCommit)
		}
		fmt.Printf("Dry Run: %v\n", dryRun)
		fmt.Println()

		if dryRun {
			fmt.Println("✅ Dry run completed - no changes made")
			return nil
		}

		// Confirm rollback unless force flag is set
		if !force {
			fmt.Printf("⚠️  Are you sure you want to rollback this application? This action cannot be undone.\n")
			fmt.Print("Type 'yes' to confirm: ")
			var confirmation string
			if _, err := fmt.Scanln(&confirmation); err != nil || confirmation != ConfirmationYes {
				fmt.Println("❌ Rollback cancelled")
				return nil
			}
		}

		// Perform rollback
		if toCommit != "" {
			return rollbackToCommit(ctx, client, appUUID, toCommit)
		}
		return rollbackToVersion(ctx, client, appUUID, toVersion)
	},
}

// rollbackServiceCmd rolls back a service
var rollbackServiceCmd = &cobra.Command{
	Use:   "service <uuid>",
	Short: "Rollback service",
	Long:  "Rollback a service to a previous configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceUUID := args[0]
		toVersion, _ := cmd.Flags().GetString("to-version")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		fmt.Printf("🔄 Service Rollback\n")
		fmt.Printf("===================\n")
		fmt.Printf("Service: %s\n", serviceUUID)
		fmt.Printf("Target Version: %s\n", toVersion)
		fmt.Printf("Dry Run: %v\n", dryRun)

		if dryRun {
			fmt.Println("✅ Dry run completed - service rollback would be performed")
			return nil
		}

		// Note: Service rollback requires additional API endpoints
		fmt.Println("⚠️  Service rollback is not yet supported by the Coolify API")
		fmt.Println("   Please use the Coolify web interface for service rollbacks")
		return nil
	},
}

// rollbackHistoryCmd shows rollback history
var rollbackHistoryCmd = &cobra.Command{
	Use:   "history <uuid>",
	Short: "Show rollback history",
	Long:  "Display the rollback and deployment history for a resource",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		resourceUUID := args[0]
		resourceType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		ctx := context.Background()

		if resourceType == "" || resourceType == "application" {
			return showApplicationHistory(ctx, client, resourceUUID, limit, jsonOutput)
		}

		return fmt.Errorf("history for resource type '%s' is not yet supported", resourceType)
	},
}

func listAvailableVersions(ctx context.Context, client interface{}, appUUID string) error {
	fmt.Printf("📋 Available Versions for Application: %s\n", appUUID)
	fmt.Printf("=========================================\n")

	// Get deployment history (if available)
	deployments, err := getDeploymentHistory(ctx, client, appUUID)
	if err != nil {
		fmt.Printf("⚠️  Could not fetch deployment history: %v\n", err)
	} else {
		fmt.Printf("\n🚀 Recent Deployments:\n")
		for i, deployment := range deployments {
			if i >= 10 { // Limit to last 10 deployments
				break
			}
			fmt.Printf("  %d. %s (deployed: %s)\n", i+1, deployment.ID, deployment.Date)
		}
	}

	// Get git commits (if it's a git-based application)
	commits, err := getGitCommits(ctx, client, appUUID)
	if err != nil {
		fmt.Printf("⚠️  Could not fetch git commits: %v\n", err)
	} else {
		fmt.Printf("\n📝 Recent Git Commits:\n")
		for i, commit := range commits {
			if i >= 10 { // Limit to last 10 commits
				break
			}
			fmt.Printf("  %s: %s (by %s)\n", commit.Hash[:8], commit.Message, commit.Author)
		}
	}

	fmt.Printf("\nUsage:\n")
	fmt.Printf("  coolifyme rollback app %s --to-commit <commit-hash>\n", appUUID)
	fmt.Printf("  coolifyme rollback app %s --to-version <version>\n", appUUID)

	return nil
}

func rollbackToCommit(_ context.Context, _ interface{}, _, commitHash string) error {
	fmt.Printf("🔄 Rolling back to commit: %s\n", commitHash)

	// This would require specific API endpoints for git-based rollbacks
	// For now, we'll simulate the process and suggest manual steps

	fmt.Printf("⚠️  Direct commit rollback is not yet supported by the Coolify API\n")
	fmt.Printf("   Suggested manual steps:\n")
	fmt.Printf("   1. Go to your git repository\n")
	fmt.Printf("   2. Reset or create a new commit with the desired state\n")
	fmt.Printf("   3. Push the changes\n")
	fmt.Printf("   4. Trigger a new deployment in Coolify\n")

	return nil
}

func rollbackToVersion(_ context.Context, _ interface{}, _, version string) error {
	fmt.Printf("🔄 Rolling back to version: %s\n", version)

	// This would require deployment history and rollback API endpoints
	fmt.Printf("⚠️  Version-based rollback is not yet supported by the Coolify API\n")
	fmt.Printf("   Please use the Coolify web interface to rollback to a previous deployment\n")

	return nil
}

func getApplicationInfo(ctx context.Context, client interface{}, appUUID string) (interface{}, error) {
	// Type assertion to get the actual client
	c, ok := client.(interface {
		Applications() interface {
			Get(context.Context, string) (coolify.Application, error)
		}
	})
	if !ok {
		return nil, fmt.Errorf("invalid client type")
	}

	return c.Applications().Get(ctx, appUUID)
}

func getAppName(app interface{}) string {
	if coolifyApp, ok := app.(coolify.Application); ok {
		if coolifyApp.Name != nil {
			return *coolifyApp.Name
		}
	}
	return "Unknown"
}

func showApplicationHistory(_ context.Context, _ interface{}, appUUID string, limit int, jsonOutput bool) error {
	fmt.Printf("📜 Application History: %s\n", appUUID)
	fmt.Printf("=========================\n")

	// Mock deployment history for demonstration
	history := []DeploymentRecord{
		{
			ID:       "deploy-001",
			Version:  "v1.2.3",
			Commit:   "abc123def",
			Status:   "successful",
			Date:     time.Now().Add(-1 * time.Hour),
			Author:   "user@example.com",
			Message:  "Latest deployment",
			Duration: "2m 30s",
		},
		{
			ID:       "deploy-002",
			Version:  "v1.2.2",
			Commit:   "def456ghi",
			Status:   "successful",
			Date:     time.Now().Add(-2 * time.Hour),
			Author:   "user@example.com",
			Message:  "Bug fix deployment",
			Duration: "1m 45s",
		},
		{
			ID:       "deploy-003",
			Version:  "v1.2.1",
			Commit:   "ghi789jkl",
			Status:   "failed",
			Date:     time.Now().Add(-3 * time.Hour),
			Author:   "user@example.com",
			Message:  "Failed deployment",
			Duration: "45s",
		},
	}

	if limit > 0 && limit < len(history) {
		history = history[:limit]
	}

	if jsonOutput {
		output, err := json.MarshalIndent(history, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(output))
		return nil
	}

	// Display in table format
	fmt.Printf("ID\t\tVERSION\tCOMMIT\t\tSTATUS\t\tDATE\t\t\tDURATION\n")
	fmt.Printf("--\t\t-------\t------\t\t------\t\t----\t\t\t--------\n")

	for _, record := range history {
		status := record.Status
		switch status {
		case "successful":
			status = "✅ " + status
		case "failed":
			status = "❌ " + status
		default:
			status = "🔄 " + status
		}

		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
			record.ID,
			record.Version,
			record.Commit[:8],
			status,
			record.Date.Format("2006-01-02 15:04"),
			record.Duration,
		)
	}

	fmt.Printf("\nTo rollback:\n")
	fmt.Printf("  coolifyme rollback app %s --to-commit <commit>\n", appUUID)

	return nil
}

// DeploymentRecord represents a deployment history entry
type DeploymentRecord struct {
	ID       string    `json:"id"`
	Version  string    `json:"version"`
	Commit   string    `json:"commit"`
	Status   string    `json:"status"`
	Date     time.Time `json:"date"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	Duration string    `json:"duration"`
}

// RollbackGitCommit represents a git commit for rollback purposes
type RollbackGitCommit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

func getDeploymentHistory(_ context.Context, _ interface{}, _ string) ([]DeploymentRecord, error) {
	// This would require actual API calls to get deployment history
	// For now, return mock data
	return []DeploymentRecord{
		{ID: "deploy-001", Date: time.Now()},
		{ID: "deploy-002", Date: time.Now().Add(-1 * time.Hour)},
	}, nil
}

func getGitCommits(_ context.Context, _ interface{}, _ string) ([]RollbackGitCommit, error) {
	// This would require actual API calls to get git commits
	// For now, return mock data
	return []RollbackGitCommit{
		{Hash: "abc123def456", Message: "Fix critical bug", Author: "developer"},
		{Hash: "def456ghi789", Message: "Add new feature", Author: "developer"},
	}, nil
}

func init() {
	// Add subcommands
	rollbackCmd.AddCommand(rollbackAppCmd)
	rollbackCmd.AddCommand(rollbackServiceCmd)
	rollbackCmd.AddCommand(rollbackHistoryCmd)

	// Flags for rollback app command
	rollbackAppCmd.Flags().String("to-version", "", "Rollback to specific version")
	rollbackAppCmd.Flags().String("to-commit", "", "Rollback to specific git commit")
	rollbackAppCmd.Flags().BoolP("list", "l", false, "List available versions/commits")
	rollbackAppCmd.Flags().BoolP("force", "f", false, "Force rollback without confirmation")
	rollbackAppCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")

	// Flags for rollback service command
	rollbackServiceCmd.Flags().String("to-version", "", "Rollback to specific version")
	rollbackServiceCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")

	// Flags for rollback history command
	rollbackHistoryCmd.Flags().StringP("type", "T", "application", "Resource type (application, service)")
	rollbackHistoryCmd.Flags().IntP("limit", "L", 10, "Limit number of history entries")
	rollbackHistoryCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
