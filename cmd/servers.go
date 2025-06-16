package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:     "servers",
	Aliases: []string{"server", "srv"},
	Short:   "Manage servers",
	Long:    "Manage Coolify servers - list, create, update, delete, and monitor servers",
}

// serversListCmd represents the servers list command
var serversListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List servers",
	Long:    "List all servers in your Coolify instance with their current status",
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
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tIP\tPORT\tUSER\tSTATUS\tPROXY\tDESCRIPTION")
		_, _ = fmt.Fprintln(w, "----\t----\t--\t----\t----\t------\t-----\t-----------")

		// Print servers
		for _, server := range servers {
			uuid := ""
			name := ""
			ip := ""
			port := ""
			user := ""
			status := ""
			proxy := ""
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
			if server.Port != nil {
				port = fmt.Sprintf("%d", *server.Port)
			}
			if server.User != nil {
				user = *server.User
			}
			if server.ValidationLogs != nil {
				status = "validated"
			} else {
				status = "unknown"
			}
			// Use the direct ProxyType field
			if server.ProxyType != nil {
				proxy = string(*server.ProxyType)
			}
			if server.Description != nil {
				description = *server.Description
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				uuid, name, ip, port, user, status, proxy, description)
		}

		return nil
	},
}

// serversCreateCmd represents the servers create command
var serversCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create server",
	Long:  "Create a new server in your Coolify instance with advanced configuration options",
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
		isBuildServer, _ := cmd.Flags().GetBool("is-build-server")
		instantValidate, _ := cmd.Flags().GetBool("instant-validate")
		proxyType, _ := cmd.Flags().GetString("proxy-type")

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

		// Validate proxy type if provided
		if proxyType != "" {
			validProxyTypes := []string{"traefik", "caddy", "none"}
			isValid := false
			for _, valid := range validProxyTypes {
				if proxyType == valid {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid proxy type: %s. Valid options: %s", proxyType, strings.Join(validProxyTypes, ", "))
			}
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

		// Add optional fields if they have specific values
		if isBuildServer {
			req.IsBuildServer = &isBuildServer
		}
		if instantValidate {
			req.InstantValidate = &instantValidate
		}
		if proxyType != "" {
			// Convert string to proper enum type
			var proxyTypeEnum coolify.CreateServerJSONBodyProxyType
			switch proxyType {
			case ProxyTraefik:
				proxyTypeEnum = coolify.CreateServerJSONBodyProxyTypeTraefik
			case "caddy":
				proxyTypeEnum = coolify.CreateServerJSONBodyProxyTypeCaddy
			case "none":
				proxyTypeEnum = coolify.CreateServerJSONBodyProxyTypeNone
			}
			req.ProxyType = &proxyTypeEnum
		}

		ctx := context.Background()

		uuid, err := client.Servers().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		fmt.Printf("✅ Server created successfully\n")
		fmt.Printf("   📛 Name: %s\n", name)
		fmt.Printf("   📦 UUID: %s\n", uuid)
		fmt.Printf("   🌐 IP: %s:%d\n", ip, port)
		fmt.Printf("   👤 User: %s\n", user)
		if proxyType != "" {
			fmt.Printf("   🔧 Proxy: %s\n", proxyType)
		}
		if isBuildServer {
			fmt.Printf("   🏗️  Build Server: Yes\n")
		}
		if instantValidate {
			fmt.Printf("   ⚡ Instant Validate: Yes\n")
		}
		return nil
	},
}

// serversGetCmd represents the servers get command
var serversGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get server details",
	Long:  "Get detailed information about a specific server including configuration and status",
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

		// Display proxy type from the direct field
		if server.ProxyType != nil {
			fmt.Printf("🔧 Proxy Type: %s\n", string(*server.ProxyType))
		}

		// Display build server setting from the Settings field
		if server.Settings != nil && server.Settings.IsBuildServer != nil && *server.Settings.IsBuildServer {
			fmt.Printf("🏗️  Build Server: Yes\n")
		}

		// Display validation status
		if server.ValidationLogs != nil {
			fmt.Printf("✅ Status: Validated\n")
		} else {
			fmt.Printf("⚠️  Status: Not validated\n")
		}

		// Display additional server information
		if server.Settings != nil {
			if server.Settings.IsReachable != nil {
				if *server.Settings.IsReachable {
					fmt.Printf("📡 Reachable: Yes\n")
				} else {
					fmt.Printf("📡 Reachable: No\n")
				}
			}
			if server.Settings.IsUsable != nil {
				if *server.Settings.IsUsable {
					fmt.Printf("⚡ Usable: Yes\n")
				} else {
					fmt.Printf("⚡ Usable: No\n")
				}
			}
		}

		// Note: Private key UUID is not returned by the API for security reasons

		return nil
	},
}

// serversUpdateCmd represents the servers update command
var serversUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update server",
	Long:  "Update an existing server configuration including advanced options",
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
		isBuildServer, _ := cmd.Flags().GetBool("is-build-server")
		instantValidate, _ := cmd.Flags().GetBool("instant-validate")
		proxyType, _ := cmd.Flags().GetString("proxy-type")

		// Validate proxy type if provided
		if cmd.Flags().Changed("proxy-type") && proxyType != "" {
			validProxyTypes := []string{"traefik", "caddy", "none"}
			isValid := false
			for _, valid := range validProxyTypes {
				if proxyType == valid {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid proxy type: %s. Valid options: %s", proxyType, strings.Join(validProxyTypes, ", "))
			}
		}

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
		if cmd.Flags().Changed("is-build-server") {
			req.IsBuildServer = &isBuildServer
		}
		if cmd.Flags().Changed("instant-validate") {
			req.InstantValidate = &instantValidate
		}
		if cmd.Flags().Changed("proxy-type") && proxyType != "" {
			// Convert string to proper enum type
			var proxyTypeEnum coolify.UpdateServerByUuidJSONBodyProxyType
			switch proxyType {
			case "traefik":
				proxyTypeEnum = coolify.Traefik
			case "caddy":
				proxyTypeEnum = coolify.Caddy
			case "none":
				proxyTypeEnum = coolify.None
			}
			req.ProxyType = &proxyTypeEnum
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
		if server.Ip != nil && server.Port != nil {
			fmt.Printf("   🌐 IP: %s:%d\n", *server.Ip, *server.Port)
		}
		return nil
	},
}

// serversDeleteCmd represents the servers delete command
var serversDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete server",
	Long:  "Delete a server from your Coolify instance (this action cannot be undone)",
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
	Long:  "Get detailed resource information for a specific server including applications, databases, and services",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Println(resources)
			return nil
		}

		// Parse the JSON response for better formatting
		var resourceData interface{}
		if err := json.Unmarshal([]byte(resources), &resourceData); err != nil {
			// If parsing fails, just display the raw response
			fmt.Printf("📊 Server Resources\n")
			fmt.Printf("==================\n")
			fmt.Printf("%s\n", resources)
			return nil
		}

		// Pretty print the JSON
		prettyJSON, err := json.MarshalIndent(resourceData, "", "  ")
		if err != nil {
			fmt.Printf("📊 Server Resources\n")
			fmt.Printf("==================\n")
			fmt.Printf("%s\n", resources)
			return nil
		}

		fmt.Printf("📊 Server Resources\n")
		fmt.Printf("==================\n")
		fmt.Printf("%s\n", string(prettyJSON))
		return nil
	},
}

// serversGetDomainsCmd represents the servers get-domains command
var serversGetDomainsCmd = &cobra.Command{
	Use:   "get-domains <uuid>",
	Short: "Get server domains",
	Long:  "Get domain configuration and routing information for a specific server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Println(domains)
			return nil
		}

		// Parse the JSON response for better formatting
		var domainData interface{}
		if err := json.Unmarshal([]byte(domains), &domainData); err != nil {
			// If parsing fails, just display the raw response
			fmt.Printf("🌐 Server Domains\n")
			fmt.Printf("================\n")
			fmt.Printf("%s\n", domains)
			return nil
		}

		// Pretty print the JSON
		prettyJSON, err := json.MarshalIndent(domainData, "", "  ")
		if err != nil {
			fmt.Printf("🌐 Server Domains\n")
			fmt.Printf("================\n")
			fmt.Printf("%s\n", domains)
			return nil
		}

		fmt.Printf("🌐 Server Domains\n")
		fmt.Printf("================\n")
		fmt.Printf("%s\n", string(prettyJSON))
		return nil
	},
}

// serversValidateCmd represents the servers validate command
var serversValidateCmd = &cobra.Command{
	Use:   "validate <uuid>",
	Short: "Validate server",
	Long:  "Validate server connection, configuration, and readiness for deployment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output := map[string]interface{}{
				"message":     result,
				"server_uuid": serverUUID,
			}
			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		fmt.Printf("✅ Server Validation\n")
		fmt.Printf("===================\n")
		fmt.Printf("Server: %s\n", serverUUID)
		fmt.Printf("Status: %s\n", result)
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
	serversCreateCmd.Flags().Bool("is-build-server", false, "Configure as build server")
	serversCreateCmd.Flags().Bool("instant-validate", false, "Validate server immediately after creation")
	serversCreateCmd.Flags().String("proxy-type", "", "Proxy type (traefik, caddy, none)")
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
	serversUpdateCmd.Flags().Bool("is-build-server", false, "Configure as build server")
	serversUpdateCmd.Flags().Bool("instant-validate", false, "Validate server after update")
	serversUpdateCmd.Flags().String("proxy-type", "", "Proxy type (traefik, caddy, none)")

	// Flags for servers delete command
	serversDeleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")

	// Flags for servers get-resources command
	serversGetResourcesCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for servers get-domains command
	serversGetDomainsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for servers validate command
	serversValidateCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
