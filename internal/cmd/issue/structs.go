package issue

type viewIssue struct {
	Issue issue `json:"issue,omitempty"`
}

type issues struct {
	Issues     []issue `json:"issues,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

type issue struct {
	ID                  int64          `json:"id,omitempty"`
	Project             issueIdName    `json:"project,omitempty"`
	Tracker             issueIdName    `json:"tracker,omitempty"`
	Status              issueStatus    `json:"status,omitempty"`
	Priority            issueIdName    `json:"priority,omitempty"`
	Author              issueIdName    `json:"author,omitempty"`
	AssignedTo          issueIdName    `json:"assigned_to,omitempty"`
	FixedVersion        issueIdName    `json:"fixed_version,omitempty"`
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
	User         issueIdName           `json:"user,omitempty"`
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

type issueIdName struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type issueStatus struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
}
