package mcp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/briandowns/spinner"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

type recorded struct {
	method string
	path   string
	query  string
	key    string
	body   map[string]map[string]interface{}
}

func setup(t *testing.T, status int, resp string) (*config.Red_t, *recorded) {
	t.Helper()
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec.method = req.Method
		rec.path = req.URL.Path
		rec.query = req.URL.RawQuery
		rec.key = req.Header.Get("X-Redmine-API-Key")
		if b, _ := io.ReadAll(req.Body); len(b) > 0 {
			_ = json.Unmarshal(b, &rec.body)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(resp))
	}))
	t.Cleanup(srv.Close)

	r := &config.Red_t{
		Client:  srv.Client(),
		Spinner: spinner.New(spinner.CharSets[11], 100, spinner.WithWriter(io.Discard)),
		Server:  &config.Server_t{Server: srv.URL, ApiKey: "secret"},
	}
	return r, rec
}

func call(t *testing.T, r *config.Red_t, name string, args map[string]interface{}) *mcpgo.CallToolResult {
	t.Helper()
	s := newServer(r, "test")
	tool := s.GetTool(name)
	if tool == nil {
		t.Fatalf("tool %q not registered", name)
	}
	req := mcpgo.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return res
}

func text(t *testing.T, res *mcpgo.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatal("empty content")
	}
	tc, ok := res.Content[0].(mcpgo.TextContent)
	if !ok {
		t.Fatalf("unexpected content type %T", res.Content[0])
	}
	return tc.Text
}

func TestToolsRegistered(t *testing.T) {
	r, _ := setup(t, 200, "{}")
	s := newServer(r, "test")
	for _, n := range []string{"list_issues", "get_issue", "create_issue", "update_issue", "add_note", "list_projects", "get_current_user"} {
		if s.GetTool(n) == nil {
			t.Errorf("missing tool %s", n)
		}
	}
}

func TestListIssuesEscapesAndClamps(t *testing.T) {
	r, rec := setup(t, 200, `{"issues":[]}`)
	res := call(t, r, "list_issues", map[string]interface{}{
		"subject":        "a&b=c",
		"assigned_to_id": "me&project_id=9",
		"limit":          1000,
		"project_id":     5,
	})
	if res.IsError {
		t.Fatal(text(t, res))
	}
	if rec.key != "secret" || rec.path != "/issues.json" {
		t.Errorf("bad request: %+v", rec)
	}
	want := "assigned_to_id=me%26project_id%3D9&limit=100&project_id=5&status_id=open&subject=~a%26b%3Dc"
	if rec.query != want {
		t.Errorf("query = %q, want %q", rec.query, want)
	}
}

func TestListIssuesInvalidLimit(t *testing.T) {
	r, rec := setup(t, 200, `{}`)
	call(t, r, "list_issues", map[string]interface{}{"limit": -3})
	if rec.query != "limit=25&status_id=open" {
		t.Errorf("query = %q", rec.query)
	}
}

func TestListProjectsNegativeOffset(t *testing.T) {
	r, rec := setup(t, 200, `{}`)
	call(t, r, "list_projects", map[string]interface{}{"offset": -5, "limit": 500})
	if rec.query != "limit=100&offset=0" {
		t.Errorf("query = %q", rec.query)
	}
}

func TestGetIssue(t *testing.T) {
	r, rec := setup(t, 200, `{"issue":{"id":7}}`)
	res := call(t, r, "get_issue", map[string]interface{}{"id": 7, "include_journals": true})
	if res.IsError || text(t, res) != `{"issue":{"id":7}}` {
		t.Fatalf("unexpected result: %v", res)
	}
	if rec.path != "/issues/7.json" || rec.query != "include=journals,allowed_statuses" {
		t.Errorf("bad request: %+v", rec)
	}
}

func TestGetIssueMissingID(t *testing.T) {
	r, _ := setup(t, 200, `{}`)
	if res := call(t, r, "get_issue", map[string]interface{}{}); !res.IsError {
		t.Error("expected error")
	}
}

func TestGetIssueNotFoundIsError(t *testing.T) {
	r, _ := setup(t, 404, `{"errors":["not found"]}`)
	res := call(t, r, "get_issue", map[string]interface{}{"id": 1})
	if !res.IsError {
		t.Errorf("expected error result, got %s", text(t, res))
	}
}

func TestUnauthorizedIsError(t *testing.T) {
	r, _ := setup(t, 401, ``)
	if res := call(t, r, "get_current_user", nil); !res.IsError {
		t.Error("expected error")
	}
}

func TestGetCurrentUser(t *testing.T) {
	r, rec := setup(t, 200, `{"user":{"id":1}}`)
	res := call(t, r, "get_current_user", nil)
	if res.IsError || rec.path != "/users/current.json" {
		t.Errorf("bad: %+v %v", rec, res)
	}
}

func TestCreateIssue(t *testing.T) {
	r, rec := setup(t, 201, `{"issue":{"id":10}}`)
	res := call(t, r, "create_issue", map[string]interface{}{
		"project_id": 2, "subject": "Hi", "tracker_id": 1, "parent_issue_id": 4,
	})
	if res.IsError {
		t.Fatal(text(t, res))
	}
	if rec.method != http.MethodPost || rec.path != "/issues.json" {
		t.Errorf("bad request: %+v", rec)
	}
	issue := rec.body["issue"]
	if issue["subject"] != "Hi" || issue["project_id"] != float64(2) || issue["tracker_id"] != float64(1) || issue["parent_issue_id"] != float64(4) {
		t.Errorf("payload = %v", issue)
	}
	if _, ok := issue["priority_id"]; ok {
		t.Error("unset fields must be omitted")
	}
}

func TestCreateIssueValidation(t *testing.T) {
	r, _ := setup(t, 201, `{}`)
	if !call(t, r, "create_issue", map[string]interface{}{"subject": "x"}).IsError {
		t.Error("missing project_id should error")
	}
	if !call(t, r, "create_issue", map[string]interface{}{"project_id": 1}).IsError {
		t.Error("missing subject should error")
	}
}

func TestCreateIssueAPIError(t *testing.T) {
	r, _ := setup(t, 422, `{"errors":["Subject cannot be blank"]}`)
	res := call(t, r, "create_issue", map[string]interface{}{"project_id": 1, "subject": "x"})
	if !res.IsError {
		t.Error("expected error")
	}
}

func TestUpdateIssue(t *testing.T) {
	r, rec := setup(t, 204, ``)
	res := call(t, r, "update_issue", map[string]interface{}{"id": 3, "done_ratio": 0, "notes": "n"})
	if res.IsError {
		t.Fatal(text(t, res))
	}
	if rec.method != http.MethodPut || rec.path != "/issues/3.json" {
		t.Errorf("bad request: %+v", rec)
	}
	issue := rec.body["issue"]
	if issue["done_ratio"] != float64(0) || issue["notes"] != "n" {
		t.Errorf("payload = %v", issue)
	}
}

func TestUpdateIssueNoFields(t *testing.T) {
	r, rec := setup(t, 204, ``)
	if !call(t, r, "update_issue", map[string]interface{}{"id": 3}).IsError {
		t.Error("expected error")
	}
	if rec.method != "" {
		t.Error("no request should be sent")
	}
}

func TestUpdateIssueAPIError(t *testing.T) {
	r, _ := setup(t, 422, `{"errors":["bad"]}`)
	if !call(t, r, "update_issue", map[string]interface{}{"id": 3, "subject": "s"}).IsError {
		t.Error("expected error")
	}
}

func TestAddNote(t *testing.T) {
	r, rec := setup(t, 204, ``)
	res := call(t, r, "add_note", map[string]interface{}{"id": 9, "notes": "hello", "private": true})
	if res.IsError {
		t.Fatal(text(t, res))
	}
	issue := rec.body["issue"]
	if issue["notes"] != "hello" || issue["private_notes"] != true {
		t.Errorf("payload = %v", issue)
	}
	if !call(t, r, "add_note", map[string]interface{}{"id": 9}).IsError {
		t.Error("missing notes should error")
	}
}
