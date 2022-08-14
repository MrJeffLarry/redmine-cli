package api

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
)

const (
	ERR_CONN_CREATE  = "Could not create connection with server, correct server details?"
	ERR_CONN_SILENCE = "No response from server, correct server details?"
	ERR_CONN_RES     = "Could not read response from server.."
	ERR_CONN_AUTH    = "Request is not authorized, do you have access or have authenticated?"
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
		return res, statusCode, errors.New(ERR_CONN_CREATE)
	}

	req.Header.Add("X-Redmine-API-Key", r.RedmineApiKey)
	resp, err = client.Do(req)
	if err != nil {
		return res, statusCode, errors.New(ERR_CONN_SILENCE + " [" + r.RedmineURL + "]")
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, statusCode, errors.New(ERR_CONN_RES)
	}

	if statusCode == http.StatusForbidden || statusCode == http.StatusUnauthorized {
		return res, statusCode, errors.New(ERR_CONN_AUTH)
	}

	return res, statusCode, nil
}

func ClientPUT(r *config.Red_t, path string, body []byte) ([]byte, int, error) {
	var err error
	var res []byte
	var statusCode int
	var req *http.Request
	var resp *http.Response

	client := &http.Client{}
	statusCode = 0

	req, err = http.NewRequest(http.MethodPut, r.RedmineURL+path, bytes.NewReader(body))
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_CREATE)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Redmine-API-Key", r.RedmineApiKey)

	resp, err = client.Do(req)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_SILENCE + " [" + r.RedmineURL + "]")
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_RES)
	}

	if statusCode == http.StatusForbidden || statusCode == http.StatusUnauthorized {
		return res, statusCode, errors.New(ERR_CONN_AUTH)
	}

	return res, statusCode, nil
}

func ClientPOST(r *config.Red_t, path string, body []byte) ([]byte, int, error) {
	var err error
	var res []byte
	var statusCode int
	var req *http.Request
	var resp *http.Response

	client := &http.Client{}
	statusCode = 0

	req, err = http.NewRequest(http.MethodPost, r.RedmineURL+path, bytes.NewReader(body))
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_CREATE)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Redmine-API-Key", r.RedmineApiKey)

	resp, err = client.Do(req)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_SILENCE + " [" + r.RedmineURL + "]")
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_RES)
	}

	if statusCode == http.StatusForbidden || statusCode == http.StatusUnauthorized {
		return res, statusCode, errors.New(ERR_CONN_AUTH)
	}

	return res, statusCode, nil
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
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_CREATE)
	}
	req.SetBasicAuth(username, password)

	resp, err = client.Do(req)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_SILENCE + " [" + server + "]")
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_RES)
	}

	return res, statusCode, nil
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
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_CREATE)
	}

	req.Header.Add("X-Redmine-API-Key", apikey)

	resp, err = client.Do(req)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_SILENCE + " [" + server + "]")
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		print.Debug(r, err.Error())
		return res, statusCode, errors.New(ERR_CONN_RES)
	}

	return res, statusCode, nil
}
