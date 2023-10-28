package global

type issuePriorities struct {
	IssuePriorities []issuePrioritie `json:"issue_priorities"`
}

type issuePrioritie struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Active    bool   `json:"active"`
}
