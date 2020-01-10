package cmd

import (
	"fmt"
	"os"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/cobra"
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
	cobra.OnInitialize(func() { config.Setup() })
}
