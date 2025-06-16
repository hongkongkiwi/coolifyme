package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Common aliases for frequently used commands
var (
	// Quick deployment aliases
	deployAppCmd = &cobra.Command{
		Use:     "deploy-app <uuid>",
		Aliases: []string{"deploy", "dep"},
		Short:   "Deploy application (alias for deploy application)",
		Long:    "Quick deployment alias for deploy application command",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Forward to the actual deploy application command
			return deployApplicationCmd().RunE(cmd, args)
		},
	}

	// Quick status check aliases
	quickStatusCmd = &cobra.Command{
		Use:     "status",
		Aliases: []string{"st", "stat"},
		Short:   "Quick status overview (alias for monitor status)",
		Long:    "Show a quick overview of all resource statuses",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Forward to the monitor status command
			return statusCmd.RunE(cmd, args)
		},
	}

	// Quick health check alias
	quickHealthCmd = &cobra.Command{
		Use:     "health",
		Aliases: []string{"ping", "check"},
		Short:   "Quick health check (alias for monitor health)",
		Long:    "Check the health of Coolify API and connected resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Forward to the monitor health command
			return healthCmd.RunE(cmd, args)
		},
	}

	// List all aliases
	listAliasesCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available aliases",
		Long:  "Display all available command aliases and their targets",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("📝 Available Command Aliases")
			fmt.Println("===========================")
			fmt.Println()

			fmt.Println("🚀 Deployment:")
			fmt.Println("   deploy-app, deploy, dep  → deploy application <uuid>")
			fmt.Println()

			fmt.Println("📊 Monitoring:")
			fmt.Println("   status, st, stat         → monitor status")
			fmt.Println("   health, ping, check      → monitor health")
			fmt.Println()

			fmt.Println("📱 Applications:")
			fmt.Println("   apps, app                → applications")
			fmt.Println("   ls-apps                  → applications list")
			fmt.Println()

			fmt.Println("🖥️  Servers:")
			fmt.Println("   servers, server, srv     → servers")
			fmt.Println("   ls-servers               → servers list")
			fmt.Println()

			fmt.Println("🔧 Services:")
			fmt.Println("   services, service, svc   → services")
			fmt.Println("   ls-services              → services list")
			fmt.Println()

			fmt.Println("💡 Tip: Use 'coolifyme <alias> --help' for more information about any command")

			return nil
		},
	}

	// Container for all alias commands
	aliasCmd = &cobra.Command{
		Use:   "alias",
		Short: "Command aliases and shortcuts",
		Long:  "Manage and view command aliases for frequently used operations",
	}

	// Quick list commands
	lsAppsCmd = &cobra.Command{
		Use:     "ls-apps",
		Aliases: []string{"la"},
		Short:   "List applications (alias for applications list)",
		Long:    "Quick alias to list all applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applicationsListCmd.RunE(cmd, args)
		},
	}

	lsServersCmd = &cobra.Command{
		Use:     "ls-servers",
		Aliases: []string{"ls"},
		Short:   "List servers (alias for servers list)",
		Long:    "Quick alias to list all servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serversListCmd.RunE(cmd, args)
		},
	}

	lsServicesCmd = &cobra.Command{
		Use:     "ls-services",
		Aliases: []string{"lsv"},
		Short:   "List services (alias for services list)",
		Long:    "Quick alias to list all services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return servicesListCmd.RunE(cmd, args)
		},
	}
)

func init() {
	// Add alias management commands
	aliasCmd.AddCommand(listAliasesCmd)

	// Copy flags from original commands to aliases where needed
	deployAppCmd.Flags().BoolP("force", "f", false, "Force deployment without confirmation")
	deployAppCmd.Flags().Bool("debug", false, "Enable debug mode for deployment")

	quickHealthCmd.Flags().BoolP("verbose", "v", false, "Verbose health check output")

	// Copy JSON flags for list commands
	lsAppsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	lsServersCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	lsServicesCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
