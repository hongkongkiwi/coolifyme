package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Global timeout and retry settings
var (
	globalTimeout    time.Duration
	globalRetryCount int
	globalRetryDelay time.Duration
)

// TimeoutConfig holds timeout and retry configuration
type TimeoutConfig struct {
	Timeout    time.Duration
	RetryCount int
	RetryDelay time.Duration
	MaxBackoff time.Duration
}

// DefaultTimeoutConfig returns the default timeout configuration
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: 1 * time.Second,
		MaxBackoff: 10 * time.Second,
	}
}

// WithTimeout wraps a function with timeout and retry logic
func WithTimeout[T any](ctx context.Context, config *TimeoutConfig, operation func(context.Context) (T, error)) (T, error) {
	var result T
	var lastErr error

	for attempt := 0; attempt <= config.RetryCount; attempt++ {
		// Create context with timeout for this attempt
		timeoutCtx, cancel := context.WithTimeout(ctx, config.Timeout)

		// Execute the operation
		result, err := operation(timeoutCtx)
		cancel()

		if err == nil {
			return result, nil
		}

		lastErr = err

		// Don't retry on last attempt
		if attempt == config.RetryCount {
			break
		}

		// Calculate backoff delay
		delay := config.RetryDelay * time.Duration(1<<attempt) // Exponential backoff
		if delay > config.MaxBackoff {
			delay = config.MaxBackoff
		}

		fmt.Printf("⚠️  Attempt %d failed: %v. Retrying in %v...\n", attempt+1, err, delay)

		// Wait before retrying
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return result, fmt.Errorf("operation failed after %d attempts: %w", config.RetryCount+1, lastErr)
}

// timeoutCmd represents the timeout command for configuring global timeouts
var timeoutCmd = &cobra.Command{
	Use:   "timeout",
	Short: "Configure global timeout and retry settings",
	Long:  "Configure global timeout and retry settings for API operations",
}

// timeoutSetCmd sets timeout configuration
var timeoutSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set timeout configuration",
	Long:  "Set global timeout and retry configuration for API operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		timeout, _ := cmd.Flags().GetDuration("timeout")
		retryCount, _ := cmd.Flags().GetInt("retry")
		retryDelay, _ := cmd.Flags().GetDuration("retry-delay")

		// Validate settings
		if timeout < 1*time.Second {
			return fmt.Errorf("timeout must be at least 1 second")
		}
		if retryCount < 0 || retryCount > 10 {
			return fmt.Errorf("retry count must be between 0 and 10")
		}
		if retryDelay < 100*time.Millisecond {
			return fmt.Errorf("retry delay must be at least 100ms")
		}

		// Set global values
		globalTimeout = timeout
		globalRetryCount = retryCount
		globalRetryDelay = retryDelay

		fmt.Printf("✅ Timeout configuration updated:\n")
		fmt.Printf("   ⏱️  Timeout: %v\n", timeout)
		fmt.Printf("   🔄 Retry Count: %d\n", retryCount)
		fmt.Printf("   ⏳ Retry Delay: %v\n", retryDelay)

		return nil
	},
}

// timeoutShowCmd shows current timeout configuration
var timeoutShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current timeout configuration",
	Long:  "Display the current global timeout and retry configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := getTimeoutConfig()

		fmt.Printf("🕐 Current Timeout Configuration\n")
		fmt.Printf("===============================\n")
		fmt.Printf("⏱️  Timeout: %v\n", config.Timeout)
		fmt.Printf("🔄 Retry Count: %d\n", config.RetryCount)
		fmt.Printf("⏳ Retry Delay: %v\n", config.RetryDelay)
		fmt.Printf("📈 Max Backoff: %v\n", config.MaxBackoff)

		return nil
	},
}

// getTimeoutConfig returns the current timeout configuration
func getTimeoutConfig() *TimeoutConfig {
	config := DefaultTimeoutConfig()

	// Apply global overrides if set
	if globalTimeout > 0 {
		config.Timeout = globalTimeout
	}
	if globalRetryCount >= 0 {
		config.RetryCount = globalRetryCount
	}
	if globalRetryDelay > 0 {
		config.RetryDelay = globalRetryDelay
	}

	return config
}

// parseTimeoutFlags extracts timeout configuration from command flags
func parseTimeoutFlags(cmd *cobra.Command) *TimeoutConfig {
	config := getTimeoutConfig()

	// Override with command-specific flags if provided
	if cmd.Flags().Changed("timeout") {
		if timeout, err := cmd.Flags().GetDuration("timeout"); err == nil && timeout > 0 {
			config.Timeout = timeout
		}
	}
	if cmd.Flags().Changed("retry") {
		if retry, err := cmd.Flags().GetInt("retry"); err == nil && retry >= 0 {
			config.RetryCount = retry
		}
	}
	if cmd.Flags().Changed("retry-delay") {
		if retryDelay, err := cmd.Flags().GetDuration("retry-delay"); err == nil && retryDelay > 0 {
			config.RetryDelay = retryDelay
		}
	}

	return config
}

// addTimeoutFlags adds timeout and retry flags to a command
func addTimeoutFlags(cmd *cobra.Command) {
	cmd.Flags().Duration("timeout", 0, "Request timeout (0 = use global default)")
	cmd.Flags().Int("retry", -1, "Number of retries (-1 = use global default)")
	cmd.Flags().Duration("retry-delay", 0, "Delay between retries (0 = use global default)")
}

func init() {
	// Add subcommands
	timeoutCmd.AddCommand(timeoutSetCmd)
	timeoutCmd.AddCommand(timeoutShowCmd)

	// Flags for timeout set command
	timeoutSetCmd.Flags().Duration("timeout", 30*time.Second, "Request timeout")
	timeoutSetCmd.Flags().Int("retry", 3, "Number of retries")
	timeoutSetCmd.Flags().Duration("retry-delay", 1*time.Second, "Delay between retries")

	// Mark required flags
	_ = timeoutSetCmd.MarkFlagRequired("timeout")
}
