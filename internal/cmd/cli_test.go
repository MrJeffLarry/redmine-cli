package cmd

import "testing"

func TestAuthLoginUserNoServer(t *testing.T) {
	r := CmdInit("Version", "GitCommit", "BuildTime")
	r.Cmd.SetArgs([]string{"auth", "login", "--username", "test", "--password", "test"})
	r.Cmd.Execute()
}
