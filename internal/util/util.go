package util

type IdName struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
