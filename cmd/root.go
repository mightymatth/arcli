package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mightymatth/arcli/client"
	"github.com/mightymatth/arcli/config"
	"github.com/spf13/cobra"
)

type Version struct {
	Version, RedmineAPIVersion string
}

func (v Version) String() string {
	return fmt.Sprintf("v%v (Redmine API v%v)", v.Version, v.RedmineAPIVersion)
}

var (
	VERSION     Version
	versionFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "arcli",
	Short: "Awesome Redmine CLI",
	Long:  `Awesome Redmine CLI. Wrapper around Redmine API`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Println(VERSION)
		} else {
			_ = cmd.Help()
		}
	},
}

var RClient *client.Client

func Execute(version string) {
	VERSION = Version{
		Version:           version,
		RedmineAPIVersion: "3.3",
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false,
		"Current arcli and supported Redmine API version")

	cobra.OnInitialize(func() { config.Setup() })

	RClient = &client.Client{
		HTTPClient: &http.Client{},
		UserAgent:  "arcli",
	}
}
