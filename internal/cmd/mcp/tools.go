package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func handleListIssues(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		q := url.Values{}

		if projectID := req.GetInt("project_id", 0); projectID > 0 {
			q.Set("project_id", strconv.Itoa(projectID))
		}
		q.Set("status_id", req.GetString("status_id", "open"))
		q.Set("limit", strconv.Itoa(clampLimit(req.GetInt("limit", 25))))
		if offset := req.GetInt("offset", 0); offset > 0 {
			q.Set("offset", strconv.Itoa(offset))
		}
		if subject := req.GetString("subject", ""); subject != "" {
			q.Set("subject", "~"+subject)
		}
		if assignedTo := req.GetString("assigned_to_id", ""); assignedTo != "" {
			q.Set("assigned_to_id", assignedTo)
		}

		path := "/issues.json?" + q.Encode()
		body, status, err := api.ClientGET(r, path)
		return handleAPIResult(body, status, err)
	}
}

func handleGetIssue(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return handleToolError("Missing required argument: id")
		}

		include := "allowed_statuses"
		if req.GetBool("include_journals", false) {
			include = "journals,allowed_statuses"
		}

		path := fmt.Sprintf("/issues/%d.json?include=%s", id, include)
		body, status, apiErr := api.ClientGET(r, path)
		return handleAPIResult(body, status, apiErr)
	}
}

func handleCreateIssue(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		projectID, err := req.RequireInt("project_id")
		if err != nil {
			return handleToolError("Missing required argument: project_id")
		}
		subject, err := req.RequireString("subject")
		if err != nil {
			return handleToolError("Missing required argument: subject")
		}

		issue := map[string]interface{}{
			"project_id": projectID,
			"subject":    subject,
		}

		if desc := req.GetString("description", ""); desc != "" {
			issue["description"] = desc
		}
		if trackerID := req.GetInt("tracker_id", 0); trackerID > 0 {
			issue["tracker_id"] = trackerID
		}
		if priorityID := req.GetInt("priority_id", 0); priorityID > 0 {
			issue["priority_id"] = priorityID
		}
		if assignedTo := req.GetInt("assigned_to_id", 0); assignedTo > 0 {
			issue["assigned_to_id"] = assignedTo
		}
		if statusID := req.GetInt("status_id", 0); statusID > 0 {
			issue["status_id"] = statusID
		}
		if versionID := req.GetInt("fixed_version_id", 0); versionID > 0 {
			issue["fixed_version_id"] = versionID
		}
		if parentID := req.GetInt("parent_issue_id", 0); parentID > 0 {
			issue["parent_issue_id"] = parentID
		}

		payload := map[string]interface{}{"issue": issue}
		body, err := json.Marshal(payload)
		if err != nil {
			return handleToolError("Failed to serialize request: %s", err.Error())
		}

		resp, status, apiErr := api.ClientPOST(r, "/issues.json", body)
		if apiErr != nil {
			return mcpgo.NewToolResultError(fmt.Sprintf("API error (HTTP %d): %s", status, apiErr.Error())), nil
		}
		if status != 201 {
			errors := api.ParseResponseError(resp)
			return mcpgo.NewToolResultError(fmt.Sprintf("Could not create issue: %v", errors.Errors)), nil
		}
		return mcpgo.NewToolResultText(string(resp)), nil
	}
}

func handleUpdateIssue(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return handleToolError("Missing required argument: id")
		}

		issue := map[string]interface{}{}

		if subject := req.GetString("subject", ""); subject != "" {
			issue["subject"] = subject
		}
		if _, ok := req.GetArguments()["description"]; ok {
			issue["description"] = req.GetString("description", "")
		}
		if statusID := req.GetInt("status_id", 0); statusID > 0 {
			issue["status_id"] = statusID
		}
		if priorityID := req.GetInt("priority_id", 0); priorityID > 0 {
			issue["priority_id"] = priorityID
		}
		if assignedTo := req.GetInt("assigned_to_id", 0); assignedTo > 0 {
			issue["assigned_to_id"] = assignedTo
		}
		if trackerID := req.GetInt("tracker_id", 0); trackerID > 0 {
			issue["tracker_id"] = trackerID
		}
		if versionID := req.GetInt("fixed_version_id", 0); versionID > 0 {
			issue["fixed_version_id"] = versionID
		}
		if doneRatio := req.GetInt("done_ratio", -1); doneRatio >= 0 {
			issue["done_ratio"] = doneRatio
		}
		if notes := req.GetString("notes", ""); notes != "" {
			issue["notes"] = notes
		}

		if len(issue) == 0 {
			return handleToolError("No fields to update were provided")
		}

		payload := map[string]interface{}{"issue": issue}
		body, err := json.Marshal(payload)
		if err != nil {
			return handleToolError("Failed to serialize request: %s", err.Error())
		}

		path := fmt.Sprintf("/issues/%d.json", id)
		resp, status, apiErr := api.ClientPUT(r, path, body)
		if apiErr != nil {
			return mcpgo.NewToolResultError(fmt.Sprintf("API error (HTTP %d): %s", status, apiErr.Error())), nil
		}
		if status < 200 || status > 299 {
			errors := api.ParseResponseError(resp)
			return mcpgo.NewToolResultError(fmt.Sprintf("Could not update issue: %v", errors.Errors)), nil
		}
		return mcpgo.NewToolResultText(fmt.Sprintf("Issue #%d updated successfully", id)), nil
	}
}

func handleAddNote(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return handleToolError("Missing required argument: id")
		}
		notes, err := req.RequireString("notes")
		if err != nil {
			return handleToolError("Missing required argument: notes")
		}

		issue := map[string]interface{}{
			"notes": notes,
		}
		if private := req.GetBool("private", false); private {
			issue["private_notes"] = true
		}

		payload := map[string]interface{}{"issue": issue}
		body, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return handleToolError("Failed to serialize request: %s", marshalErr.Error())
		}

		path := fmt.Sprintf("/issues/%d.json", id)
		resp, status, apiErr := api.ClientPUT(r, path, body)
		if apiErr != nil {
			return mcpgo.NewToolResultError(fmt.Sprintf("API error (HTTP %d): %s", status, apiErr.Error())), nil
		}
		if status < 200 || status > 299 {
			errors := api.ParseResponseError(resp)
			return mcpgo.NewToolResultError(fmt.Sprintf("Could not add note: %v", errors.Errors)), nil
		}
		return mcpgo.NewToolResultText(fmt.Sprintf("Note added to issue #%d", id)), nil
	}
}

func handleListProjects(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		limit := clampLimit(req.GetInt("limit", 25))
		offset := req.GetInt("offset", 0)
		if offset < 0 {
			offset = 0
		}

		path := fmt.Sprintf("/projects.json?limit=%d&offset=%d", limit, offset)
		body, status, err := api.ClientGET(r, path)
		return handleAPIResult(body, status, err)
	}
}

func handleGetCurrentUser(r *config.Red_t) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		body, status, err := api.ClientGET(r, "/users/current.json")
		return handleAPIResult(body, status, err)
	}
}

// clampLimit keeps the page size within Redmine's accepted range (1-100).
func clampLimit(limit int) int {
	if limit < 1 {
		return 25
	}
	if limit > 100 {
		return 100
	}
	return limit
}
