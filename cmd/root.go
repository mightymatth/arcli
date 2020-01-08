package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "arcli",
	Short:   "Awesome Redmine CLI",
	Long:    `Client for Redmine. Wrapper around Redmine API`,
	Version: "v0.0.0 (Redmine API v3.3)",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, _ = os.Create(path.Join(home, ".arcli.yaml"))
	viper.AddConfigPath(home)
	viper.SetConfigName(".arcli")

	viper.AutomaticEnv()
}
