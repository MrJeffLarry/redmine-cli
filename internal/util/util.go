package util

import (
	"fmt"
	"strconv"
	"time"
)

type IdName struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Id struct {
	ID int `json:"id,omitempty"`
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

func TimeAgo(data string) string {
	date, error := time.Parse(time.RFC3339, data)
	if error != nil {
		return data
	}

	timeDiff := time.Now().Unix() - date.Unix()

	if timeDiff > (365 * 24 * 60 * 60) {
		return date.Format("2006-01-02")
	} else if timeDiff > (24 * 60 * 60) {
		return date.Format("02 Jan")
	} else if timeDiff > (60 * 60) {
		return fmt.Sprint((timeDiff / 60 / 60)) + " hours ago"
	} else if timeDiff > 60 {
		return fmt.Sprint(timeDiff/60) + " min ago"
	} else if timeDiff < 60 {
		return fmt.Sprint(timeDiff) + " sec ago"
	}

	return data
}
