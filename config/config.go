package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	// APIKey is the key of API key in config.
	APIKey = "apikey"
	// Host is the key of host in config.
	Host = "host"
	// DefaultsMap is the key of the defaults map in config.
	DefaultsMap = "defaults"
	// AliasesMap is the key of the aliases map in config.
	AliasesMap = "aliases"
	// UserID is the user ID
	UserID = "userID"
)

// DefaultsKey represents default key.
type DefaultsKey string

const (
	// Activity represents Redmine Activity.
	Activity DefaultsKey = "activity"
)

// AvailableDefaultsKeys stores all keys that are supported as defaults.
var AvailableDefaultsKeys = []string{string(Activity)}

// Setup setups permanent configuration in local storage.
func Setup() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.OpenFile(path.Join(home, ".arcli.yaml"), os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Cannot open/write configuration file", err)
	} else {
		_ = file.Close()
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".arcli")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		log.Fatal("Cannot read in configuration", err)
	}
}

// Defaults lists all defaults saved to permanent configuration.
func Defaults() map[string]string {
	return viper.GetStringMapString(DefaultsMap)
}

// SetDefault sets default value for given key.
func SetDefault(key DefaultsKey, value string) error {
	defaults := viper.GetStringMapString(DefaultsMap)
	defaults[string(key)] = value

	viper.Set(DefaultsMap, defaults)

	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("unable to write config while adding new default")
	}

	return nil
}

// GetAliases gets all stored aliases from permanent configuration.
func GetAliases() map[string]string {
	return viper.GetStringMapString(AliasesMap)
}

// GetAlias gets the alias from permanent configuration.
func GetAlias(key string) (value string, found bool) {
	aliases := GetAliases()
	value, found = aliases[key]
	return
}

// SetAlias sets the alias to permanent configuration.
func SetAlias(key string, value string) error {
	aliases := GetAliases()
	defer func() {
		viper.Set(AliasesMap, aliases)
		err := viper.WriteConfig()
		if err != nil {
			panic("unable to write config while adding new alias")
		}
	}()

	if value == "" {
		delete(aliases, key)
		return nil
	}
	aliases[key] = value

	return nil
}
