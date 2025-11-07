package issue

import (
	"github.com/MrJeffLarry/redmine-cli/internal/cmd/global"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
)

const (
	FIELD_SUBJECT        = "Subject"
	FIELD_DESCRIPTION    = "Description"
	FIELD_CATEGORY       = "Category"
	FIELD_STATUS         = "Status"
	FIELD_PRIORITY       = "Priority"
	FIELD_TRACKER        = "Tracker"
	FIELD_NOTE           = "Notes"
	FIELD_TARGET_VERSION = "Target Version"
	FIELD_PARENT_ID      = "Parent ID"
	FIELD_ASSIGN         = "Assign"
	FIELD_PREVIEW        = "Preview"
	FIELD_SAVE           = "Save"
	FIELD_EXIT           = "Exit"
	FIELD_DEBUG          = "Print Debug"
)

type viewIssue struct {
	Issue issue `json:"issue,omitempty"`
}

type newIssueHolder struct {
	Issue newIssue `json:"issue,omitempty"`
}

type issues struct {
	Issues     []issue `json:"issues,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

type newIssue struct {
	ProjectID      int    `json:"project_id,omitempty"`
	TrackerID      int    `json:"tracker_id,omitempty"`
	StatusID       int    `json:"status_id,omitempty"`
	PriorityID     int    `json:"priority_id,omitempty"`
	Subject        string `json:"subject,omitempty"`
	Description    string `json:"description,omitempty"`
	CategoryID     int    `json:"category_id,omitempty"`
	FixedVersionID int    `json:"fixed_version_id,omitempty"`
	AssignedToID   int    `json:"assigned_to_id,omitempty"`
	ParentIssueID  int    `json:"parent_issue_id,omitempty"`
	Private        bool   `json:"is_private,omitempty"`
	EstimatedHours int    `json:"estimated_hours,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

type issue struct {
	ID                  int                  `json:"id,omitempty"`
	Project             util.IdName          `json:"project,omitempty"`
	Category             util.IdName          `json:"category,omitempty"`
	Tracker             util.IdName          `json:"tracker,omitempty"`
	Status              global.IssueStatus   `json:"status,omitempty"`
	Priority            util.IdName          `json:"priority,omitempty"`
	Author              util.IdName          `json:"author,omitempty"`
	AssignedTo          util.IdName          `json:"assigned_to,omitempty"`
	FixedVersion        util.IdName          `json:"fixed_version,omitempty"`
	Parent              util.Id              `json:"parent,omitempty"`
	Subject             string               `json:"subject,omitempty"`
	Description         string               `json:"description,omitempty"`
	StartDate           string               `json:"start_date,omitempty"`
	DueDate             string               `json:"due_date,omitempty"`
	DoneRatio           int                  `json:"done_ratio,omitempty"`
	IsPrivate           bool                 `json:"is_private,omitempty"`
	EstimatedHours      float32              `json:"estimated_hours,omitempty"`
	TotalEstimatedHours float32              `json:"total_estimated_hours,omitempty"`
	SpentHours          float32              `json:"spent_hours,omitempty"`
	TotalSpentHours     float32              `json:"total_spent_hours,omitempty"`
	CreatedOn           string               `json:"created_on,omitempty"`
	UpdatedOn           string               `json:"updated_on,omitempty"`
	Journals            []issueJournal       `json:"journals,omitempty"`
	AllowedStatuses     []global.IssueStatus `json:"allowed_statuses,omitempty"`
	Notes               string               `json:"notes,omitempty"`
}

type issueJournal struct {
	ID           int                   `json:"id,omitempty"`
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
