package util

import "strconv"

type IdName struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Errors struct {
	Errors []string `json:"errors,omitempty"`
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CheckID(id string) bool {
	if _, err := strconv.Atoi(id); err != nil {
		return false
	}
	return true
}
