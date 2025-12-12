package config

import (
	"os"
	"testing"
	"time"

	"github.com/briandowns/spinner"
)

func TestMultiInstanceSetup(t *testing.T) {
	r := InitConfig()
	r.Test = true

	// Test setting RID enables multi-mode
	r.SetRID("1")
	if !r.UseMultiMode {
		t.Error("Expected UseMultiMode to be true after SetRID")
	}

	if r.MultiConfig.Instances == nil {
		t.Error("Expected Instances map to be initialized")
	}

	if r.MultiConfig.DefaultInstance != "1" {
		t.Errorf("Expected default instance to be '1', got '%s'", r.MultiConfig.DefaultInstance)
	}
}

func TestMultiInstanceSave(t *testing.T) {
	r := InitConfig()
	r.Test = true
	r.SetRID("2")

	r.SetServer("https://redmine2.example.com")
	r.SetApiKey("test-key-2")
	r.SetUserID(2)

	// Verify config is set
	if r.Config.Server != "https://redmine2.example.com" {
		t.Errorf("Expected server 'https://redmine2.example.com', got '%s'", r.Config.Server)
	}

	// Verify multi-mode is enabled
	if !r.UseMultiMode {
		t.Error("Expected UseMultiMode to be true")
	}
}

func TestMultiInstanceClearAll(t *testing.T) {
	r := InitConfig()
	r.Test = true
	
	// Add multiple instances
	r.SetRID("1")
	r.SetServer("https://redmine1.example.com")
	r.SetApiKey("key1")
	r.Save()
	
	r.SetRID("2")
	r.SetServer("https://redmine2.example.com")
	r.SetApiKey("key2")
	r.Save()

	// Clear instance 2
	r.ClearAll()

	if r.Config.Server != "" {
		t.Error("Expected server to be cleared")
	}

	// Verify instance 2 is removed from map
	if _, exists := r.MultiConfig.Instances["2"]; exists {
		t.Error("Expected instance 2 to be removed from map")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	r := &Red_t{
		Spinner: spinner.New(spinner.CharSets[11], 100*time.Millisecond),
	}
	r.Test = true
	r.MultiConfig.Instances = make(map[string]Config_t)

	// Simulate legacy single-instance config
	r.UseMultiMode = false
	r.Config.Server = "https://legacy.example.com"
	r.Config.ApiKey = "legacy-key"

	// Verify it still works
	if r.Config.Server != "https://legacy.example.com" {
		t.Error("Expected legacy config to work")
	}

	if r.UseMultiMode {
		t.Error("Expected UseMultiMode to be false for legacy config")
	}
}

func TestSetRIDMultipleTimes(t *testing.T) {
	r := InitConfig()
	r.Test = true

	// Set first instance
	r.SetRID("1")
	r.SetServer("https://redmine1.example.com")
	r.SetApiKey("key1")
	
	if r.MultiConfig.DefaultInstance != "1" {
		t.Errorf("Expected default instance '1', got '%s'", r.MultiConfig.DefaultInstance)
	}

	// Set second instance - should NOT change default
	r.SetRID("2")
	if r.MultiConfig.DefaultInstance != "1" {
		t.Errorf("Expected default instance to remain '1', got '%s'", r.MultiConfig.DefaultInstance)
	}
}

func TestIsConfigBadMultiMode(t *testing.T) {
	r := &Red_t{
		Spinner: spinner.New(spinner.CharSets[11], 100*time.Millisecond),
	}
	r.Test = true
	r.MultiConfig.Instances = make(map[string]Config_t)

	// Empty config should be bad
	if !r.IsConfigBad() {
		t.Error("Expected empty config to be bad")
	}

	// Set valid config
	r.SetRID("1")
	r.SetServer("https://redmine.example.com")
	r.SetApiKey("valid-key")

	if r.IsConfigBad() {
		t.Error("Expected valid config to be good")
	}
}

func TestDefaultInstanceEnvironment(t *testing.T) {
	// Clean up any existing config first
	os.Unsetenv(RED_CONFIG_REDMINE_URL)
	os.Unsetenv(RED_CONFIG_REDMINE_API_KEY)
	
	// Set environment variable
	os.Setenv(RED_CONFIG_REDMINE_URL, "https://env.example.com")
	os.Setenv(RED_CONFIG_REDMINE_API_KEY, "env-key")
	defer os.Unsetenv(RED_CONFIG_REDMINE_URL)
	defer os.Unsetenv(RED_CONFIG_REDMINE_API_KEY)

	r := &Red_t{
		Spinner: spinner.New(spinner.CharSets[11], 100*time.Millisecond),
	}
	r.Test = true
	r.MultiConfig.Instances = make(map[string]Config_t)
	
	r.Config.Server = exEnv(RED_CONFIG_REDMINE_URL, "")
	r.Config.ApiKey = exEnv(RED_CONFIG_REDMINE_API_KEY, "")

	// Verify environment variables are loaded
	if r.Config.Server != "https://env.example.com" {
		t.Errorf("Expected server from env, got '%s'", r.Config.Server)
	}

	if r.Config.ApiKey != "env-key" {
		t.Errorf("Expected API key from env, got '%s'", r.Config.ApiKey)
	}
}
