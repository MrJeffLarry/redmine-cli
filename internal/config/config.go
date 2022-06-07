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

	CONFIG_REDMINE_URL     = "server"
	CONFIG_REDMINE_API_KEY = "apiKey"
	CONFIG_FILE            = "config.json"
	CONFIG_FOLDER          = ".red"

	DEBUG_FLAG = "debug"
)

type Red_t struct {
	RedmineURL    string
	RedmineApiKey string
	Debug         bool
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
func (r *Red_t) Save() error {
	sep := string(os.PathSeparator)

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return errors.New("Can't find home directory")
	}

	if _, err := os.Stat(home + sep + CONFIG_FOLDER); os.IsNotExist(err) {
		if err := os.Mkdir(home+sep+CONFIG_FOLDER, os.ModePerm); err != nil {
			if os.IsPermission(err) {
				fmt.Println("Could not create config folder", err)
			}
			fmt.Println(err)
			return err
		}
	}
	filePath := home + sep + CONFIG_FOLDER + sep + CONFIG_FILE

	//	viper.SetConfigName(CONFIG_FILE)
	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	viper.Set(CONFIG_REDMINE_URL, r.RedmineURL)
	viper.Set(CONFIG_REDMINE_API_KEY, r.RedmineApiKey)

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
func InitConfig() *Red_t {
	red := &Red_t{}

	red.RedmineURL = exEnv(RED_CONFIG_REDMINE_URL, "")
	red.RedmineApiKey = exEnv(RED_CONFIG_REDMINE_API_KEY, "")

	red.LoadConfig()

	return red
}
