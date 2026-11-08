package cmd

import (
	"testing"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/spf13/cobra"
)

func TestAuthLoginUserNoServer(t *testing.T) {
	r := CmdInit("Version", "GitCommit", "BuildTime")
	r.Cmd.SetArgs([]string{"auth", "login", "--username", "test", "--password", "test"})
	r.Cmd.Execute()
}

func TestMCPDisablesDebugOutput(t *testing.T) {
	r := CmdInit("Version", "GitCommit", "BuildTime")
	r.Config.Version = "2"
	r.Config.Servers = []config.Server_t{{Name: "test"}}

	mcpCmd, _, err := r.Cmd.Find([]string{"mcp"})
	if err != nil {
		t.Fatal(err)
	}
	mcpCmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
	r.Cmd.SetArgs([]string{"mcp", "--debug", "--rid", "test"})
	if err := r.Cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if r.Debug {
		t.Error("debug output must be disabled for the MCP command")
	}
	if r.Server == nil || r.Server.Name != "test" {
		t.Error("MCP command should still select the requested server")
	}
}
