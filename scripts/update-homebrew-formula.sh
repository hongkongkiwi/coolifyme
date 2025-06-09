#!/bin/bash

# Script to update Homebrew formula with new release
# Usage: ./scripts/update-homebrew-formula.sh v1.0.0

set -e

VERSION=${1:-$(git describe --tags --abbrev=0)}
VERSION=${VERSION#v}  # Remove 'v' prefix

if [ -z "$VERSION" ]; then
    echo "Error: No version specified and no git tags found"
    echo "Usage: $0 [version]"
    exit 1
fi

echo "Updating Homebrew formula for version $VERSION"

# GitHub release URLs
BASE_URL="https://github.com/hongkongkiwi/coolifyme/releases/download/v${VERSION}"
DARWIN_AMD64_URL="${BASE_URL}/coolifyme-darwin-amd64.tar.gz"
DARWIN_ARM64_URL="${BASE_URL}/coolifyme-darwin-arm64.tar.gz"
LINUX_AMD64_URL="${BASE_URL}/coolifyme-linux-amd64.tar.gz"
LINUX_ARM64_URL="${BASE_URL}/coolifyme-linux-arm64.tar.gz"

# Function to get SHA256 from URL
get_sha256() {
    local url=$1
    echo "Downloading $url to calculate SHA256..."
    curl -sL "$url" | shasum -a 256 | cut -d' ' -f1
}

echo "Calculating SHA256 checksums..."
DARWIN_AMD64_SHA256=$(get_sha256 "$DARWIN_AMD64_URL")
DARWIN_ARM64_SHA256=$(get_sha256 "$DARWIN_ARM64_URL")
LINUX_AMD64_SHA256=$(get_sha256 "$LINUX_AMD64_URL")
LINUX_ARM64_SHA256=$(get_sha256 "$LINUX_ARM64_URL")

echo "SHA256 checksums:"
echo "  Darwin AMD64: $DARWIN_AMD64_SHA256"
echo "  Darwin ARM64: $DARWIN_ARM64_SHA256"
echo "  Linux AMD64:  $LINUX_AMD64_SHA256"
echo "  Linux ARM64:  $LINUX_ARM64_SHA256"

# Update the formula file
cat > Formula/coolifyme.rb << EOF
class Coolifyme < Formula
  desc "A powerful command-line interface for the Coolify API"
  homepage "https://github.com/hongkongkiwi/coolifyme"
  license "MIT"
  version "$VERSION"

  on_macos do
    if Hardware::CPU.intel?
      url "$DARWIN_AMD64_URL"
      sha256 "$DARWIN_AMD64_SHA256"
    elsif Hardware::CPU.arm?
      url "$DARWIN_ARM64_URL"
      sha256 "$DARWIN_ARM64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "$LINUX_AMD64_URL"
      sha256 "$LINUX_AMD64_SHA256"
    elsif Hardware::CPU.arm?
      url "$LINUX_ARM64_URL"
      sha256 "$LINUX_ARM64_SHA256"
    end
  end

  def install
    bin.install "coolifyme-darwin-amd64" => "coolifyme" if Hardware::CPU.intel? && OS.mac?
    bin.install "coolifyme-darwin-arm64" => "coolifyme" if Hardware::CPU.arm? && OS.mac?
    bin.install "coolifyme-linux-amd64" => "coolifyme" if Hardware::CPU.intel? && OS.linux?
    bin.install "coolifyme-linux-arm64" => "coolifyme" if Hardware::CPU.arm? && OS.linux?

    # Generate shell completions
    generate_completions_from_executable(bin/"coolifyme", "completion")
  end

  test do
    assert_match "coolifyme version", shell_output("#{bin}/coolifyme --version")
    
    # Test help command
    assert_match "A powerful CLI for the Coolify API", shell_output("#{bin}/coolifyme --help")
    
    # Test that config command exists
    assert_match "Manage configuration", shell_output("#{bin}/coolifyme config --help")
  end
end
EOF

echo "Formula updated successfully!"
echo ""
echo "Next steps:"
echo "1. Create a new repository named 'homebrew-coolifyme'"
echo "2. Copy the Formula/coolifyme.rb file to that repository"
echo "3. Commit and push the changes"
echo "4. Users can then install with: brew install hongkongkiwi/coolifyme/coolifyme" 