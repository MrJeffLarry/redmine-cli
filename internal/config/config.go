package config

import (
	"errors"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	RED_CONFIG_REDMINE_URL     = "RED_CONFIG_REDMINE_URL"
	RED_CONFIG_REDMINE_API_KEY = "RED_CONFIG_REDMINE_API_KEY"
	RED_CONFIG_REDMINE_PROJECT = "RED_CONFIG_REDMINE_PROJECT"

	CONFIG_REDMINE_URL     = "server"
	CONFIG_REDMINE_API_KEY = "apiKey"
	CONFIG_REDMINE_PROJECT = "project"

	CONFIG_FILE   = "config.json"
	CONFIG_FOLDER = ".red"

	DEBUG_FLAG   = "debug"
	DEBUG_FLAG_S = "d"
)

type Red_t struct {
	RedmineURL     string
	RedmineApiKey  string
	RedmineProject string
	Debug          bool
}

//
//
//
//
//
//
func exEnv(name string, defValue string) string {
	if eValue, exName := os.LookupEnv(name); exName {
		return eValue
	} else {
		return defValue
	}
}

//
//
//
func (r *Red_t) IsConfigBad() bool {
	if len(r.RedmineURL) <= 0 {
		return true
	}
	if len(r.RedmineApiKey) <= 0 {
		return true
	}
	return false
}

//
//
//
func (r *Red_t) SetServer(server string) {
	r.RedmineURL = server
}

//
//
//
func (r *Red_t) SetApiKey(apiKey string) {
	r.RedmineApiKey = apiKey
}

//
//
//
func (r *Red_t) SetProject(id string) {
	r.RedmineProject = id
}

func createFolderPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			if os.IsPermission(err) {
				return errors.New("Could not create folder")
			}
			return err
		}
	}
	return nil
}

func ConfigPath() (string, error) {
	sep := string(os.PathSeparator)

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.New("Can't find home directory")
	}
	if err := createFolderPath(home + sep + CONFIG_FOLDER); err != nil {
		return "", errors.New("Could not create config folder")
	}

	return home + sep + CONFIG_FOLDER + sep, nil
}

func TmpPath() (string, error) {
	var path string
	var err error

	if path, err = ConfigPath(); err != nil {
		return "", err
	}

	if err = createFolderPath(path + "tmp"); err != nil {
		return "", err
	}

	return path + "tmp", nil
}

func CreateTmpFile(body string) (string, error) {
	var path string
	var err error

	if path, err = TmpPath(); err != nil {
		return "", err
	}

	f, err := os.CreateTemp(path, "tmp-*.md")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write([]byte(body)); err != nil {
		return "", err
	}
	return f.Name(), nil
}

//
//
//
func (r *Red_t) Save() error {
	var homePath string
	var err error

	if homePath, err = ConfigPath(); err != nil {
		return err
	}

	filePath := homePath + CONFIG_FILE

	//	viper.SetConfigName(CONFIG_FILE)
	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	viper.Set(CONFIG_REDMINE_URL, r.RedmineURL)
	viper.Set(CONFIG_REDMINE_API_KEY, r.RedmineApiKey)
	viper.Set(CONFIG_REDMINE_PROJECT, r.RedmineProject)

	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		if err := viper.SafeWriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

//
//
//
func (r *Red_t) LoadConfig() {
	sep := string(os.PathSeparator)

	home, err := homedir.Dir()
	if err != nil {
		return //errors.New("Can't find home directory")
	}

	filePath := home + sep + CONFIG_FOLDER + sep + CONFIG_FILE

	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return // errors.New("Can't read config file in root")
	}
	r.RedmineURL = viper.GetString(CONFIG_REDMINE_URL)
	r.RedmineApiKey = viper.GetString(CONFIG_REDMINE_API_KEY)
	r.RedmineProject = viper.GetString(CONFIG_REDMINE_PROJECT)
}

func (r *Red_t) localConfig() {
	var pwd string
	var err error
	var configPath string

	sep := string(os.PathSeparator)

	if pwd, err = os.Getwd(); err != nil {
		return
	}

	configPath = pwd + sep + CONFIG_FOLDER + sep + CONFIG_FILE

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if redmineURL := viper.GetString(CONFIG_REDMINE_URL); len(redmineURL) > 0 {
		r.RedmineURL = redmineURL
	}
	if redmineApiKey := viper.GetString(CONFIG_REDMINE_API_KEY); len(redmineApiKey) > 0 {
		r.RedmineApiKey = redmineApiKey
	}
	if redmineProject := viper.GetString(CONFIG_REDMINE_PROJECT); len(redmineProject) > 0 {
		r.RedmineProject = redmineProject
	}
}

//
//
//
func InitConfig() *Red_t {
	red := &Red_t{}

	red.RedmineURL = exEnv(RED_CONFIG_REDMINE_URL, "")
	red.RedmineApiKey = exEnv(RED_CONFIG_REDMINE_API_KEY, "")
	red.RedmineProject = exEnv(RED_CONFIG_REDMINE_PROJECT, "")

	red.LoadConfig()
	red.localConfig()

	return red
}
