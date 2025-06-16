package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update coolifyme to the latest version",
	Long: `Update coolifyme to the latest version.

If installed via Homebrew, uses 'brew upgrade coolifyme'.
Otherwise, provides instructions for manual update.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if isInstalledViaHomebrew() {
			return updateViaHomebrew(force)
		}

		return showManualUpdateInstructions()
	},
}

// isInstalledViaHomebrew checks if coolifyme was installed via Homebrew
func isInstalledViaHomebrew() bool {
	// Get the path of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	// Check if it's in a Homebrew path
	return strings.Contains(execPath, "/homebrew/") || strings.Contains(execPath, "/opt/homebrew/")
}

// updateViaHomebrew updates coolifyme using Homebrew
func updateViaHomebrew(force bool) error {
	fmt.Println("🍺 Detected Homebrew installation")

	var cmd *exec.Cmd
	if force {
		fmt.Println("🔄 Force updating coolifyme via Homebrew...")
		cmd = exec.Command("brew", "upgrade", "coolifyme", "--force")
	} else {
		fmt.Println("🔄 Updating coolifyme via Homebrew...")
		cmd = exec.Command("brew", "upgrade", "coolifyme")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update via Homebrew: %w", err)
	}

	fmt.Println("✅ coolifyme updated successfully via Homebrew!")
	return nil
}

// showManualUpdateInstructions shows instructions for manual update
func showManualUpdateInstructions() error {
	fmt.Println("📦 Manual Installation Detected")
	fmt.Println("")
	fmt.Println("To update coolifyme manually:")
	fmt.Println("")
	fmt.Println("1. Download the latest release:")
	fmt.Println("   https://github.com/hongkongkiwi/coolifyme/releases/latest")
	fmt.Println("")
	fmt.Println("2. Or use curl to download and install:")
	fmt.Println("   # macOS (Intel)")
	fmt.Println("   curl -L https://github.com/hongkongkiwi/coolifyme/releases/latest/download/coolifyme-darwin-amd64 -o coolifyme")
	fmt.Println("   chmod +x coolifyme && sudo mv coolifyme /usr/local/bin/")
	fmt.Println("")
	fmt.Println("   # macOS (Apple Silicon)")
	fmt.Println("   curl -L https://github.com/hongkongkiwi/coolifyme/releases/latest/download/coolifyme-darwin-arm64 -o coolifyme")
	fmt.Println("   chmod +x coolifyme && sudo mv coolifyme /usr/local/bin/")
	fmt.Println("")
	fmt.Println("   # Linux")
	fmt.Println("   curl -L https://github.com/hongkongkiwi/coolifyme/releases/latest/download/coolifyme-linux-amd64 -o coolifyme")
	fmt.Println("   chmod +x coolifyme && sudo mv coolifyme /usr/local/bin/")
	fmt.Println("")
	fmt.Println("3. Or build from source:")
	fmt.Println("   git clone https://github.com/hongkongkiwi/coolifyme.git")
	fmt.Println("   cd coolifyme && task install")

	return nil
}

func init() {
	updateCmd.Flags().BoolP("force", "f", false, "Force update even if already up to date (Homebrew only)")
}
