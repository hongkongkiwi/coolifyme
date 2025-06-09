# coolifyme 🚀

A powerful command-line interface for the [Coolify](https://coolify.io) API, built with Go and automatically generated from the official OpenAPI specification.

## Features

- **Auto-generated from OpenAPI spec** - Always up-to-date with the latest Coolify API
- **Easy updates** - Simply run `make update-spec` to get the latest API changes
- **Full API coverage** - Supports all Coolify API endpoints
- **Configuration management** - Store API tokens and settings locally
- **Multiple output formats** - Human-readable tables and JSON output
- **Cross-platform** - Works on macOS, Linux, and Windows

## 🚀 Installation

### Homebrew (macOS/Linux) - Recommended

```bash
# Add the tap
brew tap hongkongkiwi/coolifyme

# Install coolifyme
brew install coolifyme

# Verify installation
coolifyme --version
```

**Benefits of Homebrew installation:**
- ✅ Automatic dependency management
- ✅ Shell completions automatically installed
- ✅ Easy updates with `brew upgrade coolifyme`
- ✅ Uninstall with `brew uninstall coolifyme`

### Quick Install Script

```bash
curl -sSL https://raw.githubusercontent.com/hongkongkiwi/coolifyme/main/scripts/install.sh | bash
```

### Download Pre-built Binary

Download the latest release from [GitHub Releases](https://github.com/hongkongkiwi/coolifyme/releases) for your platform.

### Using in GitHub Actions

```yaml
- name: Setup Coolify CLI
  uses: hongkongkiwi/coolifyme@v1
  with:
    version: latest

- name: Deploy application
  env:
    COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
    COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
  run: coolifyme deploy application ${{ vars.APPLICATION_UUID }}
```

See [examples/github-actions.md](examples/github-actions.md) for more GitHub Actions examples.

### Docker

```bash
docker run --rm -v ~/.config/coolifyme:/home/coolify/.config/coolifyme ghcr.io/hongkongkiwi/coolifyme:latest --help
```

### From Source

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

# macOS (using MacPorts)
sudo port install go-task

# Ubuntu/Debian
sudo snap install task --classic

# Or via apt (Ubuntu 22.04+)
sudo apt update && sudo apt install task

# Arch Linux
sudo pacman -S go-task

# Fedora/CentOS/RHEL
sudo dnf install go-task

# Windows (using Chocolatey)
choco install go-task

# Windows (using Scoop)
scoop install task

# Universal install script
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d ~/.local/bin

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
   coolifyme config set --token YOUR_API_TOKEN
   ```

3. **List your applications:**
   ```bash
   coolifyme applications list
   ```

## Configuration

The CLI stores configuration in `~/.config/coolifyme/config.yaml`. You can manage configuration using:

```bash
# Initialize with defaults
coolifyme config init

# Set API token
coolifyme config set --token YOUR_API_TOKEN

# Set custom base URL (for self-hosted instances)
coolifyme config set --url https://your-coolify-instance.com/api/v1

# View current configuration
coolifyme config show
```

### Environment Variables

You can also configure the CLI using environment variables:

- `COOLIFY_API_TOKEN` - Your Coolify API token
- `COOLIFY_BASE_URL` - Base URL for your Coolify instance
- `COOLIFY_PROFILE` - Configuration profile to use

## Usage

### Applications

```bash
# List all applications
coolifyme applications list
coolifyme apps ls

# Get application details
coolifyme apps get <uuid>

# Create a new application (coming soon)
coolifyme apps create --name myapp --repo https://github.com/user/repo
```

### Projects

```bash
# List all projects
coolifyme projects list
coolifyme proj ls
```

### Servers

```bash
# List all servers
coolifyme servers list
coolifyme srv ls
```

### Global Flags

All commands support these global flags:

- `--token` - Override API token
- `--url` - Override base URL
- `--config` - Specify config file location
- `--profile` - Use specific configuration profile

## Development

### Project Structure

```
coolifyme/
├── cmd/                    # CLI commands
│   ├── main.go            # Main entry point
│   ├── applications.go    # Applications subcommand
│   ├── config.go          # Configuration subcommand
│   ├── projects.go        # Projects subcommand
│   └── servers.go         # Servers subcommand
├── internal/
│   ├── api/               # Generated API client (auto-generated)
│   └── config/            # Configuration management
├── pkg/
│   └── client/            # High-level API client wrapper
├── spec/
│   └── coolify-openapi.yaml  # Coolify OpenAPI specification
├── Makefile               # Build automation
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

### Updating the API Client

When Coolify releases API updates:

```bash
# Update the OpenAPI spec and regenerate client
task update-and-rebuild
```

The generated API client will be automatically updated with any new endpoints, request/response models, and API changes.

### Code Generation

This project uses [oapi-codegen](https://github.com/deepmap/oapi-codegen) to generate Go client code from the Coolify OpenAPI specification. The configuration is in `oapi-codegen.yaml`:

```yaml
generate:
  client: true          # Generate HTTP client
  models: true          # Generate data models
  embedded-spec: true   # Embed OpenAPI spec in generated code
package: coolify
output: internal/api/coolify_client.go
```

## 👨‍💻 Author

**Andy Savage** <andy@savage.hk>
- GitHub: [@hongkongkiwi](https://github.com/hongkongkiwi)
- Website: [savage.hk](https://savage.hk)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🍺 Homebrew Tap Setup (For Maintainers)

This repository includes automatic Homebrew tap management. Here's how it works:

### Automated Process

1. **Tag a release**: `git tag v1.0.0 && git push origin v1.0.0`
2. **GitHub Actions automatically**:
   - Builds cross-platform binaries
   - Creates GitHub release
   - Updates Homebrew tap at `hongkongkiwi/homebrew-coolifyme`
   - Calculates SHA256 checksums for all platforms

### Manual Setup (One-time)

1. **Create the tap repository**:
   ```bash
   # Create a new repository named 'homebrew-coolifyme'
   gh repo create hongkongkiwi/homebrew-coolifyme --public --description "Homebrew tap for coolifyme CLI"
   ```

2. **Set up GitHub token**:
   - Create a Personal Access Token with `repo` permissions
   - Add it as `HOMEBREW_TAP_TOKEN` secret in this repository's settings

3. **Initial formula**:
   ```bash
   # In the homebrew-coolifyme repository
   mkdir Formula
   cp Formula/coolifyme.rb Formula/coolifyme.rb
   git add Formula/coolifyme.rb
   git commit -m "Initial formula"
   git push
   ```

### Manual Updates

If you need to manually update the Homebrew formula:

```bash
# Run the update script
./scripts/update-homebrew-formula.sh v1.0.0

# Or for latest tag
./scripts/update-homebrew-formula.sh
```

### Users Can Install With

```bash
brew tap hongkongkiwi/coolifyme
brew install coolifyme
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## API Documentation

The CLI is built against the official Coolify API. For detailed API documentation, see:

- [Coolify API Documentation](https://app.coolify.io/docs/api)
- [OpenAPI Specification](https://github.com/coollabsio/coolify/blob/next/openapi.yaml)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `task fmt` and `task test`
6. Submit a pull request

### Adding New Commands

1. Create a new command file in `cmd/`
2. Follow the existing patterns for Cobra commands
3. Add the command to `main.go`
4. Update this README

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Coolify](https://coolify.io) - The amazing self-hosted platform
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI code generator
- [Viper](https://github.com/spf13/viper) - Configuration management

---

Built with ❤️ for the Coolify community 