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
	"github.com/spf13/cobra"
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
	Name      string `json:"name"`
	Server    string `json:"server"`
	ApiKey    string `json:"api-key"`
	Project   string `json:"project"`
	ProjectID int    `json:"project-id"`
	UserID    int    `json:"user-id"`
}

type ConfigLocal_t struct {
	Server    string        `json:"server"`
	ApiKey    string        `json:"api-key"`
	Project   string        `json:"project"`
	ProjectID int           `json:"project-id"`
	UserID    int           `json:"user-id"`
	Editor    string        `json:"editor"`
	Pager     string        `json:"pager"`
	Issue     ConfigIssue_t `json:"issue"`
}

type ConfigV1_t struct {
	Server    string        `json:"server"`
	ApiKey    string        `json:"api-key"`
	Project   string        `json:"project"`
	ProjectID int           `json:"project-id"`
	UserID    int           `json:"user-id"`
	Editor    string        `json:"editor"`
	Pager     string        `json:"pager"`
	Issue     ConfigIssue_t `json:"issue"`
}

type ConfigV2_t struct {
	Version       string        `json:"version"`
	Servers       []Server_t    `json:"servers"`
	DefaultServer int           `json:"default-server"`
	Editor        string        `json:"editor"`
	Pager         string        `json:"pager"`
	Issue         ConfigIssue_t `json:"issue"`
}

type Red_t struct {
	Spinner     *spinner.Spinner
	Client      *http.Client
	Debug       bool          `json:"debug"`
	All         bool          `json:"all"`
	Config      ConfigV2_t    `json:"config"`
	LocalConfig ConfigLocal_t `json:"local-config"`
	Server      *Server_t     `json:"server"`
	Cmd         *cobra.Command
	Term        *terminal.Terminal
	Test        bool
}

func exEnv(name string, defValue string) string {
	if eValue, exName := os.LookupEnv(name); exName {
		return eValue
	} else {
		return defValue
	}
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

func configGlobalPath() (string, error) {
	sep := string(os.PathSeparator)

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Can't find home directory")
	}

	if err := createFolderPath(home + sep + CONFIG_FOLDER); err != nil {
		return "", err
	}

	return home + sep + CONFIG_FOLDER + sep + CONFIG_FILE, nil
}

func configLocalPath() (string, error) {
	var path string
	var pwd string
	var err error

	sep := string(os.PathSeparator)

	if pwd, err = os.Getwd(); err != nil {
		return path, err
	}

	path = pwd + sep + CONFIG_FOLDER

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := createFolderPath(path); err != nil {
			return "", err
		}
	}

	return path + sep + CONFIG_FILE, nil
}

func loadGlobalConfig() (ConfigV2_t, error) {
	var config ConfigV2_t

	filePath, err := configGlobalPath()
	if err != nil {
		return config, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, errors.New("can't read config file")
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, errors.New("can't unmarshal config file")
	}

	if config.Version == "2.0" {
		return config, nil
	}

	config_v1 := ConfigV1_t{}

	if err := json.Unmarshal(data, &config_v1); err != nil {
		return config, errors.New("can't unmarshal config file")
	}

	// Migrate v1 to v2
	config.Version = "2.0"
	config.Servers = []Server_t{}
	config.Servers = append(config.Servers, Server_t{
		Name:      "default",
		Server:    config_v1.Server,
		ApiKey:    config_v1.ApiKey,
		Project:   config_v1.Project,
		ProjectID: config_v1.ProjectID,
		UserID:    config_v1.UserID,
	})
	config.DefaultServer = 0
	config.Editor = config_v1.Editor
	config.Pager = config_v1.Pager

	err = saveGlobalConfig(config)
	if err != nil {
		fmt.Println(err)
		return config, nil
	}

	return config, nil
}

func saveGlobalConfig(config ConfigV2_t) error {

	filePath, err := configGlobalPath()
	if err != nil {
		return err
	}

	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadLocalConfig() (ConfigLocal_t, error) {
	var config ConfigLocal_t
	var err error

	filePath, err := configLocalPath()
	if err != nil {
		return config, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, errors.New("can't read config file")
	}

	if err := json.Unmarshal(data, &config); err == nil {
		return config, nil
	}

	return config, nil
}

func saveLocalConfig(config ConfigLocal_t) error {

	filePath, err := configLocalPath()
	if err != nil {
		return err
	}

	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}
	return nil
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

func (r *Red_t) IsConfigBad() bool {
	// Old or empty config file
	if r.Config.Version == "" {
		return true
	}
	if len(r.Config.Servers) == 0 {
		return true
	}
	return false
}

func (r *Red_t) AddServer(name, server, apiKey, project string, projectID, userID int) error {
	r.Config.Servers = append(r.Config.Servers, Server_t{
		Name:      name,
		Server:    server,
		ApiKey:    apiKey,
		Project:   project,
		ProjectID: projectID,
		UserID:    userID,
	})
	// Set newly added server as default/active server
	r.Config.DefaultServer = len(r.Config.Servers) - 1
	r.Server = &r.Config.Servers[r.Config.DefaultServer]
	return nil
}

func (r *Red_t) RemoveServer(id int, name string) error {
	for i, server := range r.Config.Servers {
		if server.Name == name || i == id {
			r.Config.Servers = append(r.Config.Servers[:i], r.Config.Servers[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *Red_t) RemoveCurrentServer() error {
	id := r.Config.DefaultServer
	r.Config.Servers = append(r.Config.Servers[:id], r.Config.Servers[id+1:]...)
	if len(r.Config.Servers) == 0 {
		r.Config.DefaultServer = 0
		r.Server = nil
	} else {
		r.Config.DefaultServer = 0
		r.Server = &r.Config.Servers[0]
	}
	return nil
}

func (r *Red_t) SetDefaultServer(id int) error {
	if id < 0 || id >= len(r.Config.Servers) {
		return errors.New("Redmine Server ID does not exist")
	}
	r.Config.DefaultServer = id
	r.Server = &r.Config.Servers[id]
	return nil
}

func (r *Red_t) GetServers() []Server_t {
	return r.Config.Servers
}

func (r *Red_t) SetProject(id string) error {
	r.Config.Servers[r.Config.DefaultServer].Project = id
	return nil
}

func (r *Red_t) SetProjectID(id int) error {
	r.Config.Servers[r.Config.DefaultServer].ProjectID = id
	return nil
}

func (r *Red_t) SetEditor(v string) error {
	r.Config.Editor = v
	return nil
}

func (r *Red_t) SetPager(v string) error {
	r.Config.Pager = v
	return nil
}

func (r *Red_t) ClearAll() error {
	r.Config = ConfigV2_t{}
	r.LocalConfig = ConfigLocal_t{}
	return nil
}

func (r *Red_t) Save() error {
	var err error

	if r.Test {
		return nil
	}

	err = saveGlobalConfig(r.Config)
	if err != nil {
		return err
	}

	return nil
}

func (r *Red_t) SaveLocalProject(projectID int) error {
	r.LocalConfig.ProjectID = projectID
	err := saveLocalConfig(r.LocalConfig)
	return err
}

func InitConfig() *Red_t {
	red := &Red_t{}
	var err error

	red.Client = &http.Client{}
	red.Spinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	server := &Server_t{}

	// Load from environment variables
	server.Server = exEnv(RED_CONFIG_REDMINE_URL, "")
	server.ApiKey = exEnv(RED_CONFIG_REDMINE_API_KEY, "")
	server.Project = exEnv(RED_CONFIG_REDMINE_PROJECT, "")
	server.ProjectID, _ = strconv.Atoi(exEnv(RED_CONFIG_REDMINE_PROJECT_ID, ""))
	server.UserID, _ = strconv.Atoi(exEnv(RED_CONFIG_REDMINE_USER_ID, ""))
	red.Config.Editor = exEnv(RED_CONFIG_EDITOR, "")
	red.Config.Pager = exEnv(RED_CONFIG_PAGER, "")

	c, err := loadGlobalConfig()
	if err != nil {
		return red
	}
	red.Config = c

	// If no servers configured, return empty config
	if len(red.Config.Servers) == 0 {
		return red
	}
	red.Server = &red.Config.Servers[red.Config.DefaultServer]

	red.LocalConfig, err = loadLocalConfig()
	if err != nil {
		return red
	}

	// Override with local config
	if len(red.LocalConfig.Server) > 0 {
		red.Server.Server = red.LocalConfig.Server
	}
	if len(red.LocalConfig.ApiKey) > 0 {
		red.Server.ApiKey = red.LocalConfig.ApiKey
	}
	if len(red.LocalConfig.Project) > 0 {
		red.Server.Project = red.LocalConfig.Project
	}
	if red.LocalConfig.ProjectID > 0 {
		red.Server.ProjectID = red.LocalConfig.ProjectID
	}
	if red.LocalConfig.UserID > 0 {
		red.Server.UserID = red.LocalConfig.UserID
	}
	if len(red.LocalConfig.Editor) > 0 {
		red.Config.Editor = red.LocalConfig.Editor
	}
	if len(red.LocalConfig.Pager) > 0 {
		red.Config.Pager = red.LocalConfig.Pager
	}

	return red
}
