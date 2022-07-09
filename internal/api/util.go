package api

import (
	"errors"
)

func StatusCode(status int) error {
	switch status {
	case 200:
		return nil
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
