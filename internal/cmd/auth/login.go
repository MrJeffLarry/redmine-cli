package auth

import (
	"encoding/json"
	"fmt"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/terminal"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func loginCheckStatus(status int, badAuth string) bool {
	switch status {
	case 401:
		fmt.Println(badAuth)
		return false
	case 500:
		fmt.Println("Server had an internal error on our request, try again")
		return false
	case 200:
		return true
	default:
		fmt.Println("Could not process your request, please try again")
		return false
	}
}

func loginApiKey(r *config.Red_t, cmd *cobra.Command, server, apikey string) {
	var err error
	var res []byte
	var status int
	user := user{}

	if res, status, err = api.ClientAuthApiKeyGET(r, "/users/current.json", server, apikey); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.PrintDebug(r, status, string(res))

	if !loginCheckStatus(status, "Incorrect apikey, get (API access key) from: "+server+"/my/account") {
		return
	}

	if err := json.Unmarshal(res, &user); err != nil {
		print.Debug(r, status, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	r.SetApiKey(user.User.ApiKey)
	r.SetServer(server)
	r.SetUserID(user.User.ID)
	if err = r.Save(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Login done!")
}

func loginPassword(r *config.Red_t, cmd *cobra.Command, server, username string) {
	var bytePassword []byte
	var user user
	var err error

	fmt.Print("Password: ")
	if bytePassword, err = term.ReadPassword(0); err != nil {
		fmt.Println(err)
		return
	}

	res, status, err := api.ClientAuthBasicGET(r, "/users/current.json", server, username, string(bytePassword))
	if err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return
	}

	print.PrintDebug(r, status, string(res))

	fmt.Println("")

	if !loginCheckStatus(status, "Wrong username or password") {
		return
	}

	if err := json.Unmarshal(res, &user); err != nil {
		print.Debug(r, status, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return
	}

	r.SetApiKey(user.User.ApiKey)
	r.SetServer(server)
	r.SetUserID(user.User.ID)
	if err = r.Save(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Login success!")
}

func displayLogin(r *config.Red_t, cmd *cobra.Command) {
	var server string
	var username string
	var apikey string

	server = terminal.WriteLineReq("Server", 2)
	_, chooseID := terminal.ChooseString("Login method", []string{"Apikey", "Username and Password"})

	if chooseID == 0 {
		apikey = terminal.WriteLineReq("ApiKey", 2)
	} else {
		username = terminal.WriteLineReq("Username", 2)
	}

	if apikey != "" {
		loginApiKey(r, cmd, server, apikey)
		return
	}

	if username != "" {
		loginPassword(r, cmd, server, username)
		return
	}
}

func cmdAuthLogin(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to Redmine",
		Long:  "Authenticate to Redmine server",
		Run: func(cmd *cobra.Command, args []string) {
			displayLogin(r, cmd)
		},
	}
	return cmd
}
