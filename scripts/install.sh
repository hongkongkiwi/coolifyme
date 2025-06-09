#!/bin/bash

# coolifyme installation script
# Usage: curl -sSL https://raw.githubusercontent.com/hongkongkiwi/coolifyme/main/scripts/install.sh | bash

set -e

# Default values
REPO="hongkongkiwi/coolifyme"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    local os arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case $os in
        linux)
            os="linux"
            ;;
        darwin)
            os="darwin"
            ;;
        mingw*|msys*|cygwin*)
            os="windows"
            ;;
        *)
            error "Unsupported operating system: $os"
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            error "Unsupported architecture: $arch"
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest version from GitHub API
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        info "Getting latest version from GitHub API..."
        VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            error "Failed to get latest version"
        fi
        info "Latest version: $VERSION"
    fi
    echo "$VERSION"
}

# Download and install binary
install_binary() {
    local platform version binary_name download_url tmp_dir
    
    platform=$(detect_platform)
    version=$(get_latest_version)
    
    # Remove 'v' prefix if present
    version=${version#v}
    
    if [[ $platform == *"windows"* ]]; then
        binary_name="coolifyme-${platform}.exe"
    else
        binary_name="coolifyme-${platform}"
    fi
    
    download_url="https://github.com/$REPO/releases/download/v${version}/${binary_name}.tar.gz"
    
    info "Downloading coolifyme $version for $platform..."
    info "URL: $download_url"
    
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    # Download and extract
    if ! curl -sL "$download_url" | tar -xz -C "$tmp_dir"; then
        error "Failed to download or extract binary"
    fi
    
    # Verify binary exists
    if [ ! -f "$tmp_dir/$binary_name" ]; then
        error "Binary not found in archive"
    fi
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        warn "Creating install directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi
    
    # Install binary
    info "Installing to $INSTALL_DIR/coolifyme..."
    if [[ $platform == *"windows"* ]]; then
        mv "$tmp_dir/$binary_name" "$INSTALL_DIR/coolifyme.exe"
        chmod +x "$INSTALL_DIR/coolifyme.exe"
    else
        mv "$tmp_dir/$binary_name" "$INSTALL_DIR/coolifyme"
        chmod +x "$INSTALL_DIR/coolifyme"
    fi
    
    info "Installation completed successfully!"
    
    # Verify installation
    if command -v coolifyme >/dev/null 2>&1; then
        info "coolifyme is now available in your PATH"
        coolifyme --version
    else
        warn "coolifyme installed to $INSTALL_DIR but not found in PATH"
        warn "You may need to add $INSTALL_DIR to your PATH environment variable"
    fi
}

# Main execution
main() {
    info "Installing coolifyme CLI..."
    
    # Check dependencies
    for cmd in curl tar; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            error "$cmd is required but not installed"
        fi
    done
    
    install_binary
    
    info "To get started, run: coolifyme --help"
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 