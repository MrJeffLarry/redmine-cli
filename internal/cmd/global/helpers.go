package global

import (
	"encoding/json"
	"errors"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
)

func GetPriorities(r *config.Red_t) ([]util.IdName, error) {
	var payload issuePriorities
	var idNames []util.IdName

	body, status, err := api.ClientGET(r, "/enumerations/issue_priorities.json")

	print.Debug(r, "%d %s", status, string(body))

	if err != nil {
		return idNames, err
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		print.Debug(r, err.Error())
		return idNames, errors.New("Could not parse and read response from server")
	}

	for _, v := range payload.IssuePriorities {
		idNames = append(idNames, util.IdName{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	return idNames, nil
}

func GetIssueStatus(r *config.Red_t) ([]util.IdName, error) {
	var payload issueStatusHolder
	var idNames []util.IdName

	body, status, err := api.ClientGET(r, "/issue_statuses.json")

	print.Debug(r, "%d %s", status, string(body))

	if err != nil || status != 200 {
		return idNames, errors.New("Could not get statuses from server, abort")
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		print.Debug(r, err.Error())
		return idNames, errors.New("Could not parse and read response from server")
	}

	for _, status := range payload.IssueStatus {
		idname := util.IdName{
			ID:   status.ID,
			Name: status.Name,
		}
		idNames = append(idNames, idname)
	}

	return idNames, nil
}
