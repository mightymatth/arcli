package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	hostname, username, password string
)

var authCmd = &cobra.Command{
	Use:     "auth",
	Aliases: []string{"connect"},
	Short:   "Authenticate to Redmine server",
	Long:    `Authenticate to Redmine server. Save Redmine API Key for further usage.`,
	Run:     auth(),
}

var authInteractive = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i", "in"},
	Short:   "Opens login session",
	PreRun:  handleInteractiveMode(),
	Run:     auth(),
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&hostname, "server", "s", "", "Hostname of Redmine server (e.g. host.redmine.org)")
	authCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	authCmd.Flags().StringVarP(&password, "password", "p", "", "Password")

	_ = authCmd.MarkFlagRequired("server")
	_ = authCmd.MarkFlagRequired("username")
	_ = authCmd.MarkFlagRequired("password")

	authCmd.AddCommand(authInteractive)
}

func auth() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Printf("hostname: %v, username: %v, password: %v\n", hostname, username, password)

		// TODO: Login with API
		// TODO: Set API key
		//viper.Set("apiKey", apiKey)
		//_ = viper.WriteConfig()
	}
}

func handleInteractiveMode() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var err error
		if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
			fmt.Printf("stdin/stdout should be terminal")
			return
		}

		oldState, err := terminal.MakeRaw(int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = terminal.Restore(int(os.Stdout.Fd()), oldState)
		}()

		screen := struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}
		t := terminal.NewTerminal(screen, "")

		hostname = AskForText(t, "Hostname: ", false)
		username = AskForText(t, "Username: ", false)
		password = AskForText(t, "Password: ", true)
	}
}

func AskForText(t *terminal.Terminal, prefix string, hidden bool) string {
	prompt := string(t.Escape.Cyan) + prefix + string(t.Escape.Reset)

	var line string
	var err error

	if hidden {
		line, err = t.ReadPassword(prompt)
	} else {
		t.SetPrompt(prompt)
		line, err = t.ReadLine()
	}

	if err == io.EOF {
		log.Fatalln("Exiting!")
	}

	return line
}
