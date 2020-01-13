package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	ApiKey      = "apikey"
	Hostname    = "hostname"
	DefaultsMap = "defaults"
	AliasesMap  = "aliases"
)

type DefaultsKey string

const (
	Activity DefaultsKey = "activity"
)

var AvailableDefaultsKeys = []string{string(Activity)}

func Setup() {
	home, err := homedir.Dir()
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

func Defaults() map[string]string {
	return viper.GetStringMapString(DefaultsMap)
}

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

func GetAliases() map[string]string {
	return viper.GetStringMapString(AliasesMap)
}

func GetAlias(key string) (value string, found bool) {
	aliases := GetAliases()
	value, found = aliases[key]
	return
}

func SetAlias(key string, value string) error {
	aliases := GetAliases()
	defer func() {
		viper.Set(AliasesMap, aliases)
		err := viper.WriteConfig()
		if err != nil {
			panic("unable to write config while adding new alias")
		}
	}()

	if value == "" { // Remove alias if value is empty
		delete(aliases, key)
		return nil
	}
	aliases[key] = value

	return nil
}
