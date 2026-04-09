package mcp

import (
	"fmt"
	"io"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/briandowns/spinner"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

func NewCmdMCP(r *config.Red_t, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start an MCP server",
		Long: `Run red-cli as a Model Context Protocol (MCP) server over stdio.

This command exposes Redmine operations as MCP tools so AI assistants can
retrieve and update Redmine data through your existing red-cli authentication.

Available tools:
	- list_issues
	- get_issue
	- create_issue
	- update_issue
	- add_note
	- list_projects
	- get_current_user

The server communicates using JSON-RPC on stdin/stdout and is intended to be
started by an MCP client configuration (for example Claude Desktop).`,
		Example: "# Run MCP server directly (for local testing)\n" +
			"red-cli mcp\n\n" +
			"# Example MCP client configuration\n" +
			"{\n" +
			"  \"mcpServers\": {\n" +
			"    \"redmine\": {\n" +
			"      \"command\": \"red-cli\",\n" +
			"      \"args\": [\"mcp\"]\n" +
			"    }\n" +
			"  }\n" +
			"}",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer(r, version)
		},
	}
	return cmd
}

func runMCPServer(r *config.Red_t, version string) error {
	// Disable spinner output so it doesn't interfere with the JSON-RPC stdio channel
	r.Spinner = spinner.New(spinner.CharSets[11], 100, spinner.WithWriter(io.Discard))

	s := mcpserver.NewMCPServer(
		"redmine-mcp",
		version,
		mcpserver.WithToolCapabilities(false),
	)

	// Issue tools
	s.AddTool(toolListIssues(), handleListIssues(r))
	s.AddTool(toolGetIssue(), handleGetIssue(r))
	s.AddTool(toolCreateIssue(), handleCreateIssue(r))
	s.AddTool(toolUpdateIssue(), handleUpdateIssue(r))
	s.AddTool(toolAddNote(), handleAddNote(r))

	// Project tools
	s.AddTool(toolListProjects(), handleListProjects(r))

	// User tools
	s.AddTool(toolGetCurrentUser(), handleGetCurrentUser(r))

	fmt.Fprintf(io.Discard, "Starting Redmine MCP server v%s\n", version)

	return mcpserver.ServeStdio(s)
}

func toolListIssues() mcpgo.Tool {
	return mcpgo.NewTool("list_issues",
		mcpgo.WithDescription("List Redmine issues with optional filters. Returns a JSON array of issues."),
		mcpgo.WithNumber("project_id",
			mcpgo.Description("Project ID to filter issues by (optional)"),
		),
		mcpgo.WithString("status_id",
			mcpgo.Description("Filter by status: 'open' (default), 'closed', '*' for all, or a numeric status ID"),
		),
		mcpgo.WithNumber("limit",
			mcpgo.Description("Number of issues to return (default 25, max 100)"),
		),
		mcpgo.WithNumber("offset",
			mcpgo.Description("Offset for pagination (default 0)"),
		),
		mcpgo.WithString("subject",
			mcpgo.Description("Filter issues by subject text (partial match)"),
		),
		mcpgo.WithString("assigned_to_id",
			mcpgo.Description("Filter by assignee: user ID or 'me' for current user"),
		),
	)
}

func toolGetIssue() mcpgo.Tool {
	return mcpgo.NewTool("get_issue",
		mcpgo.WithDescription("Get a Redmine issue by ID, including its details such as status, priority, description, and optionally journal notes."),
		mcpgo.WithNumber("id",
			mcpgo.Description("Issue ID"),
			mcpgo.Required(),
		),
		mcpgo.WithBoolean("include_journals",
			mcpgo.Description("Include journal/notes history (default false)"),
		),
	)
}

func toolCreateIssue() mcpgo.Tool {
	return mcpgo.NewTool("create_issue",
		mcpgo.WithDescription("Create a new Redmine issue in a project."),
		mcpgo.WithNumber("project_id",
			mcpgo.Description("Project ID where the issue will be created"),
			mcpgo.Required(),
		),
		mcpgo.WithString("subject",
			mcpgo.Description("Issue subject/title"),
			mcpgo.Required(),
		),
		mcpgo.WithString("description",
			mcpgo.Description("Issue description (Textile or Markdown)"),
		),
		mcpgo.WithNumber("tracker_id",
			mcpgo.Description("Tracker ID (e.g. Bug, Feature, Task)"),
		),
		mcpgo.WithNumber("priority_id",
			mcpgo.Description("Priority ID"),
		),
		mcpgo.WithNumber("assigned_to_id",
			mcpgo.Description("User ID to assign the issue to"),
		),
		mcpgo.WithNumber("status_id",
			mcpgo.Description("Status ID for the new issue"),
		),
		mcpgo.WithNumber("fixed_version_id",
			mcpgo.Description("Target version/milestone ID"),
		),
		mcpgo.WithNumber("parent_issue_id",
			mcpgo.Description("Parent issue ID for subtask relationship"),
		),
	)
}

func toolUpdateIssue() mcpgo.Tool {
	return mcpgo.NewTool("update_issue",
		mcpgo.WithDescription("Update an existing Redmine issue's fields. Only provided fields are changed."),
		mcpgo.WithNumber("id",
			mcpgo.Description("Issue ID to update"),
			mcpgo.Required(),
		),
		mcpgo.WithString("subject",
			mcpgo.Description("New issue subject/title"),
		),
		mcpgo.WithString("description",
			mcpgo.Description("New issue description"),
		),
		mcpgo.WithNumber("status_id",
			mcpgo.Description("New status ID"),
		),
		mcpgo.WithNumber("priority_id",
			mcpgo.Description("New priority ID"),
		),
		mcpgo.WithNumber("assigned_to_id",
			mcpgo.Description("User ID to assign the issue to"),
		),
		mcpgo.WithNumber("tracker_id",
			mcpgo.Description("New tracker ID"),
		),
		mcpgo.WithNumber("fixed_version_id",
			mcpgo.Description("Target version/milestone ID"),
		),
		mcpgo.WithNumber("done_ratio",
			mcpgo.Description("Percent done (0-100)"),
		),
		mcpgo.WithString("notes",
			mcpgo.Description("Journal note to add alongside the update"),
		),
	)
}

func toolAddNote() mcpgo.Tool {
	return mcpgo.NewTool("add_note",
		mcpgo.WithDescription("Add a journal note/comment to a Redmine issue."),
		mcpgo.WithNumber("id",
			mcpgo.Description("Issue ID"),
			mcpgo.Required(),
		),
		mcpgo.WithString("notes",
			mcpgo.Description("Note text to add"),
			mcpgo.Required(),
		),
		mcpgo.WithBoolean("private",
			mcpgo.Description("Post the note as private (default false)"),
		),
	)
}

func toolListProjects() mcpgo.Tool {
	return mcpgo.NewTool("list_projects",
		mcpgo.WithDescription("List Redmine projects accessible to the authenticated user."),
		mcpgo.WithNumber("limit",
			mcpgo.Description("Number of projects to return (default 25, max 100)"),
		),
		mcpgo.WithNumber("offset",
			mcpgo.Description("Offset for pagination (default 0)"),
		),
	)
}

func toolGetCurrentUser() mcpgo.Tool {
	return mcpgo.NewTool("get_current_user",
		mcpgo.WithDescription("Get information about the currently authenticated Redmine user."),
	)
}

// handleToolError returns a tool error result from an API error.
func handleToolError(format string, args ...interface{}) (*mcpgo.CallToolResult, error) {
	return mcpgo.NewToolResultError(fmt.Sprintf(format, args...)), nil
}

// handleAPIResult returns raw JSON body as a text tool result, or an error result on failure.
func handleAPIResult(body []byte, status int, err error) (*mcpgo.CallToolResult, error) {
	if err != nil {
		return mcpgo.NewToolResultError(fmt.Sprintf("API error (HTTP %d): %s", status, err.Error())), nil
	}
	return mcpgo.NewToolResultText(string(body)), nil
}
