package project

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
)

func GetTrackers(r *config.Red_t, projectID int) ([]util.IdName, error) {
	var trackers []util.IdName
	var project singleProject

	body, status, err := api.ClientGET(r, "/projects/"+strconv.Itoa(projectID)+".json?include=trackers,issue_categories")

	print.Debug(r, "%d %s", status, string(body))

	if err != nil {
		return trackers, err
	}

	if err := json.Unmarshal(body, &project); err != nil {
		print.Debug(r, err.Error())
		return trackers, errors.New("Could not parse and read response from server")
	}

	return project.Project.Trackers, nil
}

func GetVersions(r *config.Red_t, projectID int) ([]util.IdName, error) {
	var payload versions
	var versions []util.IdName

	body, status, err := api.ClientGET(r, "/projects/"+strconv.Itoa(projectID)+"/versions.json")

	print.Debug(r, "%d %s", status, string(body))

	if err != nil {
		return versions, err
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		print.Debug(r, err.Error())
		return versions, errors.New("Could not parse and read response from server")
	}

	for _, v := range payload.Versions {
		versions = append(versions, util.IdName{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	return versions, nil
}
