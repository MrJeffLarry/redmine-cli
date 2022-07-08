package issue

import "github.com/MrJeffLarry/redmine-cli/internal/util"

type viewIssue struct {
	Issue issue `json:"issue,omitempty"`
}

type newIssueHolder struct {
	Issue newIssue `json:"issue,omitempty"`
}

type issueStatusHolder struct {
	IssueStatus []issueStatus `json:"issue_statuses,omitempty"`
}

type issues struct {
	Issues     []issue `json:"issues,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

type newIssue struct {
	ProjectID      int64  `json:"project_id,omitempty"`
	TrackerID      int64  `json:"tracker_id,omitempty"`
	StatusID       int64  `json:"status_id,omitempty"`
	PriorityID     int64  `json:"priority_id,omitempty"`
	Subject        string `json:"subject,omitempty"`
	Description    string `json:"description,omitempty"`
	CategoryID     int64  `json:"category_id,omitempty"`
	FixedVersionID int64  `json:"fixed_version_id,omitempty"`
	AssignedToID   int    `json:"assigned_to_id,omitempty"`
	ParentIssueID  int64  `json:"parent_issue_id,omitempty"`
	Private        bool   `json:"is_private,omitempty"`
	EstimatedHours int    `json:"estimated_hours,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

type issue struct {
	ID                  int64          `json:"id,omitempty"`
	Project             util.IdName    `json:"project,omitempty"`
	Tracker             util.IdName    `json:"tracker,omitempty"`
	Status              issueStatus    `json:"status,omitempty"`
	Priority            util.IdName    `json:"priority,omitempty"`
	Author              util.IdName    `json:"author,omitempty"`
	AssignedTo          util.IdName    `json:"assigned_to,omitempty"`
	FixedVersion        util.IdName    `json:"fixed_version,omitempty"`
	Subject             string         `json:"subject,omitempty"`
	Description         string         `json:"description,omitempty"`
	StartDate           string         `json:"start_date,omitempty"`
	DueDate             string         `json:"due_date,omitempty"`
	DoneRatio           int            `json:"done_ratio,omitempty"`
	IsPrivate           bool           `json:"is_private,omitempty"`
	EstimatedHours      float32        `json:"estimated_hours,omitempty"`
	TotalEstimatedHours float32        `json:"total_estimated_hours,omitempty"`
	SpentHours          float32        `json:"spent_hours,omitempty"`
	TotalSpentHours     float32        `json:"total_spent_hours,omitempty"`
	CreatedOn           string         `json:"created_on,omitempty"`
	UpdatedOn           string         `json:"updated_on,omitempty"`
	Journals            []issueJournal `json:"journals,omitempty"`
	AllowedStatuses     []issueStatus  `json:"allowed_statuses,omitempty"`
	Notes               string         `json:"notes,omitempty"`
}

type issueJournal struct {
	ID           int64                 `json:"id,omitempty"`
	User         util.IdName           `json:"user,omitempty"`
	Notes        string                `json:"notes,omitempty"`
	CreatedOn    string                `json:"created_on,omitempty"`
	PrivateNotes bool                  `json:"private_notes,omitempty"`
	Details      []issueJournalDetails `json:"details,omitempty"`
}

type issueJournalDetails struct {
	Property string `json:"property,omitempty"`
	Name     string `json:"name,omitempty"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type issueStatus struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
}
