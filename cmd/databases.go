package main

import (
	"context"
	"fmt"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/spf13/cobra"
)

// databasesCmd represents the databases command
var databasesCmd = &cobra.Command{
	Use:     "databases",
	Aliases: []string{"database", "db"},
	Short:   "Manage databases",
	Long:    "Manage Coolify databases - list, get details, start, stop, and restart databases",
}

// databasesListCmd represents the databases list command
var databasesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List databases",
	Long:    "List all databases in your Coolify instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		result, err := client.Databases().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list databases: %w", err)
		}

		// The database API currently returns a simple string
		fmt.Printf("Databases:\n%s\n", result)
		return nil
	},
}

// databasesGetCmd represents the databases get command
var databasesGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get database details",
	Long:  "Get detailed information about a specific database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		databaseUUID := args[0]

		result, err := client.Databases().Get(ctx, databaseUUID)
		if err != nil {
			return fmt.Errorf("failed to get database: %w", err)
		}

		// The database API currently returns a simple string
		fmt.Printf("Database %s:\n%s\n", databaseUUID, result)
		return nil
	},
}

// databasesStartCmd represents the databases start command
var databasesStartCmd = &cobra.Command{
	Use:   "start <uuid>",
	Short: "Start database",
	Long:  "Start a database by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		databaseUUID := args[0]

		err = client.Databases().Start(ctx, databaseUUID)
		if err != nil {
			return fmt.Errorf("failed to start database: %w", err)
		}

		fmt.Printf("✅ Database %s start request queued successfully\n", databaseUUID)
		return nil
	},
}

// databasesStopCmd represents the databases stop command
var databasesStopCmd = &cobra.Command{
	Use:   "stop <uuid>",
	Short: "Stop database",
	Long:  "Stop a database by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		databaseUUID := args[0]

		err = client.Databases().Stop(ctx, databaseUUID)
		if err != nil {
			return fmt.Errorf("failed to stop database: %w", err)
		}

		fmt.Printf("✅ Database %s stop request queued successfully\n", databaseUUID)
		return nil
	},
}

// databasesRestartCmd represents the databases restart command
var databasesRestartCmd = &cobra.Command{
	Use:   "restart <uuid>",
	Short: "Restart database",
	Long:  "Restart a database by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		databaseUUID := args[0]

		err = client.Databases().Restart(ctx, databaseUUID)
		if err != nil {
			return fmt.Errorf("failed to restart database: %w", err)
		}

		fmt.Printf("✅ Database %s restart request queued successfully\n", databaseUUID)
		return nil
	},
}

// databasesDeleteCmd represents the databases delete command
var databasesDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete a database",
	Long:  "Delete a database by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		deleteVolumes, _ := cmd.Flags().GetBool("delete-volumes")
		deleteConfigs, _ := cmd.Flags().GetBool("delete-configurations")

		options := &coolify.DeleteDatabaseByUuidParams{
			DeleteVolumes:        &deleteVolumes,
			DeleteConfigurations: &deleteConfigs,
		}

		err = client.Databases().Delete(context.Background(), args[0], options)
		if err != nil {
			return fmt.Errorf("failed to delete database: %w", err)
		}

		fmt.Printf("Database %s deleted successfully\n", args[0])
		return nil
	},
}

// databasesUpdateCmd represents the databases update command
var databasesUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update a database",
	Long:  "Update a database by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// For now, this is a stub - you'd need to implement fields to update
		req := coolify.UpdateDatabaseByUuidJSONRequestBody{}

		err = client.Databases().Update(context.Background(), args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update database: %w", err)
		}

		fmt.Printf("Database %s updated successfully\n", args[0])
		return nil
	},
}

// databasesCreateCmd represents the databases create command
var databasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database",
	Long:  "Create a new database of various types",
}

// databasesCreatePostgreSQLCmd represents the databases create postgresql command
var databasesCreatePostgreSQLCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "Create a PostgreSQL database",
	Long:  "Create a new PostgreSQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get required parameters
		envName, _ := cmd.Flags().GetString("environment")
		envUUID, _ := cmd.Flags().GetString("environment-uuid")
		projectUUID, _ := cmd.Flags().GetString("project")
		serverUUID, _ := cmd.Flags().GetString("server")

		if envName == "" && envUUID == "" {
			return fmt.Errorf("either --environment or --environment-uuid is required")
		}
		if projectUUID == "" {
			return fmt.Errorf("--project is required")
		}
		if serverUUID == "" {
			return fmt.Errorf("--server is required")
		}

		req := coolify.CreateDatabasePostgresqlJSONRequestBody{
			EnvironmentName: envName,
			EnvironmentUuid: envUUID,
			ProjectUuid:     projectUUID,
			ServerUuid:      serverUUID,
		}

		// Optional parameters
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
		}
		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			req.Description = &desc
		}
		if image, _ := cmd.Flags().GetString("image"); image != "" {
			req.Image = &image
		}
		if instant, _ := cmd.Flags().GetBool("instant-deploy"); instant {
			req.InstantDeploy = &instant
		}

		err = client.Databases().CreatePostgreSQL(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create PostgreSQL database: %w", err)
		}

		fmt.Println("PostgreSQL database created successfully")
		return nil
	},
}

// databasesCreateMySQLCmd represents the databases create mysql command
var databasesCreateMySQLCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Create a MySQL database",
	Long:  "Create a new MySQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get required parameters
		envName, _ := cmd.Flags().GetString("environment")
		envUUID, _ := cmd.Flags().GetString("environment-uuid")
		projectUUID, _ := cmd.Flags().GetString("project")
		serverUUID, _ := cmd.Flags().GetString("server")

		if envName == "" && envUUID == "" {
			return fmt.Errorf("either --environment or --environment-uuid is required")
		}
		if projectUUID == "" {
			return fmt.Errorf("--project is required")
		}
		if serverUUID == "" {
			return fmt.Errorf("--server is required")
		}

		req := coolify.CreateDatabaseMysqlJSONRequestBody{
			EnvironmentName: envName,
			EnvironmentUuid: envUUID,
			ProjectUuid:     projectUUID,
			ServerUuid:      serverUUID,
		}

		// Optional parameters
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
		}
		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			req.Description = &desc
		}
		if image, _ := cmd.Flags().GetString("image"); image != "" {
			req.Image = &image
		}
		if instant, _ := cmd.Flags().GetBool("instant-deploy"); instant {
			req.InstantDeploy = &instant
		}

		err = client.Databases().CreateMySQL(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create MySQL database: %w", err)
		}

		fmt.Println("MySQL database created successfully")
		return nil
	},
}

// databasesCreateRedisCmd represents the databases create redis command
var databasesCreateRedisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Create a Redis database",
	Long:  "Create a new Redis database",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get required parameters
		envName, _ := cmd.Flags().GetString("environment")
		envUUID, _ := cmd.Flags().GetString("environment-uuid")
		projectUUID, _ := cmd.Flags().GetString("project")
		serverUUID, _ := cmd.Flags().GetString("server")

		if envName == "" && envUUID == "" {
			return fmt.Errorf("either --environment or --environment-uuid is required")
		}
		if projectUUID == "" {
			return fmt.Errorf("--project is required")
		}
		if serverUUID == "" {
			return fmt.Errorf("--server is required")
		}

		req := coolify.CreateDatabaseRedisJSONRequestBody{
			EnvironmentName: envName,
			EnvironmentUuid: envUUID,
			ProjectUuid:     projectUUID,
			ServerUuid:      serverUUID,
		}

		// Optional parameters
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
		}
		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			req.Description = &desc
		}
		if image, _ := cmd.Flags().GetString("image"); image != "" {
			req.Image = &image
		}
		if instant, _ := cmd.Flags().GetBool("instant-deploy"); instant {
			req.InstantDeploy = &instant
		}

		err = client.Databases().CreateRedis(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create Redis database: %w", err)
		}

		fmt.Println("Redis database created successfully")
		return nil
	},
}

// databasesCreateMongoDBCmd represents the databases create mongodb command
var databasesCreateMongoDBCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "Create a MongoDB database",
	Long:  "Create a new MongoDB database",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get required parameters
		envName, _ := cmd.Flags().GetString("environment")
		envUUID, _ := cmd.Flags().GetString("environment-uuid")
		projectUUID, _ := cmd.Flags().GetString("project")
		serverUUID, _ := cmd.Flags().GetString("server")

		if envName == "" && envUUID == "" {
			return fmt.Errorf("either --environment or --environment-uuid is required")
		}
		if projectUUID == "" {
			return fmt.Errorf("--project is required")
		}
		if serverUUID == "" {
			return fmt.Errorf("--server is required")
		}

		req := coolify.CreateDatabaseMongodbJSONRequestBody{
			EnvironmentName: envName,
			EnvironmentUuid: envUUID,
			ProjectUuid:     projectUUID,
			ServerUuid:      serverUUID,
		}

		// Optional parameters
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = &name
		}
		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			req.Description = &desc
		}
		if image, _ := cmd.Flags().GetString("image"); image != "" {
			req.Image = &image
		}
		if instant, _ := cmd.Flags().GetBool("instant-deploy"); instant {
			req.InstantDeploy = &instant
		}

		err = client.Databases().CreateMongoDB(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create MongoDB database: %w", err)
		}

		fmt.Println("MongoDB database created successfully")
		return nil
	},
}

// Create stubs for other database types
var databasesCreateClickHouseCmd = &cobra.Command{
	Use:   "clickhouse",
	Short: "Create a ClickHouse database",
	Long:  "Create a new ClickHouse database",
	RunE:  createDatabaseStub("ClickHouse"),
}

var databasesCreateDragonflyCmd = &cobra.Command{
	Use:   "dragonfly",
	Short: "Create a Dragonfly database",
	Long:  "Create a new Dragonfly database",
	RunE:  createDatabaseStub("Dragonfly"),
}

var databasesCreateKeyDBCmd = &cobra.Command{
	Use:   "keydb",
	Short: "Create a KeyDB database",
	Long:  "Create a new KeyDB database",
	RunE:  createDatabaseStub("KeyDB"),
}

var databasesCreateMariaDBCmd = &cobra.Command{
	Use:   "mariadb",
	Short: "Create a MariaDB database",
	Long:  "Create a new MariaDB database",
	RunE:  createDatabaseStub("MariaDB"),
}

func createDatabaseStub(dbType string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Creating %s databases is not yet implemented in this CLI\n", dbType)
		fmt.Printf("Please use the Coolify web interface or API directly\n")
		return nil
	}
}

func init() {
	// Common flags for all database create commands
	for _, cmd := range []*cobra.Command{
		databasesCreatePostgreSQLCmd,
		databasesCreateMySQLCmd,
		databasesCreateRedisCmd,
		databasesCreateMongoDBCmd,
		databasesCreateClickHouseCmd,
		databasesCreateDragonflyCmd,
		databasesCreateKeyDBCmd,
		databasesCreateMariaDBCmd,
	} {
		cmd.Flags().String("environment", "", "Environment name")
		cmd.Flags().String("environment-uuid", "", "Environment UUID")
		cmd.Flags().String("project", "", "Project UUID (required)")
		cmd.Flags().String("server", "", "Server UUID (required)")
		cmd.Flags().String("name", "", "Database name")
		cmd.Flags().String("description", "", "Database description")
		cmd.Flags().String("image", "", "Docker image")
		cmd.Flags().Bool("instant-deploy", false, "Deploy immediately")
	}

	// Add subcommands to databases
	databasesCmd.AddCommand(databasesListCmd)
	databasesCmd.AddCommand(databasesGetCmd)
	databasesCmd.AddCommand(databasesStartCmd)
	databasesCmd.AddCommand(databasesStopCmd)
	databasesCmd.AddCommand(databasesRestartCmd)
	databasesCmd.AddCommand(databasesDeleteCmd)
	databasesCmd.AddCommand(databasesUpdateCmd)
	databasesCmd.AddCommand(databasesCreateCmd)
}
