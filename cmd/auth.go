package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/viper"

	"github.com/mightymatth/arcli/client"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	hostname, username, password string
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"connect"},
	Short:   "Authenticate to Redmine server.",
	Long:    `Authenticate to Redmine server. Save credentials for further usage.`,
	Run:     loginFunc,
}

var loginIntCmd = &cobra.Command{
	Use:     "i",
	Aliases: []string{"interactive", "in"},
	Short:   "Opens login interactive session.",
	PreRun:  interactiveLoginInputFunc,
	Run:     loginFunc,
}

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Aliases: []string{"disconnect"},
	Short:   "Logout current user.",
	Long:    "Logout current user from Redmine. It deletes user credentials.",
	Run:     logoutFunc,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	loginCmd.Flags().StringVarP(&hostname, "server", "s", "", "Hostname of Redmine server (e.g. host.redmine.org)")
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password")

	_ = loginCmd.MarkFlagRequired("server")
	_ = loginCmd.MarkFlagRequired("username")
	_ = loginCmd.MarkFlagRequired("password")

	loginCmd.AddCommand(loginIntCmd)

}

func loginFunc(_ *cobra.Command, _ []string) {
	authCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	viper.Set(config.Hostname, hostname)
	req, ReqErr := RClient.NewAuthRequest(authCtx, username, password)
	if ReqErr != nil {
		log.Fatal("user request ", ReqErr)
	}

	var userApiResponse *client.UserApiResponse
	res, ResErr := RClient.Do(req, &userApiResponse)
	var urlError *url.Error
	if errors.As(ResErr, &urlError) {
		fmt.Println("You entered wrong hostname!")
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		fmt.Println("Wrong login credentials!")
		return
	default:
		fmt.Println(fmt.Sprintf("Cannot login user (%v)", res.StatusCode))
		return
	}

	user := userApiResponse.User
	viper.Set(config.ApiKey, user.ApiKey)
	writeConfigErr := viper.WriteConfig()

	if writeConfigErr != nil {
		log.Fatal("Unable to save config! :S", writeConfigErr)
	}

	fmt.Println("You have successfully logged in!")
}

func interactiveLoginInputFunc(_ *cobra.Command, _ []string) {
	var err error
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		fmt.Printf("stdin/stdout should be terminal")
		return
	}

	oldState, err := terminal.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := terminal.NewTerminal(screen, "")

	var hostOk, userOk, passOk bool
	hostname, hostOk = AskForHostname(t)
	username, userOk = AskForText(t, "Username: ", false)
	password, passOk = AskForText(t, "Password: ", true)

	_ = terminal.Restore(int(os.Stdout.Fd()), oldState)
	if !hostOk || !userOk || !passOk {
		defer os.Exit(1)
	}
}

func AskForHostname(t *terminal.Terminal) (string, bool) {
	previousHost := viper.GetString(config.Hostname)
	var prefix string
	if previousHost != "" {
		prefix = fmt.Sprintf("Hostname (%s): ", previousHost)
	} else {
		prefix = fmt.Sprintf("Hostname: ")
	}

	userInput, ok := AskForText(t, prefix, false)

	if userInput == "" && previousHost != "" {
		return previousHost, ok
	}

	return userInput, ok
}

func AskForText(t *terminal.Terminal, prefix string, hidden bool) (string, bool) {
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
		return line, false
	}

	return line, true
}

func logoutFunc(_ *cobra.Command, _ []string) {
	viper.Set(config.ApiKey, "")
	err := viper.WriteConfig()

	if err != nil {
		log.Fatal("Unable to save configuration", err)
	}

	fmt.Println("You have successfully logged out!")
}
