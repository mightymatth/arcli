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
	ApiKey   = "apikey"
	Hostname = "hostname"
)

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
