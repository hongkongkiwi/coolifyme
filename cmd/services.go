package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:     "services",
	Aliases: []string{"service", "svc"},
	Short:   "Manage services",
	Long:    "Manage Coolify services - list, get details, start, stop, and restart services",
}

// servicesListCmd represents the services list command
var servicesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List services",
	Long:    "List all services in your Coolify instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		services, err := client.Services().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(services, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(services) == 0 {
			fmt.Println("No services found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer w.Flush()

		// Print header
		fmt.Fprintln(w, "UUID\tNAME\tTYPE")
		fmt.Fprintln(w, "----\t----\t----")

		// Print services
		for _, service := range services {
			uuid := ""
			name := ""
			serviceType := ""

			if service.Uuid != nil {
				uuid = *service.Uuid
			}
			if service.Name != nil {
				name = *service.Name
			}
			if service.ServiceType != nil {
				serviceType = *service.ServiceType
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n",
				uuid, name, serviceType)
		}

		return nil
	},
}

// servicesGetCmd represents the services get command
var servicesGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get service details",
	Long:  "Get detailed information about a specific service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serviceUUID := args[0]

		service, err := client.Services().Get(ctx, serviceUUID)
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(service, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display service details in a readable format
		fmt.Printf("Service Details:\n")
		fmt.Printf("===============\n")
		if service.Uuid != nil {
			fmt.Printf("UUID:           %s\n", *service.Uuid)
		}
		if service.Name != nil {
			fmt.Printf("Name:           %s\n", *service.Name)
		}
		if service.ServiceType != nil {
			fmt.Printf("Type:           %s\n", *service.ServiceType)
		}
		if service.Description != nil {
			fmt.Printf("Description:    %s\n", *service.Description)
		}

		return nil
	},
}

// servicesStartCmd represents the services start command
var servicesStartCmd = &cobra.Command{
	Use:   "start <uuid>",
	Short: "Start service",
	Long:  "Start a service by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serviceUUID := args[0]

		err = client.Services().Start(ctx, serviceUUID)
		if err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}

		fmt.Printf("✅ Service %s start request queued successfully\n", serviceUUID)
		return nil
	},
}

// servicesStopCmd represents the services stop command
var servicesStopCmd = &cobra.Command{
	Use:   "stop <uuid>",
	Short: "Stop service",
	Long:  "Stop a service by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serviceUUID := args[0]

		err = client.Services().Stop(ctx, serviceUUID)
		if err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}

		fmt.Printf("✅ Service %s stop request queued successfully\n", serviceUUID)
		return nil
	},
}

// servicesRestartCmd represents the services restart command
var servicesRestartCmd = &cobra.Command{
	Use:   "restart <uuid>",
	Short: "Restart service",
	Long:  "Restart a service by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serviceUUID := args[0]

		err = client.Services().Restart(ctx, serviceUUID)
		if err != nil {
			return fmt.Errorf("failed to restart service: %w", err)
		}

		fmt.Printf("✅ Service %s restart request queued successfully\n", serviceUUID)
		return nil
	},
}

func init() {
	// Add subcommands to services
	servicesCmd.AddCommand(servicesListCmd)
	servicesCmd.AddCommand(servicesGetCmd)
	servicesCmd.AddCommand(servicesStartCmd)
	servicesCmd.AddCommand(servicesStopCmd)
	servicesCmd.AddCommand(servicesRestartCmd)

	// Flags for services list command
	servicesListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for services get command
	servicesGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
