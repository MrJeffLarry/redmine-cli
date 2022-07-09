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
		return trackers, errors.New("Could not get trackers from project..")
	}

	if err := json.Unmarshal(body, &project); err != nil {
		print.Debug(r, err.Error())
		return trackers, errors.New("Could not parse and read response from server")
	}

	return project.Project.Trackers, nil
}
