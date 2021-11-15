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
	host, username, password string
)

func newLoginCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "login",
		Aliases: []string{"li"},
		Args:    cobra.ExactArgs(0),
		Short:   "Opens login interactive login session",
		PreRun:  interactiveLoginInputFunc,
		Run:     loginFunc,
	}

	c.AddCommand(newLoginInlineCmd())

	return c
}

func newLoginInlineCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "inline",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"i"},
		Short:   "Authenticate to Redmine server",
		Run:     loginFunc,
	}

	c.Flags().StringVarP(&host, "server", "s", "", "Host of Redmine server (e.g. https://host.redmine.org)")
	c.Flags().StringVarP(&username, "username", "u", "", "Username")
	c.Flags().StringVarP(&password, "password", "p", "", "Password")

	_ = c.MarkFlagRequired("server")
	_ = c.MarkFlagRequired("username")
	_ = c.MarkFlagRequired("password")

	return c
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Aliases: []string{"lo", "disconnect"},
		Short:   "Logout current user",
		Long:    "Logout current user from Redmine. It deletes user credentials.",
		Run:     logoutFunc,
	}
}

func loginFunc(_ *cobra.Command, _ []string) {
	authCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	viper.Set(config.Host, host)
	req, ReqErr := RClient.NewAuthRequest(authCtx, username, password)
	if ReqErr != nil {
		log.Println("User request:", ReqErr)
	}

	var userAPIResponse *client.UserAPIResponse
	res, ResErr := RClient.Do(req, &userAPIResponse)
	var urlError *url.Error
	if errors.As(ResErr, &urlError) {
		fmt.Println("User response:", ResErr)
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

	user := userAPIResponse.User
	viper.Set(config.APIKey, user.APIKey)
	writeConfigErr := viper.WriteConfig()

	if writeConfigErr != nil {
		log.Fatal("Unable to save config! :S", writeConfigErr)
	}

	fmt.Println("You have successfully logged in!")
}

func interactiveLoginInputFunc(_ *cobra.Command, _ []string) {
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
	host, hostOk = askForHost(t)
	username, userOk = askForText(t, "Username: ", false)
	password, passOk = askForText(t, "Password: ", true)

	_ = terminal.Restore(int(os.Stdout.Fd()), oldState)
	if !hostOk || !userOk || !passOk {
		defer os.Exit(1)
	}
}

func askForHost(t *terminal.Terminal) (string, bool) {
	previousHost := viper.GetString(config.Host)
	var prefix string
	if previousHost != "" {
		prefix = fmt.Sprintf("Host (%s): ", previousHost)
	} else {
		prefix = fmt.Sprintf("Host: ")
	}

	userInput, ok := askForText(t, prefix, false)

	if userInput == "" && previousHost != "" {
		return previousHost, ok
	}

	return userInput, ok
}

func askForText(t *terminal.Terminal, prefix string, hidden bool) (string, bool) {
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
	viper.Set(config.APIKey, "")
	err := viper.WriteConfig()

	if err != nil {
		log.Fatal("Unable to save configuration", err)
	}

	fmt.Println("You have successfully logged out!")
}
