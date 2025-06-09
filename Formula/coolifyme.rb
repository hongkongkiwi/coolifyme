class Coolifyme < Formula
  desc "A powerful command-line interface for the Coolify API"
  homepage "https://github.com/hongkongkiwi/coolifyme"
  url "https://github.com/hongkongkiwi/coolifyme/releases/download/v1.0.0/coolifyme-darwin-amd64.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"
  version "1.0.0"

  depends_on "go" => :build

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/hongkongkiwi/coolifyme/releases/download/v1.0.0/coolifyme-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_AMD64"
    elsif Hardware::CPU.arm?
      url "https://github.com/hongkongkiwi/coolifyme/releases/download/v1.0.0/coolifyme-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_ARM64"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/hongkongkiwi/coolifyme/releases/download/v1.0.0/coolifyme-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_AMD64"
    elsif Hardware::CPU.arm?
      url "https://github.com/hongkongkiwi/coolifyme/releases/download/v1.0.0/coolifyme-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_ARM64"
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