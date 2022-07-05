package project

import "github.com/MrJeffLarry/redmine-cli/internal/util"

type projects struct {
	Projects   []project `json:"projects,omitempty"`
	TotalCount int       `json:"total_count,omitempty"`
	Offset     int       `json:"offset,omitempty"`
	Limit      int       `json:"limit,omitempty"`
}

type singleProject struct {
	Project project `json:"project,omitempty"`
}

type project struct {
	ID              int64         `json:"id,omitempty"`
	Name            string        `json:"name,omitempty"`
	Identifier      string        `json:"identifier,omitempty"`
	Description     string        `json:"description,omitempty"`
	IsPublic        bool          `json:"is_public,omitempty"`
	CreatedOn       string        `json:"created_on,omitempty"`
	UpdatedOn       string        `json:"updated_on,omitempty"`
	HomePage        string        `json:"homepage,omitempty"`
	Status          int           `json:"status,omitempty"`
	Parent          util.IdName   `json:"parent,omitempty"`
	DefaultVersion  util.IdName   `json:"default_version,omitempty"`
	DefaultAssignee util.IdName   `json:"default_assignee,omitempty"`
	Trackers        []util.IdName `json:"trackers,omitempty"`
}
