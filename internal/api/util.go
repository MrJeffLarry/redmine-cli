package api

import (
	"encoding/json"
	"errors"
)

type ResponseError struct {
	Errors []string `json:"errors"`
}

func ParseResponseError(body []byte) ResponseError {
	var response ResponseError

	if err := json.Unmarshal(body, &response); err != nil {
		response.Errors = append(response.Errors, "Could not parse response from server!")
	}

	return response
}

func StatusCode(status int) error {

	if status >= 200 && status <= 299 {
		return nil
	}

	switch status {
	case 401:
		return errors.New("Wrong login details, please try again")
	case 404:
		return errors.New("Does not exist")
	case 403:
		return errors.New("Does not have access to view or access this")
	case 500:
		return errors.New("Server has internal error, please try again")
	}

	return errors.New("Unknown error status code")
}
