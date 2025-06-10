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
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		servers, err := client.Servers().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list servers: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(servers, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(servers) == 0 {
			fmt.Println("No servers found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tIP\tSTATUS\tDESCRIPTION")
		_, _ = fmt.Fprintln(w, "----\t----\t--\t------\t-----------")

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

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", uuid, name, ip, status, description)
		}

		return nil
	},
}

// serversCreateCmd represents the servers create command
var serversCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create server",
	Long:  "Create a new server in your Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		ip, _ := cmd.Flags().GetString("ip")
		user, _ := cmd.Flags().GetString("user")
		port, _ := cmd.Flags().GetInt32("port")
		privateKeyUuid, _ := cmd.Flags().GetString("private-key-uuid")

		// Validate required fields
		if name == "" {
			return fmt.Errorf("server name is required (--name)")
		}
		if ip == "" {
			return fmt.Errorf("server IP is required (--ip)")
		}
		if user == "" {
			return fmt.Errorf("server user is required (--user)")
		}
		if privateKeyUuid == "" {
			return fmt.Errorf("private key UUID is required (--private-key-uuid)")
		}

		// Create request body
		portInt := int(port)
		req := coolify.CreateServerJSONRequestBody{
			Name:           &name,
			Description:    &description,
			Ip:             &ip,
			User:           &user,
			Port:           &portInt,
			PrivateKeyUuid: &privateKeyUuid,
		}

		ctx := context.Background()

		uuid, err := client.Servers().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		fmt.Printf("✅ Server created successfully\n")
		fmt.Printf("   📛 Name: %s\n", name)
		fmt.Printf("   📦 UUID: %s\n", uuid)
		fmt.Printf("   🌐 IP: %s\n", ip)
		return nil
	},
}

// serversGetCmd represents the servers get command
var serversGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get server details",
	Long:  "Get detailed information about a specific server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serverUUID := args[0]

		server, err := client.Servers().Get(ctx, serverUUID)
		if err != nil {
			return fmt.Errorf("failed to get server: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display server details
		fmt.Printf("📄 Server Details\n")
		fmt.Printf("================\n")

		if server.Uuid != nil {
			fmt.Printf("📦 UUID: %s\n", *server.Uuid)
		}
		if server.Name != nil {
			fmt.Printf("📛 Name: %s\n", *server.Name)
		}
		if server.Description != nil && *server.Description != "" {
			fmt.Printf("📝 Description: %s\n", *server.Description)
		}
		if server.Ip != nil {
			fmt.Printf("🌐 IP: %s\n", *server.Ip)
		}
		if server.User != nil {
			fmt.Printf("👤 User: %s\n", *server.User)
		}
		if server.Port != nil {
			fmt.Printf("🔌 Port: %d\n", *server.Port)
		}
		// Note: Private key UUID is not returned by the API for security reasons

		return nil
	},
}

// serversUpdateCmd represents the servers update command
var serversUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update server",
	Long:  "Update an existing server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		ip, _ := cmd.Flags().GetString("ip")
		user, _ := cmd.Flags().GetString("user")
		port, _ := cmd.Flags().GetInt32("port")
		privateKeyUuid, _ := cmd.Flags().GetString("private-key-uuid")

		// Create request body with only provided values
		req := coolify.UpdateServerByUuidJSONRequestBody{}

		if name != "" {
			req.Name = &name
		}
		if description != "" {
			req.Description = &description
		}
		if ip != "" {
			req.Ip = &ip
		}
		if user != "" {
			req.User = &user
		}
		if cmd.Flags().Changed("port") {
			portInt := int(port)
			req.Port = &portInt
		}
		if privateKeyUuid != "" {
			req.PrivateKeyUuid = &privateKeyUuid
		}

		ctx := context.Background()
		serverUUID := args[0]

		server, err := client.Servers().Update(ctx, serverUUID, req)
		if err != nil {
			return fmt.Errorf("failed to update server: %w", err)
		}

		fmt.Printf("✅ Server updated successfully\n")
		if server.Uuid != nil {
			fmt.Printf("   📦 UUID: %s\n", *server.Uuid)
		}
		if server.Name != nil {
			fmt.Printf("   📛 Name: %s\n", *server.Name)
		}
		return nil
	},
}

// serversDeleteCmd represents the servers delete command
var serversDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete server",
	Long:  "Delete a server from your Coolify instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serverUUID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("⚠️  Are you sure you want to delete server %s? This action cannot be undone.\n", serverUUID)
			fmt.Print("Type 'yes' to confirm: ")
			var confirmation string
			if _, err := fmt.Scanln(&confirmation); err != nil || confirmation != "yes" {
				fmt.Println("❌ Deletion cancelled")
				return nil
			}
		}

		err = client.Servers().Delete(ctx, serverUUID)
		if err != nil {
			return fmt.Errorf("failed to delete server: %w", err)
		}

		fmt.Printf("✅ Server %s deleted successfully\n", serverUUID)
		return nil
	},
}

// serversGetResourcesCmd represents the servers get-resources command
var serversGetResourcesCmd = &cobra.Command{
	Use:   "get-resources <uuid>",
	Short: "Get server resources",
	Long:  "Get resource information for a specific server",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serverUUID := args[0]

		resources, err := client.Servers().GetResources(ctx, serverUUID)
		if err != nil {
			return fmt.Errorf("failed to get server resources: %w", err)
		}

		fmt.Printf("📊 Server Resources\n")
		fmt.Printf("==================\n")
		fmt.Printf("%s\n", resources)
		return nil
	},
}

// serversGetDomainsCmd represents the servers get-domains command
var serversGetDomainsCmd = &cobra.Command{
	Use:   "get-domains <uuid>",
	Short: "Get server domains",
	Long:  "Get domain information for a specific server",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serverUUID := args[0]

		domains, err := client.Servers().GetDomains(ctx, serverUUID)
		if err != nil {
			return fmt.Errorf("failed to get server domains: %w", err)
		}

		fmt.Printf("🌐 Server Domains\n")
		fmt.Printf("================\n")
		fmt.Printf("%s\n", domains)
		return nil
	},
}

// serversValidateCmd represents the servers validate command
var serversValidateCmd = &cobra.Command{
	Use:   "validate <uuid>",
	Short: "Validate server",
	Long:  "Validate server connection and configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serverUUID := args[0]

		result, err := client.Servers().Validate(ctx, serverUUID)
		if err != nil {
			return fmt.Errorf("failed to validate server: %w", err)
		}

		fmt.Printf("✅ Server Validation\n")
		fmt.Printf("===================\n")
		fmt.Printf("%s\n", result)
		return nil
	},
}

func init() {
	// Add subcommands to servers
	serversCmd.AddCommand(serversListCmd)
	serversCmd.AddCommand(serversCreateCmd)
	serversCmd.AddCommand(serversGetCmd)
	serversCmd.AddCommand(serversUpdateCmd)
	serversCmd.AddCommand(serversDeleteCmd)
	serversCmd.AddCommand(serversGetResourcesCmd)
	serversCmd.AddCommand(serversGetDomainsCmd)
	serversCmd.AddCommand(serversValidateCmd)

	// Flags for servers list command
	serversListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for servers create command
	serversCreateCmd.Flags().StringP("name", "n", "", "Server name (required)")
	serversCreateCmd.Flags().StringP("description", "d", "", "Server description")
	serversCreateCmd.Flags().StringP("ip", "i", "", "Server IP address (required)")
	serversCreateCmd.Flags().StringP("user", "u", "", "SSH user (required)")
	serversCreateCmd.Flags().Int32P("port", "p", 22, "SSH port")
	serversCreateCmd.Flags().StringP("private-key-uuid", "k", "", "Private key UUID (required)")
	_ = serversCreateCmd.MarkFlagRequired("name")
	_ = serversCreateCmd.MarkFlagRequired("ip")
	_ = serversCreateCmd.MarkFlagRequired("user")
	_ = serversCreateCmd.MarkFlagRequired("private-key-uuid")

	// Flags for servers get command
	serversGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for servers update command
	serversUpdateCmd.Flags().StringP("name", "n", "", "Server name")
	serversUpdateCmd.Flags().StringP("description", "d", "", "Server description")
	serversUpdateCmd.Flags().StringP("ip", "i", "", "Server IP address")
	serversUpdateCmd.Flags().StringP("user", "u", "", "SSH user")
	serversUpdateCmd.Flags().Int32P("port", "p", 22, "SSH port")
	serversUpdateCmd.Flags().StringP("private-key-uuid", "k", "", "Private key UUID")

	// Flags for servers delete command
	serversDeleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
}
