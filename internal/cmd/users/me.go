package users

import (
	"encoding/json"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/spf13/cobra"
)

type user struct {
	User userInfo `json:"user,omitempty"`
}

type userInfo struct {
	ID              int64  `json:"id,omitempty"`
	Login           string `json:"login,omitempty"`
	Admin           bool   `json:"admin,omitempty"`
	FirstName       string `json:"firstname,omitempty"`
	LastName        string `json:"lastname,omitempty"`
	CreatedOn       string `json:"created_on,omitempty"`
	UpdatedOn       string `json:"updated_on,omitempty"`
	LastLoginOn     string `json:"last_login_on,omitempty"`
	PasswdChangedOn string `json:"passwd_changed_on,omitempty"`
	AvatarUrl       string `json:"avatar_url,omitempty"`
	TwofaScheme     string `json:"twofa_scheme,omitempty"`
	ApiKey          string `json:"api_key,omitempty"`
}

func displayMeGET(r *config.Red_t, path string) {
	var err error
	var body []byte
	var status int
	user := user{}

	if body, status, err = api.ClientGET(r, path); err != nil {
		fmt.Println(status, "Could not get response from client", err)
		return
	}

	print.PrintDebug(r, status, string(body))

	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println(err)
		fmt.Println(status, "Could not parse and read response from server")
		return
	}

	fmt.Printf("#%d %s %s (%s)\n\nCreated: %s\nUpdated: %s\nLast Login: %s\n",
		user.User.ID,
		user.User.FirstName,
		user.User.LastName,
		user.User.Login,
		user.User.CreatedOn,
		user.User.UpdatedOn,
		user.User.LastLoginOn,
	)
}

func cmdUsersMe(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Display my info",
		Long:  "Display account and other information about my self",
		Run: func(cmd *cobra.Command, args []string) {
			displayMeGET(r, "/users/current.json")
		},
	}
	return cmd
}
