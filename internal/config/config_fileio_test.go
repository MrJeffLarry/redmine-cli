package config

import (
	"os"
	"testing"
)

func TestCreateTmpFile(t *testing.T) {
	name, err := CreateTmpFile("test body")
	if err != nil {
		t.Fatalf("CreateTmpFile failed: %v", err)
	}
	defer os.Remove(name)
	if _, err := os.Stat(name); err != nil {
		t.Errorf("Temp file not created: %v", err)
	}
}

func TestCreateFolderPath_PermissionDenied(t *testing.T) {
	// Try to create a folder in a location that should not be writable (root)
	if os.Geteuid() == 0 {
		t.Skip("Test not valid as root user")
	}
	err := createFolderPath("/root/should-not-exist")
	if err == nil {
		t.Error("Expected permission error, got nil")
	}
}

func TestConfigGlobalPath_Error(t *testing.T) {
	// Backup and unset HOME and USERPROFILE to force error cross-platform
	home := os.Getenv("HOME")
	userprofile := os.Getenv("USERPROFILE")
	defer func() {
		os.Setenv("HOME", home)
		os.Setenv("USERPROFILE", userprofile)
	}()
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")

	_, err := configGlobalPath()
	if err == nil {
		t.Skip("os.UserHomeDir() still succeeded; cannot force error on this platform/environment")
	}
}

func TestSaveAndLoadLocalConfig(t *testing.T) {
	// Use a temp dir for local config
	dir := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(old)
	cfg := ConfigLocal_t{Server: "s", ApiKey: "k", Project: "p", ProjectID: 1, UserID: 2, Editor: "e", Pager: "pg"}
	err := saveLocalConfig(cfg)
	if err != nil {
		t.Fatalf("saveLocalConfig failed: %v", err)
	}
	cfg2, err := loadLocalConfig()
	if err != nil {
		t.Fatalf("loadLocalConfig failed: %v", err)
	}
	if cfg2.Server != "s" || cfg2.ApiKey != "k" || cfg2.ProjectID != 1 {
		t.Error("Loaded config does not match saved config")
	}
}
