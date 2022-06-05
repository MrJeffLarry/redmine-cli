package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
)

func ClientGET(r *config.Red_t, path string) ([]byte, int, error) {
	var err error
	var res []byte
	var statusCode int
	var req *http.Request
	var resp *http.Response

	statusCode = 0
	client := &http.Client{}

	req, err = http.NewRequest(http.MethodGet, r.RedmineURL+path, nil)
	if err != nil {
		fmt.Println(err.Error())
		return res, statusCode, errors.New("Servern hade problem och svara på vår förfrågan")
	}
	req.Header.Add("X-Redmine-API-Key", r.RedmineApiKey)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Servern hade problem och svara på vår förfrågan")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Kunde inte typa redmines svar")
	}

	return body, resp.StatusCode, nil
}

func ClientPOST(r *config.Red_t, path string) ([]byte, error) {
	return []byte{}, nil
}
