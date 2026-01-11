package config

import (
	"os"
	"testing"
	"time"

	"github.com/briandowns/spinner"
)

func TestMultiInstanceSetup(t *testing.T) {
	// Build a Red_t directly and use current API (Servers slice)
	r := &Red_t{}
	r.Test = true

	// Add a server and verify servers slice and default server
	if err := r.AddServer("default", "https://redmine1.example.com", "", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	if len(r.Config.Servers) == 0 {
		t.Error("Expected Servers slice to be initialized and contain at least one server")
	}

	if r.Config.DefaultServer != 0 {
		t.Errorf("Expected default server index to be 0, got %d", r.Config.DefaultServer)
	}
}

func TestMultiInstanceSave(t *testing.T) {
	r := &Red_t{}
	r.Test = true

	// Add server and set values using available methods
	if err := r.AddServer("inst2", "https://redmine2.example.com", "", "", 0, 2); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	// Set API key for default (which is 0)
	r.Config.Servers[0].ApiKey = "multi-key-123"
	r.Config.Servers[0].Project = "multi-project-456"
	r.Config.Servers[0].UserID = 2

	// Verify config is set
	if r.Config.Servers[0].Server != "https://redmine2.example.com" {
		t.Errorf("Expected server 'https://redmine2.example.com', got '%s'", r.Config.Servers[0].Server)
	}
}

func TestMultiInstanceClearAll(t *testing.T) {
	r := &Red_t{}
	r.Test = true

	// Add multiple servers
	if err := r.AddServer("one", "https://redmine1.example.com", "key1", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}
	if err := r.AddServer("two", "https://redmine2.example.com", "key2", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	// Clear all
	r.ClearAll()

	if len(r.Config.Servers) != 0 {
		t.Error("Expected servers to be cleared")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	r := &Red_t{}
	r.Test = true

	// Legacy: empty/old config should be considered bad
	r.Config = ConfigV2_t{}
	if !r.IsConfigBad() {
		t.Error("Expected empty config to be considered bad")
	}

	// Migrate / valid config
	if err := r.AddServer("default", "https://legacy.example.com", "legacy-key", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}
	r.Config.Version = "2.0"
	if r.IsConfigBad() {
		t.Error("Expected non-empty config to be good")
	}
}

func TestSetRIDMultipleTimes(t *testing.T) {
	r := &Red_t{}
	r.Test = true

	// Add first server and ensure default remains index 0
	if err := r.AddServer("one", "https://redmine1.example.com", "key1", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}
	if r.Config.DefaultServer != 0 {
		t.Errorf("Expected default server 0, got %d", r.Config.DefaultServer)
	}

	// Add second server - default should still be 0 until explicitly changed
	if err := r.AddServer("two", "https://redmine2.example.com", "key2", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}
	if r.Config.DefaultServer != 0 {
		t.Errorf("Expected default server to remain 0, got %d", r.Config.DefaultServer)
	}
}

func TestIsConfigBadMultiMode(t *testing.T) {
	r := &Red_t{}
	r.Test = true

	// Empty config
	r.Config = ConfigV2_t{}
	if !r.IsConfigBad() {
		t.Error("Expected empty config to be bad")
	}

	// Valid config
	if err := r.AddServer("one", "https://redmine.example.com", "valid-key", "", 0, 0); err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}
	r.Config.Version = "2.0"
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

	// Load env into a server entry
	srv := Server_t{
		Server: exEnv(RED_CONFIG_REDMINE_URL, ""),
		ApiKey: exEnv(RED_CONFIG_REDMINE_API_KEY, ""),
	}
	r.Config.Servers = []Server_t{srv}

	// Verify environment variables are loaded
	if r.Config.Servers[0].Server != "https://env.example.com" {
		t.Errorf("Expected server from env, got '%s'", r.Config.Servers[0].Server)
	}

	if r.Config.Servers[0].ApiKey != "env-key" {
		t.Errorf("Expected API key from env, got '%s'", r.Config.Servers[0].ApiKey)
	}
}
