package config

import (
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

	CONFIG_REDMINE_URL        = "server"
	CONFIG_REDMINE_API_KEY    = "api-key"
	CONFIG_REDMINE_PROJECT    = "project"
	CONFIG_REDMINE_PROJECT_ID = "project-id"
	CONFIG_REDMINE_USER_ID    = "user-id"
	CONFIG_EDITOR             = "editor"
	CONFIG_PAGER              = "pager"
	CONFIG_ISSUE              = "issue"
	CONFIG_INSTANCES          = "instances"
	CONFIG_DEFAULT_INSTANCE   = "default-instance"

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

type MultiInstanceConfig_t struct {
	Instances       map[string]Config_t `json:"instances"`
	DefaultInstance string              `json:"default-instance"`
}

type Red_t struct {
	Spinner         *spinner.Spinner
	Client          *http.Client
	Debug           bool     `json:"debug"`
	All             bool     `json:"all"`
	Config          Config_t `json:"config"`
	Cmd             *cobra.Command
	Term            *terminal.Terminal
	Test            bool
	RID             string                `json:"rid"`
	MultiConfig     MultiInstanceConfig_t `json:"multi-config"`
	UseMultiMode    bool                  `json:"use-multi-mode"`
}

func exEnv(name string, defValue string) string {
	if eValue, exName := os.LookupEnv(name); exName {
		return eValue
	} else {
		return defValue
	}
}

func (r *Red_t) IsConfigBad() bool {
	if len(r.Config.Server) <= 0 {
		return true
	}
	if len(r.Config.ApiKey) <= 0 {
		return true
	}
	return false
}

func (r *Red_t) SetServer(server string) {
	r.Config.Server = server
}

func (r *Red_t) SetApiKey(apiKey string) {
	r.Config.ApiKey = apiKey
}

func (r *Red_t) SetProject(id string) {
	r.Config.Project = id
}

func (r *Red_t) SetProjectID(id int) {
	r.Config.ProjectID = id
}

func (r *Red_t) SetUserID(id int) {
	r.Config.UserID = id
}

func (r *Red_t) SetEditor(v string) {
	r.Config.Editor = v
}

func (r *Red_t) SetPager(v string) {
	r.Config.Pager = v
}

func (r *Red_t) SetRID(rid string) {
	r.RID = rid
	r.UseMultiMode = true
	
	// Initialize multi-config if needed
	if r.MultiConfig.Instances == nil {
		r.MultiConfig.Instances = make(map[string]Config_t)
	}
	
	// If no default instance set, make this one the default
	if r.MultiConfig.DefaultInstance == "" {
		r.MultiConfig.DefaultInstance = rid
	}
}

func (r *Red_t) ClearAll() {
	if r.UseMultiMode && r.RID != "" {
		// Clear only the current instance
		delete(r.MultiConfig.Instances, r.RID)
		
		// If clearing the default instance, reset it
		if r.MultiConfig.DefaultInstance == r.RID {
			// Find another instance to be default, or set to empty
			for key := range r.MultiConfig.Instances {
				r.MultiConfig.DefaultInstance = key
				break
			}
			if len(r.MultiConfig.Instances) == 0 {
				r.MultiConfig.DefaultInstance = ""
			}
		}
	}
	
	r.Config.Server = ""
	r.Config.ApiKey = ""
	r.Config.UserID = 0
	r.Config.Project = ""
	r.Config.ProjectID = 0
	r.Config.Editor = ""
	r.Config.Pager = ""
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

	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	if r.UseMultiMode {
		// Save in multi-instance mode
		r.MultiConfig.Instances[r.RID] = r.Config
		viper.Set(CONFIG_INSTANCES, r.MultiConfig.Instances)
		viper.Set(CONFIG_DEFAULT_INSTANCE, r.MultiConfig.DefaultInstance)
	} else {
		// Save in single-instance mode (backward compatible)
		viper.Set(CONFIG_REDMINE_URL, r.Config.Server)
		viper.Set(CONFIG_REDMINE_API_KEY, r.Config.ApiKey)
		viper.Set(CONFIG_REDMINE_PROJECT, r.Config.Project)
		viper.Set(CONFIG_REDMINE_PROJECT_ID, r.Config.ProjectID)
		viper.Set(CONFIG_REDMINE_USER_ID, r.Config.UserID)
		viper.Set(CONFIG_EDITOR, r.Config.Editor)
		viper.Set(CONFIG_PAGER, r.Config.Pager)
	}

	if r.Test {
		return nil
	}

	if err := viper.WriteConfig(); err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			return err
		}
	}
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

	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Check if multi-instance config exists
	if viper.IsSet(CONFIG_INSTANCES) {
		r.UseMultiMode = true
		r.MultiConfig.Instances = make(map[string]Config_t)
		if err := viper.UnmarshalKey(CONFIG_INSTANCES, &r.MultiConfig.Instances); err != nil {
			return errors.New("can't unmarshal multi-instance config")
		}
		r.MultiConfig.DefaultInstance = viper.GetString(CONFIG_DEFAULT_INSTANCE)
		
		// If no RID specified, use default or "1"
		if r.RID == "" {
			if r.MultiConfig.DefaultInstance != "" {
				r.RID = r.MultiConfig.DefaultInstance
			} else {
				r.RID = "1"
			}
		}
		
		// Load the specific instance config
		if cfg, ok := r.MultiConfig.Instances[r.RID]; ok {
			r.Config = cfg
		}
	} else {
		// Legacy single-instance mode
		r.UseMultiMode = false
		if err := viper.Unmarshal(&r.Config); err != nil {
			return errors.New("can't unmarshal config file")
		}
	}

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
	
	// Initialize multi-instance structures
	red.MultiConfig.Instances = make(map[string]Config_t)
	red.UseMultiMode = false
	red.RID = ""

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
