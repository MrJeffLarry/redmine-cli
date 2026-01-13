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
	var idNames []util.IdName

	body, status, err := api.ClientGET(r, "/projects/"+strconv.Itoa(projectID)+"/versions.json")

	print.Debug(r, "%d %s", status, string(body))

	if err != nil {
		return idNames, err
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		print.Debug(r, err.Error())
		return idNames, errors.New("Could not parse and read response from server")
	}

	for _, v := range payload.Versions {
		idNames = append(idNames, util.IdName{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	idNames = append(idNames, util.IdName{
		ID:   -1,
		Name: "--None--",
	})

	return idNames, nil
}

func GetAssigns(r *config.Red_t, projectID int) ([]util.IdName, error) {
	var payload memberships
	var idNames []util.IdName

	body, status, err := api.ClientGET(r, "/projects/"+strconv.Itoa(projectID)+"/memberships.json")

	print.Debug(r, "%d %s", status, string(body))

	if err != nil {
		return idNames, err
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		print.Debug(r, err.Error())
		return idNames, errors.New("Could not parse and read response from server")
	}

	if r.Server.UserID > 0 {
		idNames = append(idNames, util.IdName{
			ID:   r.Server.UserID,
			Name: "Me",
		})
	}

	for _, v := range payload.Memberships {
		if v.User.ID <= 0 || v.User.ID == r.Server.UserID { // Skip group or me
			continue
		}

		idNames = append(idNames, util.IdName{
			ID:   v.User.ID,
			Name: v.User.Name,
		})
	}

	idNames = append(idNames, util.IdName{
		ID:   -1,
		Name: "--None--",
	})

	return idNames, nil
}
