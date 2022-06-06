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

	client := &http.Client{}
	statusCode = 0

	req, err = http.NewRequest(http.MethodGet, r.RedmineURL+path, nil)
	if err != nil {
		fmt.Println(err.Error())
		return res, statusCode, errors.New("Server could not handle our request")
	}
	req.Header.Add("X-Redmine-API-Key", r.RedmineApiKey)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Server could not handle our request")
	}
	defer resp.Body.Close()

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Kunde inte typa redmines svar")
	}

	return res, resp.StatusCode, nil
}

func ClientPOST(r *config.Red_t, path string) ([]byte, error) {
	return []byte{}, nil
}

func ClientAuthBasicGET(r *config.Red_t, path, server, username, password string) ([]byte, int, error) {
	var err error
	var res []byte
	var statusCode int
	var req *http.Request
	var resp *http.Response

	client := &http.Client{}
	statusCode = 0

	req, err = http.NewRequest(http.MethodGet, server+path, nil)
	if err != nil {
		fmt.Println(err.Error())
		return res, statusCode, errors.New("Server could not handle our request")
	}
	req.SetBasicAuth(username, password)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Server could not handle our request")
	}
	defer resp.Body.Close()

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Could not read response from server, please try again")
	}

	return res, resp.StatusCode, nil
}

func ClientAuthApiKeyGET(r *config.Red_t, path, server, apikey string) ([]byte, int, error) {
	var err error
	var res []byte
	var statusCode int
	var req *http.Request
	var resp *http.Response

	client := &http.Client{}
	statusCode = 0

	req, err = http.NewRequest(http.MethodGet, server+path, nil)
	if err != nil {
		fmt.Println(err.Error())
		return res, statusCode, errors.New("Server could not handle our request")
	}
	req.Header.Add("X-Redmine-API-Key", apikey)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Server could not handle our request")
	}
	defer resp.Body.Close()

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return res, resp.StatusCode, errors.New("Could not read response from server, please try again")
	}

	return res, resp.StatusCode, nil
}
