package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "coolifyme-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Set HOME to our temp directory
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() {
		_ = os.Setenv("HOME", originalHome)
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.BaseURL != "https://app.coolify.io/api/v1" {
		t.Errorf("Expected default BaseURL, got %s", config.BaseURL)
	}

	if config.Profile != "default" {
		t.Errorf("Expected default profile, got %s", config.Profile)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "coolifyme-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Set HOME to our temp directory
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() {
		_ = os.Setenv("HOME", originalHome)
	}()

	// Create and save a config
	cfg := &Config{
		APIToken: "test-token",
		BaseURL:  "https://test.example.com/api/v1",
		Profile:  "test",
	}

	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config back
	loadedCfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.APIToken != cfg.APIToken {
		t.Errorf("Expected APIToken %s, got %s", cfg.APIToken, loadedCfg.APIToken)
	}

	if loadedCfg.BaseURL != cfg.BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", cfg.BaseURL, loadedCfg.BaseURL)
	}
}

func TestGetConfigDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "coolifyme-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Set HOME to our temp directory
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() {
		_ = os.Setenv("HOME", originalHome)
	}()

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Failed to get config dir: %v", err)
	}

	expected := filepath.Join(tmpDir, ".config", "coolifyme")
	if configDir != expected {
		t.Errorf("Expected config dir %s, got %s", expected, configDir)
	}
}
