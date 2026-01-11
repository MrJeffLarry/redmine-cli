package auth

import (
	"encoding/json"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

func loginApiKey(r *config.Red_t, cmd *cobra.Command, name, server, apikey string) bool {
	var err error
	var res []byte
	var status int
	user := user{}

	if res, status, err = api.ClientAuthApiKeyGET(r, "/users/current.json", server, apikey); err != nil {
		print.Error("StatusCode %d, %s", status, err.Error())
		return false
	}

	print.Debug(r, "%d %s", status, string(res))

	if err = api.StatusCode(status); err != nil {
		print.Error(err.Error())
		return false
	}

	if err := json.Unmarshal(res, &user); err != nil {
		print.Debug(r, err.Error())
		print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
		return false
	}

	r.AddServer(name, server, user.User.ApiKey, "", 0, user.User.ID)

	if err = r.Save(); err != nil {
		print.Error(err.Error())
		return false
	}

	print.OK("Login done!")
	return true
}

func loginPassword(r *config.Red_t, cmd *cobra.Command, name, server, username string) {
	var password string
	var user user
	var err error
	var retry bool

	for true {

		if retry {
			if !r.Term.Confirm("Try again?") {
				return
			}
		}
		retry = true

		if username, err = r.Term.PromptStringRequire("Username", ""); err != nil {
			print.Debug(r, err.Error())
			print.Error("Could not read input, please try again or submit issue")
			continue
		}

		if password, err = r.Term.PromptPassword("Password", ""); err != nil {
			print.Error(err.Error())
			continue
		}

		res, status, err := api.ClientAuthBasicGET(r, "/users/current.json", server, username, password)
		if err != nil {
			print.Error("StatusCode %d, %s", status, err.Error())
			continue
		}

		print.Debug(r, "%d %s", status, string(res))

		if err = api.StatusCode(status); err != nil {
			print.Info("If Two-factor authentication is enabled, this login method will not work as it is not supported, please use Apikey instead, you will find the Apikey (Api access key) at /my/account\n")
			print.Error(err.Error())
			continue
		}

		if err := json.Unmarshal(res, &user); err != nil {
			print.Debug(r, err.Error())
			print.Error("StatusCode %d, %s", status, "Could not parse and read response from server")
			return
		}
		break
	}

	r.AddServer(name, server, user.User.ApiKey, "", 0, user.User.ID)

	if err = r.Save(); err != nil {
		print.Error(err.Error())
		return
	}

	print.OK("Login success!")
}

func displayLogin(r *config.Red_t, cmd *cobra.Command) {
	var alias string
	var server string
	var username string
	var apikey string
	var err error

	print.Info(text.FgGreen.Sprint("Welcome to Red an Redmine CLI\n") +
		"Before login make sure you have enabled `Enable REST web service`\nfind it in Administration -> Settings -> API or use url /settings?tab=api\nYou find ApiKey (API access key) from /my/account\n\n")

	if alias, err = r.Term.PromptStringRequire("Redmine server alias (Used for multi redmine server support)", ""); err != nil {
		print.Debug(r, err.Error())
		print.Error("Could not read input, please try again or submit issue")
		return
	}

	if server, err = r.Term.PromptStringRequire("Server URL (https://example.com)", ""); err != nil {
		print.Debug(r, err.Error())
		print.Error("Could not read input, please try again or submit issue")
		return
	}

	_, chooseID := r.Term.ChooseString("Login method", []string{"Apikey", "Username and Password"})

	if chooseID == 0 {
		if apikey, err = r.Term.PromptPassword("ApiKey", ""); err != nil {
			print.Debug(r, err.Error())
			print.Error("Could not read input, please try again or submit issue")
			return
		}
		loginApiKey(r, cmd, alias, server, apikey)
		return
	}

	loginPassword(r, cmd, alias, server, username)
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
