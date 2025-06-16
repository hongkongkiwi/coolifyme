package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hongkongkiwi/coolifyme/internal/config"
	"github.com/spf13/cobra"
)

// Interactive setup wizard
var initInteractiveCmd = &cobra.Command{
	Use:   "init-interactive",
	Short: "Interactive setup wizard",
	Long:  "Guided setup wizard to configure coolifyme for first-time use",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println("🚀 Welcome to coolifyme interactive setup!")
		fmt.Println("=====================================")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		// Profile name
		fmt.Print("📛 Profile name [default]: ")
		profileName, _ := reader.ReadString('\n')
		profileName = strings.TrimSpace(profileName)
		if profileName == "" {
			profileName = "default"
		}

		// API Token
		fmt.Print("🔑 Coolify API Token: ")
		apiToken, _ := reader.ReadString('\n')
		apiToken = strings.TrimSpace(apiToken)
		if apiToken == "" {
			return fmt.Errorf("API token is required")
		}

		// Base URL
		fmt.Print("🌐 Coolify URL [https://app.coolify.io/api/v1]: ")
		baseURL, _ := reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
		if baseURL == "" {
			baseURL = "https://app.coolify.io/api/v1"
		}

		// Output format
		fmt.Print("📄 Default output format (table/json/yaml) [table]: ")
		outputFormat, _ := reader.ReadString('\n')
		outputFormat = strings.TrimSpace(outputFormat)
		if outputFormat == "" {
			outputFormat = "table"
		}

		// Log level
		fmt.Print("📝 Log level (debug/info/warn/error) [info]: ")
		logLevel, _ := reader.ReadString('\n')
		logLevel = strings.TrimSpace(logLevel)
		if logLevel == "" {
			logLevel = "info"
		}

		// Create profile
		fmt.Println("\n⚙️  Creating profile...")

		cfg := &config.Config{
			APIToken:     apiToken,
			BaseURL:      baseURL,
			Profile:      profileName,
			OutputFormat: outputFormat,
			LogLevel:     logLevel,
		}

		if err := config.CreateProfile(profileName, apiToken, baseURL); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}

		if err := config.SetDefaultProfile(profileName); err != nil {
			return fmt.Errorf("failed to set default profile: %w", err)
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Println("✅ Setup completed successfully!")
		fmt.Printf("   📛 Profile: %s\n", profileName)
		fmt.Printf("   🌐 URL: %s\n", baseURL)
		fmt.Printf("   📄 Output: %s\n", outputFormat)
		fmt.Printf("   📝 Log Level: %s\n", logLevel)
		fmt.Println()
		fmt.Println("🎉 You can now use coolifyme! Try: coolifyme apps list")

		return nil
	},
}

// Interactive application creation wizard
var appCreateWizardCmd = &cobra.Command{
	Use:   "create-wizard",
	Short: "Interactive application creation wizard",
	Long:  "Guided wizard to create a new application with all necessary configuration",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println("🚀 Application Creation Wizard")
		fmt.Println("=============================")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("📦 Loading projects and servers...")
		// This would require API calls to list projects and servers
		// For now, we'll ask for UUIDs directly

		// Repository URL
		fmt.Print("📁 Git repository URL: ")
		repo, _ := reader.ReadString('\n')
		repo = strings.TrimSpace(repo)
		if repo == "" {
			return fmt.Errorf("repository URL is required")
		}

		// Branch
		fmt.Print("🌿 Git branch [main]: ")
		branch, _ := reader.ReadString('\n')
		branch = strings.TrimSpace(branch)
		if branch == "" {
			branch = "main"
		}

		// Build pack
		fmt.Print("🏗️  Build pack (nixpacks/static/dockerfile/dockercompose) [nixpacks]: ")
		buildPack, _ := reader.ReadString('\n')
		buildPack = strings.TrimSpace(buildPack)
		if buildPack == "" {
			buildPack = "nixpacks"
		}

		// Project UUID
		fmt.Print("📦 Project UUID: ")
		project, _ := reader.ReadString('\n')
		project = strings.TrimSpace(project)
		if project == "" {
			return fmt.Errorf("project UUID is required")
		}

		// Server UUID
		fmt.Print("🖥️  Server UUID: ")
		server, _ := reader.ReadString('\n')
		server = strings.TrimSpace(server)
		if server == "" {
			return fmt.Errorf("server UUID is required")
		}

		// Environment
		fmt.Print("🌍 Environment [production]: ")
		environment, _ := reader.ReadString('\n')
		environment = strings.TrimSpace(environment)
		if environment == "" {
			environment = "production"
		}

		fmt.Println("\n📋 Configuration Summary:")
		fmt.Printf("   📁 Repository: %s\n", repo)
		fmt.Printf("   🌿 Branch: %s\n", branch)
		fmt.Printf("   🏗️  Build Pack: %s\n", buildPack)
		fmt.Printf("   📦 Project: %s\n", project)
		fmt.Printf("   🖥️  Server: %s\n", server)
		fmt.Printf("   🌍 Environment: %s\n", environment)
		fmt.Println()

		fmt.Print("✅ Create application? (y/N): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))

		if confirm != "y" && confirm != "yes" {
			fmt.Println("❌ Application creation cancelled")
			return nil
		}

		fmt.Println("🚀 Creating application...")
		// This would use the actual create application API
		fmt.Println("⚠️  Application creation wizard is not fully implemented yet")
		fmt.Println("   Use: coolifyme apps create --repo URL --project UUID --server UUID --environment ENV")

		return nil
	},
}

// Interactive server setup wizard
var serverAddWizardCmd = &cobra.Command{
	Use:   "add-wizard",
	Short: "Interactive server setup wizard",
	Long:  "Guided wizard to add a new server with all necessary configuration",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println("🖥️  Server Setup Wizard")
		fmt.Println("======================")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		// Server name
		fmt.Print("📛 Server name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			return fmt.Errorf("server name is required")
		}

		// Server IP
		fmt.Print("🌐 Server IP address: ")
		ip, _ := reader.ReadString('\n')
		ip = strings.TrimSpace(ip)
		if ip == "" {
			return fmt.Errorf("server IP is required")
		}

		// SSH user
		fmt.Print("👤 SSH user [root]: ")
		user, _ := reader.ReadString('\n')
		user = strings.TrimSpace(user)
		if user == "" {
			user = "root"
		}

		// SSH port
		fmt.Print("🔌 SSH port [22]: ")
		portStr, _ := reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)
		port := 22
		if portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			}
		}

		// Private key UUID
		fmt.Print("🔑 Private key UUID: ")
		privateKey, _ := reader.ReadString('\n')
		privateKey = strings.TrimSpace(privateKey)
		if privateKey == "" {
			return fmt.Errorf("private key UUID is required")
		}

		// Proxy type
		fmt.Print("🔧 Proxy type (traefik/caddy/none) [traefik]: ")
		proxy, _ := reader.ReadString('\n')
		proxy = strings.TrimSpace(proxy)
		if proxy == "" {
			proxy = "traefik"
		}

		// Build server
		fmt.Print("🏗️  Is build server? (y/N): ")
		buildServerStr, _ := reader.ReadString('\n')
		buildServerStr = strings.TrimSpace(strings.ToLower(buildServerStr))
		buildServer := buildServerStr == "y" || buildServerStr == "yes"

		// Description
		fmt.Print("📝 Description (optional): ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)

		fmt.Println("\n📋 Server Configuration Summary:")
		fmt.Printf("   📛 Name: %s\n", name)
		fmt.Printf("   🌐 IP: %s:%d\n", ip, port)
		fmt.Printf("   👤 User: %s\n", user)
		fmt.Printf("   🔑 Private Key: %s\n", privateKey)
		fmt.Printf("   🔧 Proxy: %s\n", proxy)
		fmt.Printf("   🏗️  Build Server: %t\n", buildServer)
		if description != "" {
			fmt.Printf("   📝 Description: %s\n", description)
		}
		fmt.Println()

		fmt.Print("✅ Add server? (y/N): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))

		if confirm != "y" && confirm != "yes" {
			fmt.Println("❌ Server setup cancelled")
			return nil
		}

		fmt.Println("🚀 Adding server...")
		// This would use the actual create server API
		fmt.Printf("coolifyme servers create --name \"%s\" --ip \"%s\" --user \"%s\" --port %d --private-key-uuid \"%s\" --proxy-type \"%s\"",
			name, ip, user, port, privateKey, proxy)
		if buildServer {
			fmt.Print(" --is-build-server")
		}
		if description != "" {
			fmt.Printf(" --description \"%s\"", description)
		}
		fmt.Println()
		fmt.Println("⚠️  Server setup wizard is not fully implemented yet")
		fmt.Println("   Use the command above to create the server")

		return nil
	},
}
