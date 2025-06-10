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

// privateKeysCmd represents the private keys command
var privateKeysCmd = &cobra.Command{
	Use:     "keys",
	Aliases: []string{"key", "private-keys", "pk"},
	Short:   "Manage private keys",
	Long:    "Manage Coolify private keys - list, create, get, update, and delete private keys",
}

// privateKeysListCmd represents the private keys list command
var privateKeysListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List private keys",
	Long:    "List all private keys in your Coolify instance",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		keys, err := client.PrivateKeys().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list private keys: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(keys, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		if len(keys) == 0 {
			fmt.Println("No private keys found")
			return nil
		}

		// Create a tabwriter for nicely formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer func() {
			_ = w.Flush()
		}()

		// Print header
		_, _ = fmt.Fprintln(w, "UUID\tNAME\tDESCRIPTION\tFINGERPRINT")
		_, _ = fmt.Fprintln(w, "----\t----\t-----------\t-----------")

		// Print private keys
		for _, key := range keys {
			uuid := ""
			name := ""
			description := ""
			fingerprint := ""

			if key.Uuid != nil {
				uuid = *key.Uuid
			}
			if key.Name != nil {
				name = *key.Name
			}
			if key.Description != nil {
				description = *key.Description
			}
			if key.Fingerprint != nil {
				fingerprint = *key.Fingerprint
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", uuid, name, description, fingerprint)
		}

		return nil
	},
}

// privateKeysGetCmd represents the private keys get command
var privateKeysGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get private key details",
	Long:  "Get detailed information about a specific private key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		keyUUID := args[0]

		key, err := client.PrivateKeys().Get(ctx, keyUUID)
		if err != nil {
			return fmt.Errorf("failed to get private key: %w", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(key, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Display private key details in a readable format
		fmt.Printf("Private Key Details:\n")
		fmt.Printf("===================\n")
		if key.Uuid != nil {
			fmt.Printf("UUID:         %s\n", *key.Uuid)
		}
		if key.Name != nil {
			fmt.Printf("Name:         %s\n", *key.Name)
		}
		if key.Description != nil {
			fmt.Printf("Description:  %s\n", *key.Description)
		}
		if key.Fingerprint != nil {
			fmt.Printf("Fingerprint:  %s\n", *key.Fingerprint)
		}
		if key.PublicKey != nil {
			fmt.Printf("Public Key:   %s\n", *key.PublicKey)
		}

		return nil
	},
}

// privateKeysCreateCmd represents the private keys create command
var privateKeysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create private key",
	Long:  "Create a new private key",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		privateKey, _ := cmd.Flags().GetString("private-key")

		if privateKey == "" {
			return fmt.Errorf("private key content is required")
		}

		req := coolify.CreatePrivateKeyJSONRequestBody{
			Name:        &name,
			Description: &description,
			PrivateKey:  privateKey,
		}

		result, err := client.PrivateKeys().Create(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create private key: %w", err)
		}

		fmt.Printf("✅ Private key created successfully\n")
		fmt.Printf("   UUID: %s\n", result)

		return nil
	},
}

// privateKeysUpdateCmd represents the private keys update command
var privateKeysUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update private key",
	Long:  "Update an existing private key",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		privateKey, _ := cmd.Flags().GetString("private-key")

		if privateKey == "" {
			return fmt.Errorf("private key content is required")
		}

		req := coolify.UpdatePrivateKeyJSONRequestBody{
			Name:        &name,
			Description: &description,
			PrivateKey:  privateKey,
		}

		result, err := client.PrivateKeys().Update(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to update private key: %w", err)
		}

		fmt.Printf("✅ Private key updated successfully\n")
		fmt.Printf("   UUID: %s\n", result)

		return nil
	},
}

// privateKeysDeleteCmd represents the private keys delete command
var privateKeysDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete private key",
	Long:  "Delete a private key by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx := context.Background()
		keyUUID := args[0]

		err = client.PrivateKeys().Delete(ctx, keyUUID)
		if err != nil {
			return fmt.Errorf("failed to delete private key: %w", err)
		}

		fmt.Printf("✅ Private key %s deleted successfully\n", keyUUID)
		return nil
	},
}

func init() {
	// Add subcommands to private keys
	privateKeysCmd.AddCommand(privateKeysListCmd)
	privateKeysCmd.AddCommand(privateKeysGetCmd)
	privateKeysCmd.AddCommand(privateKeysCreateCmd)
	privateKeysCmd.AddCommand(privateKeysUpdateCmd)
	privateKeysCmd.AddCommand(privateKeysDeleteCmd)

	// Flags for list command
	privateKeysListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for get command
	privateKeysGetCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Flags for create command
	privateKeysCreateCmd.Flags().StringP("name", "n", "", "Name of the private key")
	privateKeysCreateCmd.Flags().StringP("description", "d", "", "Description of the private key")
	privateKeysCreateCmd.Flags().StringP("private-key", "k", "", "Private key content (required)")
	_ = privateKeysCreateCmd.MarkFlagRequired("private-key")

	// Flags for update command
	privateKeysUpdateCmd.Flags().StringP("name", "n", "", "Name of the private key")
	privateKeysUpdateCmd.Flags().StringP("description", "d", "", "Description of the private key")
	privateKeysUpdateCmd.Flags().StringP("private-key", "k", "", "Private key content (required)")
	_ = privateKeysUpdateCmd.MarkFlagRequired("private-key")
}
