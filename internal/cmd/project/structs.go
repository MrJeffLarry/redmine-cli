package project

type projects struct {
	Projects   []project `json:"projects,omitempty"`
	TotalCount int       `json:"total_count,omitempty"`
	Offset     int       `json:"offset,omitempty"`
	Limit      int       `json:"limit,omitempty"`
}

type project struct {
	ID              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Identifier      string `json:"identifier,omitempty"`
	Description     string `json:"description,omitempty"`
	IsPublic        bool   `json:"is_public,omitempty"`
	CreatedOn       string `json:"created_on,omitempty"`
	UpdatedOn       string `json:"updated_on,omitempty"`
	HomePage        string `json:"homepage,omitempty"`
	Status          int    `json:"status,omitempty"`
	Parent          idName `json:"parent,omitempty"`
	DefaultVersion  idName `json:"default_version,omitempty"`
	DefaultAssignee idName `json:"default_assignee,omitempty"`
}

type idName struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
