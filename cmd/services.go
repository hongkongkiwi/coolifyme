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
	RunE: func(cmd *cobra.Command, _ []string) error {
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
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tTYPE")
		_, _ = fmt.Fprintln(w, "----\t----\t----")

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

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n",
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
	RunE: func(_ *cobra.Command, args []string) error {
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
	RunE: func(_ *cobra.Command, args []string) error {
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
	RunE: func(_ *cobra.Command, args []string) error {
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

// servicesCreateCmd represents the services create command
var servicesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create service",
	Long:  "Create a new service",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		serviceType, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		project, _ := cmd.Flags().GetString("project")
		environment, _ := cmd.Flags().GetString("environment")
		server, _ := cmd.Flags().GetString("server")
		dockerCompose, _ := cmd.Flags().GetString("docker-compose")
		instantDeploy, _ := cmd.Flags().GetBool("instant-deploy")

		// Validate required fields
		if project == "" {
			return fmt.Errorf("project UUID is required (--project)")
		}
		if server == "" {
			return fmt.Errorf("server UUID is required (--server)")
		}
		if environment == "" {
			return fmt.Errorf("environment name is required (--environment)")
		}

		// Create request body
		req := coolify.CreateServiceJSONRequestBody{
			ProjectUuid:     project,
			ServerUuid:      server,
			EnvironmentName: environment,
			InstantDeploy:   &instantDeploy,
		}

		if serviceType != "" {
			serviceTypeEnum := coolify.CreateServiceJSONBodyType(serviceType)
			req.Type = &serviceTypeEnum
		}
		if name != "" {
			req.Name = &name
		}
		if description != "" {
			req.Description = &description
		}
		if dockerCompose != "" {
			req.DockerComposeRaw = &dockerCompose
		}

		ctx := context.Background()
		uuid, err := client.Services().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}

		fmt.Printf("✅ Service created successfully\n")
		fmt.Printf("   📦 UUID: %s\n", uuid)
		return nil
	},
}

// servicesDeleteCmd represents the services delete command
var servicesDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete service",
	Long:  "Delete a service by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		deleteConfigurations, _ := cmd.Flags().GetBool("delete-configurations")
		deleteVolumes, _ := cmd.Flags().GetBool("delete-volumes")
		dockerCleanup, _ := cmd.Flags().GetBool("docker-cleanup")
		deleteConnectedNetworks, _ := cmd.Flags().GetBool("delete-connected-networks")

		options := &coolify.DeleteServiceByUuidParams{
			DeleteConfigurations:    &deleteConfigurations,
			DeleteVolumes:           &deleteVolumes,
			DockerCleanup:           &dockerCleanup,
			DeleteConnectedNetworks: &deleteConnectedNetworks,
		}

		ctx := context.Background()
		serviceUUID := args[0]

		err = client.Services().Delete(ctx, serviceUUID, options)
		if err != nil {
			return fmt.Errorf("failed to delete service: %w", err)
		}

		fmt.Printf("✅ Service %s deleted successfully\n", serviceUUID)
		return nil
	},
}

// servicesUpdateCmd represents the services update command
var servicesUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update service",
	Long:  "Update a service by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values - only set fields that were explicitly provided
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		dockerCompose, _ := cmd.Flags().GetString("docker-compose")

		// Create request body with only provided fields
		req := coolify.UpdateServiceByUuidJSONRequestBody{}

		// Only set fields if they were provided
		if cmd.Flags().Changed("name") {
			req.Name = &name
		}
		if cmd.Flags().Changed("description") {
			req.Description = &description
		}
		if cmd.Flags().Changed("docker-compose") {
			req.DockerComposeRaw = dockerCompose
		}

		ctx := context.Background()
		serviceUUID := args[0]

		uuid, err := client.Services().Update(ctx, serviceUUID, req)
		if err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}

		fmt.Printf("✅ Service updated successfully\n")
		fmt.Printf("   📦 UUID: %s\n", uuid)
		return nil
	},
}

// servicesListEnvsCmd represents the services list-envs command
var servicesListEnvsCmd = &cobra.Command{
	Use:   "list-envs <uuid>",
	Short: "List environment variables",
	Long:  "List all environment variables for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serviceUUID := args[0]

		envs, err := client.Services().ListEnvs(ctx, serviceUUID)
		if err != nil {
			return fmt.Errorf("failed to list environment variables: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(envs, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(envs) == 0 {
			fmt.Println("No environment variables found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "KEY\tVALUE")
		_, _ = fmt.Fprintln(w, "---\t-----")

		// Print environment variables
		for _, env := range envs {
			key := ""
			value := ""

			if env.Key != nil {
				key = *env.Key
			}
			if env.Value != nil {
				value = *env.Value
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\n", key, value)
		}

		return nil
	},
}

// servicesCreateEnvCmd represents the services create-env command
var servicesCreateEnvCmd = &cobra.Command{
	Use:   "create-env <uuid>",
	Short: "Create environment variable",
	Long:  "Create an environment variable for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")
		isPreview, _ := cmd.Flags().GetBool("is-preview")
		isBuildTime, _ := cmd.Flags().GetBool("is-build-time")
		isLiteral, _ := cmd.Flags().GetBool("is-literal")
		isMultiline, _ := cmd.Flags().GetBool("is-multiline")
		isShownOnce, _ := cmd.Flags().GetBool("is-shown-once")

		// Validate required fields
		if key == "" {
			return fmt.Errorf("environment variable key is required (--key)")
		}
		if value == "" {
			return fmt.Errorf("environment variable value is required (--value)")
		}

		// Create request body
		req := coolify.CreateEnvByServiceUuidJSONRequestBody{
			Key:         &key,
			Value:       &value,
			IsPreview:   &isPreview,
			IsBuildTime: &isBuildTime,
			IsLiteral:   &isLiteral,
			IsMultiline: &isMultiline,
			IsShownOnce: &isShownOnce,
		}

		ctx := context.Background()
		serviceUUID := args[0]

		uuid, err := client.Services().CreateEnv(ctx, serviceUUID, req)
		if err != nil {
			return fmt.Errorf("failed to create environment variable: %w", err)
		}

		fmt.Printf("✅ Environment variable created successfully\n")
		fmt.Printf("   🔑 Key: %s\n", key)
		fmt.Printf("   📦 UUID: %s\n", uuid)
		return nil
	},
}

// servicesUpdateEnvCmd represents the services update-env command
var servicesUpdateEnvCmd = &cobra.Command{
	Use:   "update-env <uuid>",
	Short: "Update environment variable",
	Long:  "Update an environment variable for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")
		isPreview, _ := cmd.Flags().GetBool("is-preview")
		isBuildTime, _ := cmd.Flags().GetBool("is-build-time")
		isLiteral, _ := cmd.Flags().GetBool("is-literal")
		isMultiline, _ := cmd.Flags().GetBool("is-multiline")
		isShownOnce, _ := cmd.Flags().GetBool("is-shown-once")

		// Validate required fields
		if key == "" {
			return fmt.Errorf("environment variable key is required (--key)")
		}
		if value == "" {
			return fmt.Errorf("environment variable value is required (--value)")
		}

		// Create request body
		req := coolify.UpdateEnvByServiceUuidJSONRequestBody{
			Key:         key,
			Value:       value,
			IsPreview:   &isPreview,
			IsBuildTime: &isBuildTime,
			IsLiteral:   &isLiteral,
			IsMultiline: &isMultiline,
			IsShownOnce: &isShownOnce,
		}

		ctx := context.Background()
		serviceUUID := args[0]

		uuid, err := client.Services().UpdateEnv(ctx, serviceUUID, req)
		if err != nil {
			return fmt.Errorf("failed to update environment variable: %w", err)
		}

		fmt.Printf("✅ Environment variable updated successfully\n")
		fmt.Printf("   🔑 Key: %s\n", key)
		fmt.Printf("   📦 UUID: %s\n", uuid)
		return nil
	},
}

// servicesUpdateEnvsCmd represents the services update-envs command
var servicesUpdateEnvsCmd = &cobra.Command{
	Use:   "update-envs <uuid>",
	Short: "Bulk update environment variables",
	Long:  "Update multiple environment variables for a service from a file or JSON string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get flag values
		envDataFlag, _ := cmd.Flags().GetString("env-data")
		envFile, _ := cmd.Flags().GetString("env-file")

		if envDataFlag == "" && envFile == "" {
			return fmt.Errorf("either --env-data or --env-file is required")
		}

		var envVarsList []interface{}
		if envFile != "" {
			// Read environment variables from file
			content, err := os.ReadFile(envFile)
			if err != nil {
				return fmt.Errorf("failed to read env file: %w", err)
			}
			if err := json.Unmarshal(content, &envVarsList); err != nil {
				return fmt.Errorf("failed to parse env file JSON: %w", err)
			}
		} else {
			// Parse environment variables from JSON string
			if err := json.Unmarshal([]byte(envDataFlag), &envVarsList); err != nil {
				return fmt.Errorf("failed to parse env data JSON: %w", err)
			}
		}

		// Convert to the expected structure
		var envStructs []struct {
			IsBuildTime *bool   `json:"is_build_time,omitempty"`
			IsLiteral   *bool   `json:"is_literal,omitempty"`
			IsMultiline *bool   `json:"is_multiline,omitempty"`
			IsPreview   *bool   `json:"is_preview,omitempty"`
			IsShownOnce *bool   `json:"is_shown_once,omitempty"`
			Key         *string `json:"key,omitempty"`
			Value       *string `json:"value,omitempty"`
		}

		// Parse each environment variable
		for _, item := range envVarsList {
			itemData, _ := json.Marshal(item)
			var envVar struct {
				IsBuildTime *bool   `json:"is_build_time,omitempty"`
				IsLiteral   *bool   `json:"is_literal,omitempty"`
				IsMultiline *bool   `json:"is_multiline,omitempty"`
				IsPreview   *bool   `json:"is_preview,omitempty"`
				IsShownOnce *bool   `json:"is_shown_once,omitempty"`
				Key         *string `json:"key,omitempty"`
				Value       *string `json:"value,omitempty"`
			}
			if err := json.Unmarshal(itemData, &envVar); err == nil {
				envStructs = append(envStructs, envVar)
			}
		}

		// Create request body
		req := coolify.UpdateEnvsByServiceUuidJSONRequestBody{
			Data: envStructs,
		}

		ctx := context.Background()
		serviceUUID := args[0]

		message, err := client.Services().UpdateEnvs(ctx, serviceUUID, req)
		if err != nil {
			return fmt.Errorf("failed to bulk update environment variables: %w", err)
		}

		fmt.Printf("✅ Environment variables updated successfully\n")
		fmt.Printf("   💬 Message: %s\n", message)
		return nil
	},
}

// servicesDeleteEnvCmd represents the services delete-env command
var servicesDeleteEnvCmd = &cobra.Command{
	Use:   "delete-env <service-uuid> <env-uuid>",
	Short: "Delete environment variable",
	Long:  "Delete an environment variable from a service",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		serviceUUID := args[0]
		envUUID := args[1]

		uuid, err := client.Services().DeleteEnv(ctx, serviceUUID, envUUID)
		if err != nil {
			return fmt.Errorf("failed to delete environment variable: %w", err)
		}

		fmt.Printf("✅ Environment variable deleted successfully\n")
		fmt.Printf("   📦 UUID: %s\n", uuid)
		return nil
	},
}

func init() {
	// Add subcommands to services
	servicesCmd.AddCommand(servicesListCmd)
	servicesCmd.AddCommand(servicesGetCmd)
	servicesCmd.AddCommand(servicesCreateCmd)
	servicesCmd.AddCommand(servicesUpdateCmd)
	servicesCmd.AddCommand(servicesDeleteCmd)
	servicesCmd.AddCommand(servicesStartCmd)
	servicesCmd.AddCommand(servicesStopCmd)
	servicesCmd.AddCommand(servicesRestartCmd)
	servicesCmd.AddCommand(servicesListEnvsCmd)
	servicesCmd.AddCommand(servicesCreateEnvCmd)
	servicesCmd.AddCommand(servicesUpdateEnvCmd)
	servicesCmd.AddCommand(servicesUpdateEnvsCmd)
	servicesCmd.AddCommand(servicesDeleteEnvCmd)

	// Flags for services list command
	servicesListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for services get command
	servicesGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for services create command
	servicesCreateCmd.Flags().StringP("project", "p", "", "Project UUID (required)")
	servicesCreateCmd.Flags().StringP("server", "s", "", "Server UUID (required)")
	servicesCreateCmd.Flags().StringP("environment", "e", "", "Environment name (required)")
	servicesCreateCmd.Flags().StringP("type", "t", "", "Service type")
	servicesCreateCmd.Flags().StringP("name", "n", "", "Service name")
	servicesCreateCmd.Flags().StringP("description", "d", "", "Service description")
	servicesCreateCmd.Flags().StringP("docker-compose", "c", "", "Docker compose file content")
	servicesCreateCmd.Flags().BoolP("instant-deploy", "i", false, "Deploy service immediately after creation")
	_ = servicesCreateCmd.MarkFlagRequired("project")
	_ = servicesCreateCmd.MarkFlagRequired("server")
	_ = servicesCreateCmd.MarkFlagRequired("environment")

	// Flags for services update command
	servicesUpdateCmd.Flags().StringP("name", "n", "", "Service name")
	servicesUpdateCmd.Flags().StringP("description", "d", "", "Service description")
	servicesUpdateCmd.Flags().StringP("docker-compose", "c", "", "Docker compose file content")
	servicesUpdateCmd.Flags().BoolP("instant-deploy", "i", false, "Deploy service immediately after update")

	// Flags for services delete command
	servicesDeleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")

	// Flags for environment variable list command
	servicesListEnvsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for environment variable create command
	servicesCreateEnvCmd.Flags().StringP("key", "k", "", "Environment variable key (required)")
	servicesCreateEnvCmd.Flags().StringP("value", "v", "", "Environment variable value (required)")
	servicesCreateEnvCmd.Flags().BoolP("is-preview", "p", false, "Is preview environment variable")
	servicesCreateEnvCmd.Flags().BoolP("is-build-time", "b", false, "Is build time environment variable")
	servicesCreateEnvCmd.Flags().BoolP("is-literal", "l", false, "Is literal environment variable")
	servicesCreateEnvCmd.Flags().BoolP("is-multiline", "m", false, "Is multiline environment variable")
	servicesCreateEnvCmd.Flags().BoolP("is-shown-once", "o", false, "Is shown once environment variable")
	_ = servicesCreateEnvCmd.MarkFlagRequired("key")
	_ = servicesCreateEnvCmd.MarkFlagRequired("value")

	// Flags for environment variable update command
	servicesUpdateEnvCmd.Flags().StringP("key", "k", "", "Environment variable key (required)")
	servicesUpdateEnvCmd.Flags().StringP("value", "v", "", "Environment variable value (required)")
	servicesUpdateEnvCmd.Flags().BoolP("is-preview", "p", false, "Is preview environment variable")
	servicesUpdateEnvCmd.Flags().BoolP("is-build-time", "b", false, "Is build time environment variable")
	servicesUpdateEnvCmd.Flags().BoolP("is-literal", "l", false, "Is literal environment variable")
	servicesUpdateEnvCmd.Flags().BoolP("is-multiline", "m", false, "Is multiline environment variable")
	servicesUpdateEnvCmd.Flags().BoolP("is-shown-once", "o", false, "Is shown once environment variable")
	_ = servicesUpdateEnvCmd.MarkFlagRequired("key")
	_ = servicesUpdateEnvCmd.MarkFlagRequired("value")

	// Flags for bulk environment variable update command
	servicesUpdateEnvsCmd.Flags().StringP("env-data", "d", "", "JSON string containing environment variables")
	servicesUpdateEnvsCmd.Flags().StringP("env-file", "f", "", "File containing environment variables in JSON format")
}
