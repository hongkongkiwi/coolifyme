# Using coolifyme in GitHub Actions

This document provides examples of how to use the coolifyme CLI in GitHub Actions workflows.

## Basic Setup

### Method 1: Using the GitHub Action (Recommended)

```yaml
name: Deploy to Coolify

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Coolify CLI
        uses: hongkongkiwi/coolifyme@v1
        with:
          version: latest
      
      - name: Deploy application
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          coolifyme deploy application ${{ vars.APPLICATION_UUID }}
```

### Method 2: Using Installation Script

```yaml
name: Deploy to Coolify

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Install coolifyme
        run: |
          curl -sSL https://raw.githubusercontent.com/hongkongkiwi/coolifyme/main/scripts/install.sh | bash
          echo "/usr/local/bin" >> $GITHUB_PATH
      
      - name: Deploy application
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          coolifyme deploy application ${{ vars.APPLICATION_UUID }}
```

### Method 3: Using Docker

```yaml
name: Deploy to Coolify

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    container:
      image: ghcr.io/hongkongkiwi/coolifyme:latest
    steps:
      - name: Deploy application
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          coolifyme deploy application ${{ vars.APPLICATION_UUID }}
```

## Advanced Examples

### Conditional Deployment

```yaml
name: Conditional Deploy

on:
  push:
    branches: [main, develop]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Coolify CLI
        uses: hongkongkiwi/coolifyme@v1
      
      - name: Deploy to staging
        if: github.ref == 'refs/heads/develop'
        env:
          COOLIFY_API_TOKEN: ${{ secrets.STAGING_COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.STAGING_COOLIFY_BASE_URL }}
        run: |
          coolifyme deploy application ${{ vars.STAGING_APPLICATION_UUID }}
      
      - name: Deploy to production
        if: github.ref == 'refs/heads/main'
        env:
          COOLIFY_API_TOKEN: ${{ secrets.PROD_COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.PROD_COOLIFY_BASE_URL }}
        run: |
          coolifyme deploy application ${{ vars.PROD_APPLICATION_UUID }}
```

### Multi-Application Deployment

```yaml
name: Multi-App Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        app:
          - name: frontend
            uuid: ${{ vars.FRONTEND_UUID }}
          - name: backend
            uuid: ${{ vars.BACKEND_UUID }}
          - name: worker
            uuid: ${{ vars.WORKER_UUID }}
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Coolify CLI
        uses: hongkongkiwi/coolifyme@v1
      
      - name: Deploy ${{ matrix.app.name }}
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          echo "Deploying ${{ matrix.app.name }}..."
          coolifyme deploy application ${{ matrix.app.uuid }}
```

### With Monitoring and Notifications

```yaml
name: Deploy with Monitoring

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Coolify CLI
        uses: hongkongkiwi/coolifyme@v1
      
      - name: Pre-deployment check
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          # Check application status before deployment
          coolifyme applications get ${{ vars.APPLICATION_UUID }}
      
      - name: Deploy application
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          coolifyme deploy application ${{ vars.APPLICATION_UUID }}
      
      - name: Monitor deployment
        env:
          COOLIFY_API_TOKEN: ${{ secrets.COOLIFY_API_TOKEN }}
          COOLIFY_BASE_URL: ${{ secrets.COOLIFY_BASE_URL }}
        run: |
          # Watch deployment status for 5 minutes
          timeout 300 coolifyme watch deployments ${{ vars.APPLICATION_UUID }} --interval 10s || true
      
      - name: Notify on success
        if: success()
        uses: 8398a7/action-slack@v3
        with:
          status: success
          text: "🚀 Application deployed successfully!"
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
      
      - name: Notify on failure
        if: failure()
        uses: 8398a7/action-slack@v3
        with:
          status: failure
          text: "❌ Application deployment failed!"
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

### Custom Configuration

```yaml
name: Deploy with Custom Config

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Coolify CLI
        uses: hongkongkiwi/coolifyme@v1
        with:
          version: v1.0.0  # Pin to specific version
      
      - name: Configure coolifyme
        run: |
          # Create custom configuration
          mkdir -p ~/.config/coolifyme
cat > ~/.config/coolifyme/config.yaml << EOF
          default:
            api_token: "${{ secrets.COOLIFY_API_TOKEN }}"
            base_url: "${{ secrets.COOLIFY_BASE_URL }}"
            timeout: "60s"
            retry_attempts: 3
          EOF
      
      - name: Deploy with force flag
        run: |
          coolifyme deploy application ${{ vars.APPLICATION_UUID }} --force
```

## Environment Variables

When using coolifyme in GitHub Actions, you can set these environment variables:

- `COOLIFY_API_TOKEN`: Your Coolify API token
- `COOLIFY_BASE_URL`: Your Coolify instance URL
- `COOLIFY_PROFILE`: Configuration profile to use (default: "default")
- `COOLIFY_TIMEOUT`: Request timeout (default: "30s")
- `COOLIFY_DEBUG`: Enable debug logging (set to "true")

## Secrets Configuration

Add these secrets to your GitHub repository:

- `COOLIFY_API_TOKEN`: Your Coolify API token
- `COOLIFY_BASE_URL`: Your Coolify instance URL (e.g., https://app.coolify.io/api/v1)

For staging/production environments, you might want separate secrets:
- `STAGING_COOLIFY_API_TOKEN` / `PROD_COOLIFY_API_TOKEN`
- `STAGING_COOLIFY_BASE_URL` / `PROD_COOLIFY_BASE_URL`

## Variables Configuration

Add these variables to your GitHub repository:

- `APPLICATION_UUID`: UUID of your application in Coolify
- `FRONTEND_UUID`, `BACKEND_UUID`, etc.: UUIDs for different applications

## Tips

1. **Pin versions**: Use specific versions instead of `latest` for production workflows
2. **Use matrix builds**: Deploy multiple applications in parallel
3. **Monitor deployments**: Use the watch command to monitor deployment progress
4. **Handle failures**: Add proper error handling and notifications
5. **Cache dependencies**: The GitHub Action automatically handles caching
6. **Security**: Never expose API tokens in logs or outputs 