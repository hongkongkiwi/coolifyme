# Contributing to coolifyme

Thank you for your interest in contributing to coolifyme! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/yourusername/coolifyme.git
   cd coolifyme
   ```

2. **Install development tools**
   ```bash
   task install-tools
   ```

3. **Set up your development environment**
   ```bash
   # Update API spec and generate client
   task update-and-rebuild
   
   # Run tests
   task test
   
   # Check code quality
   task lint
   ```

## 🔧 Development Workflow

### Code Generation
The API client is auto-generated from the Coolify OpenAPI specification:

```bash
# Update OpenAPI spec and regenerate client
task update-spec
task generate

# Or do both in one command
task update-and-rebuild
```

### Testing
```bash
# Run all tests
task test

# Run tests with coverage
task test-coverage

# Run security scan
task security-scan
```

### Code Quality
```bash
# Format code
task fmt

# Run linter
task lint

# Check all quality metrics
task fmt lint test security-scan
```

## 📋 Coding Standards

- **Go Style**: Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **Testing**: Write tests for new functionality
- **Documentation**: Update documentation for user-facing changes
- **Commits**: Use conventional commit messages

### Conventional Commits
```
feat: add new watch command for monitoring deployments
fix: resolve authentication error handling
docs: update README with new installation methods
refactor: improve error handling in client package
test: add unit tests for config package
```

## 🏗️ Project Structure

```
coolifyme/
├── cmd/                    # CLI commands
├── internal/
│   ├── api/               # Generated API client
│   ├── config/            # Configuration management
│   ├── logger/            # Logging utilities
│   └── output/            # Output formatting
├── pkg/
│   └── client/            # High-level client wrapper
├── spec/                  # OpenAPI specifications
├── examples/              # Example configurations
└── completions/           # Shell completions
```

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment information**
   - Operating system
   - Go version
   - coolifyme version

2. **Steps to reproduce**
3. **Expected behavior**
4. **Actual behavior**
5. **Relevant logs or error messages**

## 🎯 Feature Requests

For feature requests, please:

1. Check existing issues to avoid duplicates
2. Describe the use case and problem you're trying to solve
3. Provide examples of how the feature would be used
4. Consider backward compatibility

## 🔄 Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write code following our standards
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes**
   ```bash
   task fmt lint test security-scan
   ```

4. **Commit your changes**
   ```bash
   git commit -m "feat: add your feature description"
   ```

5. **Push and create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **PR Review**
   - All checks must pass
   - Code review approval required
   - Documentation updates if needed

## 🏷️ Release Process

Releases are automated through GitHub Actions:

1. **Create and push a tag**
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

2. **GitHub Actions will**
   - Run all tests and quality checks
   - Build binaries for multiple platforms
   - Create a GitHub release
   - Generate release notes

## 💡 Development Tips

### Adding New Commands
1. Create a new file in `cmd/` (e.g., `cmd/newcommand.go`)
2. Follow the pattern of existing commands
3. Add the command to `main.go`
4. Write tests for the new functionality
5. Update documentation

### Working with the API Client
The API client is generated from the OpenAPI spec. To add support for new endpoints:

1. Check if the endpoint exists in the latest OpenAPI spec
2. If not, wait for it to be added to Coolify's spec
3. Run `task update-and-rebuild` to regenerate the client
4. Add high-level wrappers in `pkg/client/`

### Debugging
Enable debug logging:
```bash
coolifyme --debug your-command
```

## 📞 Getting Help

- **Issues**: [GitHub Issues](https://github.com/hongkongkiwi/coolifyme/issues)
- **Discussions**: [GitHub Discussions](https://github.com/hongkongkiwi/coolifyme/discussions)
- **Documentation**: Check the README and examples

## 👨‍💻 Author

**Andy Savage** <andy@savage.hk>
- GitHub: [@hongkongkiwi](https://github.com/hongkongkiwi)
- Website: [savage.hk](https://savage.hk)

## 📄 License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).