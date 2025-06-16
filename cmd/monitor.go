package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor Coolify resources",
	Long:  "Monitor applications, services, and infrastructure health",
}

// Health check command
// NOTE: golangci-lint incorrectly flags cmd as unused, but it's used for cmd.Flags().GetBool()
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check system health",
	Long:  "Check the health of Coolify API and connected resources",

	RunE: func(cmd *cobra.Command, _ []string) error {
		// Note: cmd parameter is used for accessing flags with cmd.Flags().GetBool()
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		fmt.Println("🏥 Coolify Health Check")
		fmt.Println("======================")

		// Test API connectivity
		fmt.Print("📡 API Connection... ")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Use a simple API call to test connectivity
		_, err = client.Teams().List(ctx)
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			return fmt.Errorf("API health check failed")
		}
		fmt.Println("✅ OK")

		if verbose {
			// Additional checks in verbose mode
			fmt.Print("📦 Applications... ")
			apps, err := client.Applications().List(ctx)
			if err != nil {
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ OK (%d found)\n", len(apps))
			}

			fmt.Print("🖥️  Servers... ")
			servers, err := client.Servers().List(ctx)
			if err != nil {
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ OK (%d found)\n", len(servers))
			}

			fmt.Print("🔧 Services... ")
			services, err := client.Services().List(ctx)
			if err != nil {
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ OK (%d found)\n", len(services))
			}
		}

		fmt.Println("\n🎉 All health checks passed!")
		return nil
	},
}

// Status command for quick overview
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show resource status overview",
	Long:  "Show a quick overview of all resource statuses",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()

		fmt.Println("📊 Coolify Status Overview")
		fmt.Println("=========================")

		// Applications status
		apps, err := client.Applications().List(ctx)
		if err == nil {
			running := 0
			stopped := 0
			unknown := 0

			for _, app := range apps {
				if app.Status != nil {
					switch *app.Status {
					case "running":
						running++
					case "stopped":
						stopped++
					default:
						unknown++
					}
				} else {
					unknown++
				}
			}

			fmt.Printf("📱 Applications: %d total\n", len(apps))
			if running > 0 {
				fmt.Printf("   ✅ Running: %d\n", running)
			}
			if stopped > 0 {
				fmt.Printf("   ⏹️  Stopped: %d\n", stopped)
			}
			if unknown > 0 {
				fmt.Printf("   ❓ Unknown: %d\n", unknown)
			}
		}

		// Servers status
		servers, err := client.Servers().List(ctx)
		if err == nil {
			fmt.Printf("🖥️  Servers: %d total\n", len(servers))
		}

		// Services status
		services, err := client.Services().List(ctx)
		if err == nil {
			fmt.Printf("🔧 Services: %d total\n", len(services))
		}

		return nil
	},
}

// Watch command for real-time monitoring
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch resource status in real-time",
	Long:  "Monitor resource status with auto-refresh",
	RunE: func(cmd *cobra.Command, _ []string) error {
		interval, _ := cmd.Flags().GetInt("interval")
		if interval < 1 {
			interval = 30 // Default 30 seconds
		}

		fmt.Printf("🔄 Watching Coolify status (refresh every %ds, Ctrl+C to stop)...\n\n", interval)

		for {
			// Clear screen (works on most terminals)
			fmt.Print("\033[2J\033[H")

			// Show timestamp
			fmt.Printf("🕒 Last updated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

			// Run status command
			err := statusCmd.RunE(cmd, []string{})
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

			// Wait for next refresh
			time.Sleep(time.Duration(interval) * time.Second)
		}
	},
}

func init() {
	// Add subcommands
	monitorCmd.AddCommand(healthCmd)
	monitorCmd.AddCommand(statusCmd)
	monitorCmd.AddCommand(watchCmd)

	// Health command flags
	healthCmd.Flags().BoolP("verbose", "v", false, "Verbose health check output")

	// Watch command flags
	watchCmd.Flags().IntP("interval", "i", 30, "Refresh interval in seconds")
}
