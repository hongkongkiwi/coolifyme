# coolifyme 🚀

A powerful and feature-rich command-line interface for the [Coolify](https://coolify.io) API.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/hongkongkiwi/coolifyme)](https://github.com/hongkongkiwi/coolifyme/releases)

## Features

### Core Features
✨ **Complete API Coverage**: Support for all 75 Coolify API endpoints (100% coverage)  
🔧 **Profile Management**: Multiple configuration profiles like AWS CLI  
🔍 **Debug Logging**: Detailed API request/response logging with `--debug` flag  
🎨 **Multiple Output Formats**: JSON, YAML, and table formats  
🚀 **Shell Completion**: Bash, Zsh, Fish, and PowerShell completion support  
⚙️ **Environment Variables**: Flexible configuration via environment variables  
📝 **Rich CLI**: Industry-standard CLI patterns with verbose/quiet modes  
🔐 **Secure**: API tokens are handled securely and masked in logs  
📦 **Easy Installation**: Single binary with no dependencies

### Industry-Standard CLI Features
🔍 **Search & Filtering**: Universal search across resources with wildcard patterns and advanced filtering  
⏱️ **Timeouts & Retry Logic**: Global timeout configuration with exponential backoff retry logic  
🎨 **Advanced Output Formatting**: Multiple formats (JSON, YAML, CSV) with custom templates and sorting  
🔄 **Rollback Operations**: Safe rollbacks with deployment history and dry-run support  
🎯 **Interactive Wizards**: Guided setup for first-time configuration and complex operations  
⚡ **Bulk Operations**: Mass operations with concurrency control and dry-run support  
📊 **Monitoring & Health Checks**: Real-time status monitoring and system health verification  
🔗 **Command Aliases**: Quick shortcuts for frequently used commands  
🚀 **Auto-Updates**: Smart update detection with Homebrew integration

## Table of Contents

- [coolifyme 🚀](#coolifyme-)
  - [Features](#features)
    - [Core Features](#core-features)
    - [Industry-Standard CLI Features](#industry-standard-cli-features)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Binary Releases](#binary-releases)
    - [Build from Source](#build-from-source)
      - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
  - [Configuration](#configuration)
    - [Profiles](#profiles)
    - [Environment Variables](#environment-variables)
  - [Usage](#usage)
    - [Global Options](#global-options)
    - [Applications](#applications)
    - [Deployments](#deployments)
    - [Servers](#servers)
    - [Services](#services)
    - [Databases](#databases)
  - [Industry-Standard CLI Features](#industry-standard-cli-features-1)
    - [Search \& Filtering System 🔍](#search--filtering-system-)
    - [Global Timeouts \& Retry Logic ⏱️](#global-timeouts--retry-logic-️)
    - [Advanced Output Formatting 🎨](#advanced-output-formatting-)
    - [Rollback Operations 🔄](#rollback-operations-)
    - [Interactive Wizards 🧙](#interactive-wizards-)
    - [Bulk Operations 📦](#bulk-operations-)
    - [Monitoring \& Health Checks 📊](#monitoring--health-checks-)
    - [Command Aliases 🚀](#command-aliases-)
    - [Auto-Updates 🔄](#auto-updates-)
    - [Environment Variables Management](#environment-variables-management)
  - [Shell Completion](#shell-completion)
    - [Bash](#bash)
    - [Zsh](#zsh)
    - [Fish](#fish)
    - [PowerShell](#powershell)
  - [Debug and Logging](#debug-and-logging)
    - [Debug Mode](#debug-mode)
    - [Logging Levels](#logging-levels)
    - [Sample Debug Output](#sample-debug-output)
  - [Output Formats](#output-formats)
  - [Development](#development)
    - [Project Structure](#project-structure)
    - [Building](#building)
    - [API Coverage](#api-coverage)
  - [Contributing](#contributing)
  - [License](#license)
  - [Support](#support)
  - [Acknowledgments](#acknowledgments)

## Installation

### Binary Releases

Download the latest binary from the [releases page](https://github.com/hongkongkiwi/coolifyme/releases):

```bash
# macOS (Intel)
curl -L https://github.com/hongkongkiwi/coolifyme/releases/latest/download/coolifyme-darwin-amd64 -o coolifyme
chmod +x coolifyme
sudo mv coolifyme /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/hongkongkiwi/coolifyme/releases/latest/download/coolifyme-darwin-arm64 -o coolifyme
chmod +x coolifyme
sudo mv coolifyme /usr/local/bin/

# Linux
curl -L https://github.com/hongkongkiwi/coolifyme/releases/latest/download/coolifyme-linux-amd64 -o coolifyme
chmod +x coolifyme
sudo mv coolifyme /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/hongkongkiwi/coolifyme.git
cd coolifyme
task install
```

#### Prerequisites

- Go 1.21 or later
- [Task](https://taskfile.dev/) (for build automation)

**Install Task:**

```bash
# macOS (using Homebrew)
brew install go-task/tap/go-task

# Ubuntu/Debian
sudo snap install task --classic

# Or via Go
go install github.com/go-task/task/v3/cmd/task@latest
```

## Quick Start

1. **Initialize configuration:**
   ```bash
   coolifyme config init
   ```

2. **Set your API token:**
   ```bash
   coolifyme config profile set --token YOUR_API_TOKEN
   ```

3. **Or create a new profile for your environment:**
   ```bash
   coolifyme config profile create production \
     --token YOUR_API_TOKEN \
     --url https://your-coolify-instance.com/api/v1
   ```

4. **List your applications:**
   ```bash
   coolifyme applications list
   ```

5. **Deploy an application:**
   ```bash
   coolifyme deploy application app-uuid-here
   ```

## Configuration

### Profiles

coolifyme supports multiple configuration profiles, similar to AWS CLI, allowing you to manage different Coolify instances or environments:

```bash
# Create profiles for different environments
coolifyme config profile create production --token TOKEN1 --url https://coolify.prod.com/api/v1
coolifyme config profile create staging --token TOKEN2 --url https://coolify.staging.com/api/v1
coolifyme config profile create local --token TOKEN3 --url http://localhost:8000/api/v1

# List all profiles
coolifyme config profile list

# Switch between profiles
coolifyme config profile use production

# Use a profile for a single command
coolifyme --profile staging applications list

# Update current profile
coolifyme config profile set --token NEW_TOKEN
```

Configuration is stored in `~/.config/coolifyme/config.yaml`:

```yaml
default_profile: production
profiles:
  production:
    name: production
    api_token: your_production_token
    base_url: https://coolify.yourdomain.com/api/v1
  staging:
    name: staging
    api_token: your_staging_token
    base_url: https://staging.coolify.yourdomain.com/api/v1
global_settings:
  output_format: table
  log_level: info
  color_output: true
```

### Environment Variables

Configure coolifyme using environment variables:

```bash
# API configuration
export COOLIFY_API_TOKEN="your_api_token"
export COOLIFY_BASE_URL="https://your-coolify-instance.com/api/v1"
export COOLIFY_PROFILE="production"

# Output and logging
export COOLIFY_LOG_LEVEL="debug"
export COOLIFY_OUTPUT_FORMAT="json"

# Backward compatibility
export COOLIFYME_API_TOKEN="your_api_token"  # Also supported
export COOLIFYME_BASE_URL="your_base_url"    # Also supported
```

## Usage

### Global Options

All commands support these global options:

```bash
  --color string     colorize output (auto, always, never) (default "auto")
  --config string    config file (default is ~/.config/coolifyme/config.yaml)
  --debug            debug output (shows API calls)
  -o, --output string    output format (json, yaml, table)
  -p, --profile string   configuration profile to use
  -q, --quiet            quiet output (errors only)
  -s, --server string    Coolify server URL
  -t, --token string     API token
  -v, --verbose          verbose output
```

### Applications

```bash
# List all applications
coolifyme applications list
coolifyme apps ls

# Get application details
coolifyme apps get <uuid>

# Start/stop/restart applications
coolifyme apps start <uuid>
coolifyme apps stop <uuid>
coolifyme apps restart <uuid>

# View application logs
coolifyme apps logs <uuid> --lines 100

# Manage environment variables
coolifyme apps env list <uuid>
coolifyme apps env export <uuid> --file .env
coolifyme apps env import <uuid> --file .env
coolifyme apps env sync <uuid> --file .env
coolifyme apps env cleanup <uuid> --file .env  # Remove non-existent vars

# Environment variable operations with preview
coolifyme apps env import <uuid> --file .env --dry-run
coolifyme apps env sync <uuid> --file .env --dry-run
```

### Deployments

```bash
# Deploy an application
coolifyme deploy application <uuid>
coolifyme deploy app <uuid> --force

# Deploy from specific branch or PR
coolifyme deploy app <uuid> --branch main
coolifyme deploy app <uuid> --pr 123

# Deploy multiple applications
coolifyme deploy multiple <uuid1> <uuid2> <uuid3>

# Monitor deployment
coolifyme deploy watch <deployment-uuid>
coolifyme deploy logs <deployment-uuid>

# List deployments
coolifyme deployments list
coolifyme deployments list-by-app <app-uuid>
```

### Servers

```bash
# List all servers
coolifyme servers list
coolifyme srv ls

# Get server details
coolifyme srv get <uuid>

# Create a new server
coolifyme srv create \
  --name "production-server" \
  --ip "192.168.1.100" \
  --user "root" \
  --private-key-uuid "key-uuid" \
  --proxy-type "traefik"

# Update server configuration
coolifyme srv update <uuid> --name "new-name"

# Validate server connection
coolifyme srv validate <uuid>

# Get server resources and domains
coolifyme srv get-resources <uuid>
coolifyme srv get-domains <uuid>

# Delete a server
coolifyme srv delete <uuid> --force
```

### Services

```bash
# List all services
coolifyme services list
coolifyme svc ls

# Get service details
coolifyme svc get <uuid>

# Create a service
coolifyme svc create \
  --type "docker-compose" \
  --name "my-service" \
  --project "project-uuid" \
  --server "server-uuid" \
  --environment "production"

# Update and delete services
coolifyme svc update <uuid> --name "new-name"
coolifyme svc delete <uuid> --force

# Manage service environment variables
coolifyme svc env list <uuid>
coolifyme svc env create <uuid> --key "DATABASE_URL" --value "postgres://..."
coolifyme svc env update <uuid> <env-uuid> --value "new-value"
coolifyme svc env delete <uuid> <env-uuid>

# Bulk operations
coolifyme svc start-all
coolifyme svc stop-all
coolifyme svc restart-all
```

### Databases

```bash
# List all databases
coolifyme databases list
coolifyme db ls

# Create databases
coolifyme db create postgresql --project "uuid" --server "uuid" --environment "prod"
coolifyme db create mysql --project "uuid" --server "uuid" --environment "prod"
coolifyme db create redis --project "uuid" --server "uuid" --environment "prod"
coolifyme db create mongodb --project "uuid" --server "uuid" --environment "prod"

# Specialized databases
coolifyme db create clickhouse --project "uuid" --server "uuid" --environment "prod"
coolifyme db create dragonfly --project "uuid" --server "uuid" --environment "prod"
coolifyme db create keydb --project "uuid" --server "uuid" --environment "prod"
coolifyme db create mariadb --project "uuid" --server "uuid" --environment "prod"

# Get database details
coolifyme db get <uuid>

# Delete a database
coolifyme db delete <uuid> --force
```

## Industry-Standard CLI Features

### Search & Filtering System 🔍

Universal search capabilities across all Coolify resources:

```bash
# Search across all resources
coolifyme search "my-app"
coolifyme search "prod-*" --type applications
coolifyme search "nginx" --status running

# Advanced filtering
coolifyme find --name "api-*" --status running --type services
coolifyme find --status failed --type applications

# Output control
coolifyme search "web" --json --limit 10
coolifyme search "database" --case-sensitive
```

**Features:**
- Cross-resource search: applications, services, servers, databases
- Wildcard pattern support (`app-*`, `prod-*`)
- Case-sensitive and case-insensitive modes
- Status and tag filtering with result limiting
- JSON output for scripting and automation

### Global Timeouts & Retry Logic ⏱️

Robust timeout and retry configuration for API operations:

```bash
# Configure global timeouts
coolifyme timeout set --timeout 60s --retry 5 --retry-delay 2s
coolifyme timeout show

# Per-command overrides
coolifyme applications list --timeout 30s --retry 3
coolifyme deploy app uuid --timeout 300s --retry 1
```

**Features:**
- Request timeout configuration (default: 30s)
- Retry count with exponential backoff (default: 3 retries)
- Customizable retry delays and maximum backoff
- Per-command timeout overrides
- Smart retry logic with jitter to prevent thundering herd

### Advanced Output Formatting 🎨

Powerful output formatting for different use cases:

```bash
# Multiple output formats
coolifyme apps list --output json
coolifyme apps list --output yaml
coolifyme apps list --output csv
coolifyme apps list --output wide

# Custom table columns
coolifyme apps list --output "table(name,status,url)"
coolifyme apps list --columns name,status,url

# Template formatting
coolifyme apps list --output "custom({name} is {status})"

# Sorting and filtering
coolifyme apps list --sort-by name --sort-reverse
coolifyme apps list --no-headers --show-kind

# Examples and help
coolifyme format examples
```

**Supported Formats:**
- **JSON/YAML**: Machine-readable for automation
- **Table**: Human-readable with column control
- **CSV**: Spreadsheet integration
- **Custom**: Template-based output formatting
- **Wide**: Extended information display
- **Name-only**: Just resource names

### Rollback Operations 🔄

Safe rollback capabilities with history tracking:

```bash
# Application rollbacks
coolifyme rollback app <uuid> --list                    # Show available versions
coolifyme rollback app <uuid> --to-commit abc123        # Rollback to git commit
coolifyme rollback app <uuid> --to-version v1.2.3       # Rollback to version
coolifyme rollback app <uuid> --to-commit abc123 --dry-run  # Preview rollback

# Deployment history
coolifyme rollback history <uuid> --limit 20
coolifyme rollback history <uuid> --json

# Service rollbacks
coolifyme rollback service <uuid> --to-version v1.1.0
```

**Features:**
- Git commit-based and version-based rollbacks
- Deployment history tracking and visualization
- Dry-run mode for safe previewing
- Force rollback with `--force` flag
- Confirmation prompts for safety

### Interactive Wizards 🧙

Guided setup and configuration wizards:

```bash
# First-time setup wizard
coolifyme init-interactive

# Interactive application creation
coolifyme applications create-wizard

# Interactive server setup
coolifyme servers add-wizard
```

These wizards guide you through complex operations with prompts, validation, and helpful descriptions.

### Bulk Operations 📦

Efficiently manage multiple resources with built-in concurrency control:

```bash
# Start all applications (with dry-run support)
coolifyme applications start-all --dry-run
coolifyme applications start-all --concurrent 10

# Stop all applications
coolifyme applications stop-all --concurrent 5

# Restart all applications
coolifyme applications restart-all

# Deploy all services
coolifyme services deploy-all --dry-run --concurrent 3
```

**Features:**
- `--dry-run`: Preview what would be executed without making changes
- `--concurrent N`: Control parallelism (default: 5)
- Progress tracking and detailed result summaries
- Error handling for individual operations

### Monitoring & Health Checks 📊

Comprehensive monitoring tools for your Coolify infrastructure:

```bash
# Quick health check
coolifyme health
coolifyme monitor health --verbose

# Status overview
coolifyme status
coolifyme monitor status

# Real-time monitoring (auto-refresh)
coolifyme monitor watch --interval 30
```

**Health Check Features:**
- API connectivity verification
- Resource counting and validation
- Timeout handling for reliability
- Verbose mode for detailed diagnostics

### Command Aliases 🚀

Quick shortcuts for frequently used commands:

```bash
# Deployment aliases
coolifyme deploy-app <uuid>    # Short for: deploy application
coolifyme deploy <uuid>        # Even shorter
coolifyme dep <uuid>           # Shortest

# Status aliases
coolifyme status               # Quick status overview
coolifyme st                   # Short form
coolifyme ping                 # Health check

# List aliases
coolifyme ls-apps              # List applications
coolifyme ls-servers           # List servers  
coolifyme ls-services          # List services

# View all aliases
coolifyme alias list
```

**Available Aliases:**
- **Deployment**: `deploy-app`, `deploy`, `dep` → `deploy application`
- **Monitoring**: `status`, `st`, `stat` → `monitor status`
- **Health**: `health`, `ping`, `check` → `monitor health`
- **Listing**: `ls-apps`, `ls-servers`, `ls-services` → respective list commands

### Auto-Updates 🔄

Smart update management with Homebrew integration:

```bash
# Check for updates and install
coolifyme update

# Force update check
coolifyme update --force

# Just check version without updating
coolifyme version
```

**Update Features:**
- Auto-detects Homebrew installation
- Runs `brew upgrade coolifyme` if installed via Homebrew
- Falls back to manual instructions for other installations
- Version information with build details

### Environment Variables Management

coolifyme provides powerful `.env` file management capabilities:

```bash
# Export application env vars to .env file
coolifyme apps env export <app-uuid> --file .env

# Import env vars from .env file to application
coolifyme apps env import <app-uuid> --file .env --dry-run
coolifyme apps env import <app-uuid> --file .env

# Bidirectional sync between .env file and application
coolifyme apps env sync <app-uuid> --file .env --dry-run
coolifyme apps env sync <app-uuid> --file .env

# Clean up .env file (remove variables that don't exist in app)
coolifyme apps env cleanup <app-uuid> --file .env --backup
```

**Features:**
- 🔍 **Dry-run mode**: Preview changes before applying
- 📄 **Automatic backups**: Create backups before modifying files
- 🔄 **Bidirectional sync**: Keep .env files and applications in sync
- 🧹 **Cleanup**: Remove stale variables from .env files
- 📝 **Multiline support**: Handle complex environment variables

## Shell Completion

Enable shell completion for better CLI experience:

### Bash
```bash
# Add to current session
source <(coolifyme completion bash)

# Add to ~/.bashrc for persistence
echo 'source <(coolifyme completion bash)' >> ~/.bashrc

# Or install system-wide (Linux)
coolifyme completion bash | sudo tee /etc/bash_completion.d/coolifyme

# Or install system-wide (macOS with Homebrew)
coolifyme completion bash > /usr/local/etc/bash_completion.d/coolifyme
```

### Zsh
```bash
# Enable completion in zsh (add to ~/.zshrc if not already present)
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Add completion
coolifyme completion zsh > "${fpath[1]}/_coolifyme"

# Restart shell
exec zsh
```

### Fish
```bash
# Add to current session
coolifyme completion fish | source

# Add permanently
coolifyme completion fish > ~/.config/fish/completions/coolifyme.fish
```

### PowerShell
```powershell
# Add to current session
coolifyme completion powershell | Out-String | Invoke-Expression

# Add to profile for persistence
coolifyme completion powershell > coolifyme.ps1
# Then source it from your PowerShell profile
```

## Debug and Logging

coolifyme provides comprehensive logging and debugging capabilities:

### Debug Mode
```bash
# Show all API requests and responses
coolifyme --debug applications list

# See detailed HTTP calls with timing
coolifyme --debug deploy application app-uuid
```

### Logging Levels
```bash
# Set log level globally
coolifyme config set --log-level debug

# Or use flags for specific commands
coolifyme --verbose servers list   # Info level + verbose output
coolifyme --quiet deploy app uuid  # Errors only
```

### Sample Debug Output
```
2024-01-15 10:30:45 DEBUG API Request method=GET url=https://app.coolify.io/api/v1/applications headers="Accept: application/json; Authorization: [REDACTED]; Content-Type: application/json"
2024-01-15 10:30:45 DEBUG API Response method=GET url=https://app.coolify.io/api/v1/applications status="200 OK" duration=245ms headers="Content-Type: application/json; ..."
```

## Output Formats

Support for multiple output formats:

```bash
# Table format (default) - human readable
coolifyme applications list

# JSON format - for scripting
coolifyme --output json applications list

# YAML format - structured and readable
coolifyme --output yaml applications list

# Set default format
coolifyme config set --output json
```

## Development

### Project Structure

```
coolifyme/
├── cmd/                    # CLI commands
│   ├── main.go            # Main entry point with enhanced global flags
│   ├── applications.go    # Applications subcommand with env management
│   ├── config.go          # Configuration with profile management
│   ├── deploy.go          # Deployment commands
│   ├── servers.go         # Server management
│   ├── services.go        # Service management
│   └── databases.go       # Database management
├── internal/
│   ├── api/               # Generated API client (auto-generated)
│   ├── config/            # Configuration management with profiles
│   └── logger/            # Enhanced logging system
├── pkg/
│   └── client/            # High-level API client with debug logging
├── spec/
│   └── coolify-openapi.yaml  # Coolify OpenAPI specification
├── example-config.yaml    # Example configuration file
├── Taskfile.yml           # Build automation
└── oapi-codegen.yaml      # Code generation configuration
```

### Building

```bash
# Generate API client and build
task

# Just generate API client from OpenAPI spec
task generate

# Just build the binary
task build

# Update OpenAPI spec from official source
task update-spec

# Update spec and rebuild everything
task update-and-rebuild

# Clean build artifacts
task clean

# Run tests
task test

# Format code
task fmt

# See all available tasks
task --list
```

### API Coverage

coolifyme provides **100% coverage** of the Coolify API with 75/75 endpoints:

- ✅ **Applications**: 19/19 endpoints (list, get, create, update, delete, logs, env management, start/stop/restart)
- ✅ **Servers**: 8/8 endpoints (CRUD operations, validation, resources, domains)
- ✅ **Teams**: 5/5 endpoints (team management)
- ✅ **Projects**: 6/6 endpoints (project and environment management)
- ✅ **Private Keys**: 5/5 endpoints (SSH key management)
- ✅ **Resources**: 1/1 endpoints (resource information)
- ✅ **API Management**: 3/3 endpoints (version, enable/disable, healthcheck)
- ✅ **Services**: 13/13 endpoints (complete service lifecycle management)
- ✅ **Databases**: 15/15 endpoints (all database types including specialized ones)
- ✅ **Deployments**: Multiple deployment endpoints with enhanced monitoring

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📚 **Documentation**: Check this README and command help (`coolifyme --help`)
- 🐛 **Issues**: [GitHub Issues](https://github.com/hongkongkiwi/coolifyme/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/hongkongkiwi/coolifyme/discussions)
- 🌟 **Feature Requests**: [GitHub Issues](https://github.com/hongkongkiwi/coolifyme/issues/new?template=feature_request.md)

## Acknowledgments

- [Coolify](https://coolify.io) - For creating an amazing self-hosting platform
- [Cobra](https://github.com/spf13/cobra) - For the excellent CLI framework
- [Viper](https://github.com/spf13/viper) - For configuration management
- All contributors who make this project better

---

Built with ❤️ for the Coolify community 