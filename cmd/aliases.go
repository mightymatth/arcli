package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mightymatth/arcli/config"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
)

var aliasesCmd = &cobra.Command{
	Use:     "aliases",
	Aliases: []string{"a", "alias"},
	Short:   "Words that can be used instead of issue or project ids",
}

var aliasesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "all"},
	Short:   "List of all user aliases",
	Run: func(cmd *cobra.Command, args []string) {
		drawAliases()
	},
}

var aliasesAddCmd = &cobra.Command{
	Use:     "add [aliasName] [id]",
	Aliases: []string{"set", "new"},
	Args:    ValidAliasesAddArgs(),
	Short:   "Add alias entry",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.SetAlias(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("'%v: %v' has been successfully added to aliases.\n", args[0], args[1])
	},
}

var aliasesDeleteCmd = &cobra.Command{
	Use:     "delete [aliasName]",
	Aliases: []string{"remove", "rm", "del"},
	Args:    ValidAliasesDeleteArgs(),
	Short:   "Remove alias entry",
	Run: func(cmd *cobra.Command, args []string) {
		_, found := config.GetAlias(args[0])
		if !found {
			fmt.Printf("Alias with key '%v' does not exist, so can't be deleted.\n", args[0])
			return
		}

		err := config.SetAlias(args[0], "")
		if err != nil {
			fmt.Println("Cannot delete alias:", err)
			return
		}

		fmt.Printf("Alias with key '%v' has been deleted.\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	aliasesCmd.AddCommand(aliasesListCmd)
	aliasesCmd.AddCommand(aliasesAddCmd)
	aliasesCmd.AddCommand(aliasesDeleteCmd)
}

func drawAliases() {
	aliases := config.GetAliases()
	if len(aliases) == 0 {
		fmt.Println("You have no previously aliases set.")
		fmt.Printf("These can be set with: '%v'\n", aliasesAddCmd.UseLine())
		return
	}

	t := utils.NewTable()
	t.AppendHeader(table.Row{"Alias", "ID"})
	for key, val := range aliases {
		t.AppendRow(table.Row{key, val})
	}

	t.Render()
}

func ValidAliasesAddArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(2)(cmd, args)
		if err != nil {
			return err
		}

		keyPattern := "^[[:alnum:]-_]{1,30}$"
		if !regexp.MustCompile(keyPattern).MatchString(args[0]) {
			return fmt.Errorf("alias key must have pattern '%v'", keyPattern)
		}

		_, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("alias value must be integer")
		}

		return nil
	}
}

func ValidAliasesDeleteArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}

		return nil
	}
}
