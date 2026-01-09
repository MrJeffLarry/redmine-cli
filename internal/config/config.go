package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MrJeffLarry/redmine-cli/internal/terminal"
	"github.com/briandowns/spinner"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	RED_CONFIG_REDMINE_URL        = "RED_CONFIG_REDMINE_URL"
	RED_CONFIG_REDMINE_API_KEY    = "RED_CONFIG_REDMINE_API_KEY"
	RED_CONFIG_REDMINE_PROJECT    = "RED_CONFIG_REDMINE_PROJECT"
	RED_CONFIG_REDMINE_PROJECT_ID = "RED_CONFIG_REDMINE_PROJECT_ID"
	RED_CONFIG_REDMINE_USER_ID    = "RED_CONFIG_REDMINE_USER_ID"
	RED_CONFIG_EDITOR             = "RED_CONFIG_EDITOR"
	RED_CONFIG_PAGER              = "RED_CONFIG_PAGER"
	RED_CONFIG_ISSUE_VIEW_JOURNAL = "RED_CONFIG_ISSUE_VIEW_JOURNAL"

	CONFIG_SERVERS            = "servers"
	CONFIG_SERVER_ID          = "id"
	CONFIG_SERVER_NAME        = "name"
	CONFIG_REDMINE_URL        = "server"
	CONFIG_REDMINE_API_KEY    = "api-key"
	CONFIG_REDMINE_PROJECT    = "project"
	CONFIG_REDMINE_PROJECT_ID = "project-id"
	CONFIG_REDMINE_USER_ID    = "user-id"
	CONFIG_EDITOR             = "editor"
	CONFIG_PAGER              = "pager"
	CONFIG_ISSUE              = "issue"

	CONFIG_FILE   = "config.json"
	CONFIG_FOLDER = ".red"

	DEBUG_FLAG   = "debug"
	DEBUG_FLAG_S = "d"

	ALL_FLAG = "all"
	RID_FLAG = "rid"
)

type ConfigIssue_t struct {
	ViewJournalAlways bool `json:"view-journal"`
}

type Server_t struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Server    string `json:"server"`
	ApiKey    string `mapstructure:"api-key"`
	Project   string `json:"project"`
	ProjectID int    `mapstructure:"project-id"`
	UserID    int    `mapstructure:"user-id"`
}

type Config_t struct {
	Server    string        `json:"server"`
	ApiKey    string        `mapstructure:"api-key"`
	Project   string        `json:"project"`
	ProjectID int           `mapstructure:"project-id"`
	UserID    int           `mapstructure:"user-id"`
	Editor    string        `json:"editor"`
	Pager     string        `json:"pager"`
	Issue     ConfigIssue_t `json:"issue"`
}

type ConfigV2_t struct {
	Version       string        `json:"version"`
	Servers       []Server_t    `json:"servers"`
	DefaultServer int           `json:"default-server`
	Editor        string        `json:"editor"`
	Pager         string        `json:"pager"`
	Issue         ConfigIssue_t `json:"issue"`
}

type Red_t struct {
	Spinner       *spinner.Spinner
	Client        *http.Client
	Debug         bool     `json:"debug"`
	All           bool     `json:"all"`
	ConfigVersion string   `json:"version"`
	Server        Server_t `json:"server"`
	SaveConfig    bool     `json:"save-config"`
	Editor        string   `json:"editor"`
	Pager         string   `json:"pager"`
	Cmd           *cobra.Command
	Term          *terminal.Terminal
	Test          bool
}

func exEnv(name string, defValue string) string {
	if eValue, exName := os.LookupEnv(name); exName {
		return eValue
	} else {
		return defValue
	}
}

func (r *Red_t) IsConfigBad() bool {
	// Old or empty config file
	if r.ConfigVersion == "" {
		return true
	}
	if r.Server.Server == "" || r.Server.ApiKey == "" {
		return true
	}
	return false
}

func (r *Red_t) AddServer(s Server_t) {

}

func (r *Red_t) RemoveServer(id int) {

}

func (r *Red_t) SetServer(server string) {
	r.Server.Server = server
	r.SaveConfig = true
}

func (r *Red_t) SetApiKey(apiKey string) {
	r.Server.ApiKey = apiKey
	r.SaveConfig = true
}

func (r *Red_t) SetProject(id string) {
	r.Server.Project = id
	r.SaveConfig = true
}

func (r *Red_t) SetProjectID(id int) {
	r.Server.ProjectID = id
	r.SaveConfig = true
}

func (r *Red_t) SetUserID(id int) {
	r.Server.UserID = id
	r.SaveConfig = true
}

func (r *Red_t) SetEditor(v string) {
	r.Editor = v
	r.SaveConfig = true
}

func (r *Red_t) SetPager(v string) {
	r.Pager = v
	r.SaveConfig = true
}

func (r *Red_t) ClearAll() {
	r.Server = Server_t{}
	r.SaveConfig = true
}

func createFolderPath(path string) error {
	_, err := os.Stat(path)
	if os.IsPermission(err) {
		return errors.New("We are not allowed to access folder in " + path + " please check permissions")
	}

	if os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			if os.IsPermission(err) {
				return errors.New("We are not allowed to create folder in " + path)
			}
			return err
		}
	}
	return nil
}

func ConfigLocalPath() (string, error) {
	var configPath string
	var pwd string
	var err error

	sep := string(os.PathSeparator)

	if pwd, err = os.Getwd(); err != nil {
		return configPath, err
	}

	configPath = pwd + sep + CONFIG_FOLDER

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createFolderPath(configPath); err != nil {
			return "", err
		}
		return configPath + sep, nil
	}

	return configPath + sep, nil
}

func ConfigPath() (string, error) {
	sep := string(os.PathSeparator)

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.New("Can't find home directory")
	}

	if err := createFolderPath(home + sep + CONFIG_FOLDER); err != nil {
		return "", err
	}

	return home + sep + CONFIG_FOLDER + sep, nil
}

func CreateTmpFile(body string) (string, error) {
	var err error

	f, err := os.CreateTemp(os.TempDir(), "red-*.md")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write([]byte(body)); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func saveLocal(r *Red_t, name string, value interface{}) error {
	var err error
	var configPath string

	if configPath, err = ConfigLocalPath(); err != nil {
		return err
	}

	viper.Reset()
	viper.SetConfigFile(configPath + CONFIG_FILE)
	viper.SetConfigType("json")

	viper.ReadInConfig() // ignore if we can or not read we will try write in any way

	viper.Set(name, value)

	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		if err := viper.SafeWriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Red_t) Save() error {
	var homePath string
	var err error

	if homePath, err = ConfigPath(); err != nil {
		return err
	}

	filePath := homePath + CONFIG_FILE

	config := ConfigV2_t{}
	config.Version = "2.0"
	config.Servers = []Server_t{}
	config.Servers = append(config.Servers, r.Server)
	config.DefaultServer = r.Server.Id
	config.Editor = r.Editor
	config.Pager = r.Pager

	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if r.Test {
		return nil
	}

	err = os.WriteFile(filePath, body, 0644)
	if err == nil || r.Test {
		return nil
	}

	// viper.SetConfigFile(filePath)
	// viper.SetConfigType("json")

	// viper.Set(CONFIG_REDMINE_URL, r.Config.Server)
	// viper.Set(CONFIG_REDMINE_API_KEY, r.Config.ApiKey)
	// viper.Set(CONFIG_REDMINE_PROJECT, r.Config.Project)
	// viper.Set(CONFIG_REDMINE_PROJECT_ID, r.Config.ProjectID)
	// viper.Set(CONFIG_REDMINE_USER_ID, r.Config.UserID)
	// viper.Set(CONFIG_EDITOR, r.Config.Editor)
	// viper.Set(CONFIG_PAGER, r.Config.Pager)

	// if err := viper.WriteConfig(); err != nil {
	// 	if err := viper.SafeWriteConfig(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (r *Red_t) SaveLocalProject(projectID int) error {
	return saveLocal(r, CONFIG_REDMINE_PROJECT_ID, projectID)
}

func (r *Red_t) LoadConfig() error {
	sep := string(os.PathSeparator)

	home, err := homedir.Dir()
	if err != nil {
		return errors.New("can't find home directory")
	}

	filePath := home + sep + CONFIG_FOLDER + sep + CONFIG_FILE

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	config_v2 := ConfigV2_t{}

	if err := json.Unmarshal(data, &config_v2); err == nil {
		r.ConfigVersion = config_v2.Version
		for _, server := range config_v2.Servers {
			if server.Id == config_v2.DefaultServer {
				r.Server = server
				break
			}
		}
		r.Editor = config_v2.Editor
		r.Pager = config_v2.Pager

		// Valid config v2
		if !r.IsConfigBad() {
			return nil
		}
	}

	config_v1 := Config_t{}

	if err := json.Unmarshal(data, &config_v1); err != nil {
		return errors.New("can't unmarshal config file")
	}

	r.Server = Server_t{
		Id:        1,
		Name:      "default",
		Server:    config_v1.Server,
		ApiKey:    config_v1.ApiKey,
		Project:   config_v1.Project,
		ProjectID: config_v1.ProjectID,
		UserID:    config_v1.UserID,
	}
	r.Editor = config_v1.Editor
	r.Pager = config_v1.Pager

	return nil
}

func (r *Red_t) localConfig() error {
	var pwd string
	var err error
	var configPath string

	sep := string(os.PathSeparator)

	if pwd, err = os.Getwd(); err != nil {
		return errors.New("can't find current directory")
	}

	configPath = pwd + sep + CONFIG_FOLDER + sep + CONFIG_FILE

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err = viper.ReadInConfig(); err != nil {
		return errors.New("can't read config file in local")
	}

	var c Config_t
	if err := viper.Unmarshal(&c); err != nil {
		return errors.New("can't unmarshal config file")
	}

	if len(c.Server) > 0 {
		r.Config.Server = c.Server
	}
	if len(c.ApiKey) > 0 {
		r.Config.ApiKey = c.ApiKey
	}
	if len(c.Project) > 0 {
		r.Config.Project = c.Project
	}
	if c.ProjectID > 0 {
		r.Config.ProjectID = c.ProjectID
	}
	if c.UserID > 0 {
		r.Config.UserID = c.UserID
	}
	return nil
}

func InitConfig() *Red_t {
	red := &Red_t{}

	red.Client = &http.Client{}
	red.Spinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	red.Config.Server = exEnv(RED_CONFIG_REDMINE_URL, "")
	red.Config.ApiKey = exEnv(RED_CONFIG_REDMINE_API_KEY, "")
	red.Config.Project = exEnv(RED_CONFIG_REDMINE_PROJECT, "")
	red.Config.ProjectID, _ = strconv.Atoi(exEnv(RED_CONFIG_REDMINE_PROJECT_ID, ""))
	red.Config.UserID, _ = strconv.Atoi(exEnv(RED_CONFIG_REDMINE_USER_ID, ""))
	red.Config.Editor = exEnv(RED_CONFIG_EDITOR, "")
	red.Config.Pager = exEnv(RED_CONFIG_PAGER, "")

	if err := red.LoadConfig(); err != nil {
		fmt.Println(err)
	}
	if err := red.localConfig(); err != nil {
		fmt.Println(err)
	}

	return red
}
