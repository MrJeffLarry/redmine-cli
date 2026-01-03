package global

import (
	"github.com/MrJeffLarry/redmine-cli/internal/util"
)

type issueCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Project   util.IdName `json:"project"`
	AssignedTo   util.IdName `json:"assignedTo"`
}

type issueCategories struct {
	IssueCategories []issueCategory `json:"issue_categories"`
}

type issuePriorities struct {
	IssuePriorities []issuePriority `json:"issue_priorities"`
}

type issueStatusHolder struct {
	IssueStatus []IssueStatus `json:"issue_statuses,omitempty"`
}

type issuePriority struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Active    bool   `json:"active"`
}

type IssueStatus struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
}
