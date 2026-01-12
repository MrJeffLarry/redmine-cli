package config

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	if r.Config.DefaultServer != 1 {
		t.Errorf("Expected default server to 1, got %d", r.Config.DefaultServer)
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

func TestSaveMethod_ErrorHandling(t *testing.T) {
	r := &Red_t{}
	r.Config = ConfigV2_t{Version: "2.0"}
	r.Test = false // Save should try to write
	// Use a temp dir for config file
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)
	// Remove config file to force creation
	os.RemoveAll(dir + "/.red")
	err := r.Save()
	if err != nil {
		t.Errorf("Save() failed: %v", err)
	}
}

func TestSaveLocalProject(t *testing.T) {
	r := &Red_t{}
	r.Test = false
	dir := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(old)
	err := r.SaveLocalProject(42)
	if err != nil {
		t.Errorf("SaveLocalProject failed: %v", err)
	}
}

func TestRemoveServerByNameAndID(t *testing.T) {
	r := &Red_t{}
	r.AddServer("one", "url1", "key1", "p1", 1, 1)
	r.AddServer("two", "url2", "key2", "p2", 2, 2)
	err := r.RemoveServer(0, "")
	if err != nil {
		t.Errorf("RemoveServer by ID failed: %v", err)
	}
	r.AddServer("three", "url3", "key3", "p3", 3, 3)
	err = r.RemoveServer(-1, "three")
	if err != nil {
		t.Errorf("RemoveServer by name failed: %v", err)
	}
}

func TestRemoveCurrentServer_Empty(t *testing.T) {
	r := &Red_t{}
	r.AddServer("one", "url1", "key1", "p1", 1, 1)
	r.RemoveCurrentServer()
	if r.Server != nil {
		t.Error("Expected Server to be nil after removing only server")
	}
}

func TestSwitchDefaultServer(t *testing.T) {
	r := &Red_t{}
	r.AddServer("one", "url1", "key1", "p1", 1, 1)
	r.AddServer("two", "url2", "key2", "p2", 2, 2)
	err := r.SetDefaultServer(0)
	if err != nil {
		t.Errorf("SwitchDefaultServer failed: %v", err)
	}
	if r.Config.DefaultServer != 0 {
		t.Errorf("Expected DefaultServer to be 0, got %d", r.Config.DefaultServer)
	}

	err = r.SetDefaultServer(1)
	if err != nil {
		t.Errorf("SwitchDefaultServer failed: %v", err)
	}
	r.RemoveCurrentServer()
	if r.Config.DefaultServer != 0 {
		t.Errorf("Expected DefaultServer to reset to 0 after removal, got %d", r.Config.DefaultServer)
	}
	r.RemoveCurrentServer()
	err = r.SetDefaultServer(0)
	if err == nil {
		t.Error("Expected error when setting default server on empty config")
	}
}
func TestInitConfig_EmptyConfig(t *testing.T) {
	// Clean up environment variables
	os.Unsetenv(RED_CONFIG_REDMINE_URL)
	os.Unsetenv(RED_CONFIG_REDMINE_API_KEY)
	os.Unsetenv(RED_CONFIG_EDITOR)

	// Set up temp directory for config
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)

	red := InitConfig()

	// Verify basic initialization
	if red == nil {
		t.Fatal("InitConfig returned nil")
	}
	if red.Client == nil {
		t.Error("Expected Client to be initialized")
	}
	if red.Spinner == nil {
		t.Error("Expected Spinner to be initialized")
	}
	if red.Server != nil {
		t.Error("Expected Server to be nil with empty config")
	}
}

func TestInitConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv(RED_CONFIG_REDMINE_URL, "https://env.example.com")
	os.Setenv(RED_CONFIG_REDMINE_API_KEY, "env-api-key")
	os.Setenv(RED_CONFIG_REDMINE_PROJECT, "env-project")
	os.Setenv(RED_CONFIG_REDMINE_PROJECT_ID, "123")
	os.Setenv(RED_CONFIG_REDMINE_USER_ID, "456")
	os.Setenv(RED_CONFIG_EDITOR, "vim")
	os.Setenv(RED_CONFIG_PAGER, "less")
	defer func() {
		os.Unsetenv(RED_CONFIG_REDMINE_URL)
		os.Unsetenv(RED_CONFIG_REDMINE_API_KEY)
		os.Unsetenv(RED_CONFIG_REDMINE_PROJECT)
		os.Unsetenv(RED_CONFIG_REDMINE_PROJECT_ID)
		os.Unsetenv(RED_CONFIG_REDMINE_USER_ID)
		os.Unsetenv(RED_CONFIG_EDITOR)
		os.Unsetenv(RED_CONFIG_PAGER)
	}()

	// Set up temp directory for config
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)

	red := InitConfig()

	if red == nil {
		t.Fatal("InitConfig returned nil")
	}

	// Verify environment variables are loaded into config
	if red.Config.Editor != "vim" {
		t.Errorf("Expected Editor 'vim', got '%s'", red.Config.Editor)
	}
	if red.Config.Pager != "less" {
		t.Errorf("Expected Pager 'less', got '%s'", red.Config.Pager)
	}
}

func TestInitConfig_WithExistingGlobalConfig(t *testing.T) {
	// Set up temp directory and create config file
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)

	// Create a valid global config
	config := ConfigV2_t{
		Version: "2.0",
		Servers: []Server_t{
			{
				Name:      "test-server",
				Server:    "https://test.example.com",
				ApiKey:    "test-key",
				Project:   "test-project",
				ProjectID: 789,
				UserID:    101,
			},
		},
		DefaultServer: 0,
		Editor:        "nano",
		Pager:         "more",
	}

	// Save config to file using filepath.Join and check errors
	redDir := filepath.Join(dir, ".red")
	err := os.Mkdir(redDir, 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create .red dir: %v", err)
	}
	configData, _ := json.Marshal(config)
	err = os.WriteFile(filepath.Join(redDir, "config.json"), configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	red := InitConfig()

	if red.Server == nil {
		t.Fatal("Expected Server to be set")
	}
	if red.Server.Name != "test-server" {
		t.Errorf("Expected Server name 'test-server', got '%s'", red.Server.Name)
	}
	if red.Config.Editor != "nano" {
		t.Errorf("Expected Editor 'nano', got '%s'", red.Config.Editor)
	}
}

func TestInitConfig_WithLocalConfigOverride(t *testing.T) {
	// Set up temp directory
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	oldCwd, _ := os.Getwd()
	os.Setenv("HOME", dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Chdir(oldCwd)
	}()

	// Create global config
	globalConfig := ConfigV2_t{
		Version: "2.0",
		Servers: []Server_t{
			{
				Name:      "global-server",
				Server:    "https://global.example.com",
				ApiKey:    "global-key",
				Project:   "global-project",
				ProjectID: 100,
				UserID:    200,
			},
		},
		DefaultServer: 0,
		Editor:        "global-editor",
		Pager:         "global-pager",
	}

	// Create local config that overrides some values
	localConfig := ConfigLocal_t{
		Server:    "https://local.example.com",
		ApiKey:    "local-key",
		ProjectID: 999,
		Editor:    "local-editor",
	}

	// Save configs using filepath.Join for cross-platform compatibility
	redDir := filepath.Join(dir, ".red")
	err = os.Mkdir(redDir, 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create global .red dir: %v", err)
	}
	globalData, _ := json.Marshal(globalConfig)
	err = os.WriteFile(filepath.Join(redDir, "config.json"), globalData, 0644)
	if err != nil {
		t.Fatalf("Failed to write global config: %v", err)
	}

	err = os.Mkdir(".red", 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create local .red dir: %v", err)
	}
	localData, _ := json.Marshal(localConfig)
	err = os.WriteFile(filepath.Join(".red", "config.json"), localData, 0644)
	if err != nil {
		t.Fatalf("Failed to write local config: %v", err)
	}

	red := InitConfig()

	if red.Server == nil {
		t.Fatal("Expected Server to be set, got nil")
	}

	// Verify local config overrides global config
	if red.Server.Server != "https://local.example.com" {
		t.Errorf("Expected local server override, got '%s'", red.Server.Server)
	}
	if red.Server.ApiKey != "local-key" {
		t.Errorf("Expected local API key override, got '%s'", red.Server.ApiKey)
	}
	if red.Server.ProjectID != 999 {
		t.Errorf("Expected local project ID override, got %d", red.Server.ProjectID)
	}
	if red.Config.Editor != "local-editor" {
		t.Errorf("Expected local editor override, got '%s'", red.Config.Editor)
	}
}

func TestInitConfig_MigrateV1ToV2(t *testing.T) {
	// Set up temp directory
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)

	// Create a v1 config (no version field)
	v1Config := ConfigV1_t{
		Server:    "https://v1.example.com",
		ApiKey:    "v1-key",
		Project:   "v1-project",
		ProjectID: 111,
		UserID:    222,
		Editor:    "v1-editor",
		Pager:     "v1-pager",
	}

	// Save v1 config using filepath.Join for cross-platform compatibility
	redDir := filepath.Join(dir, ".red")
	err := os.Mkdir(redDir, 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create .red dir: %v", err)
	}
	v1Data, _ := json.Marshal(v1Config)
	err = os.WriteFile(filepath.Join(redDir, "config.json"), v1Data, 0644)
	if err != nil {
		t.Fatalf("Failed to write v1 config: %v", err)
	}

	red := InitConfig()

	// Verify migration to v2
	if red.Config.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", red.Config.Version)
	}
	if len(red.Config.Servers) != 1 {
		t.Errorf("Expected 1 server after migration, got %d", len(red.Config.Servers))
	}
	if len(red.Config.Servers) > 0 && red.Config.Servers[0].Name != "default" {
		t.Errorf("Expected migrated server name 'default', got '%s'", red.Config.Servers[0].Name)
	}
	if red.Server == nil {
		t.Fatal("Expected Server to be set, got nil")
	}
	if red.Server.Server != "https://v1.example.com" {
		t.Errorf("Expected migrated server URL, got '%s'", red.Server.Server)
	}
}
