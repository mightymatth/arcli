package cmd

import (
	"log"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/mightymatth/arcli/utils"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"me"},
	Short:   "Overall account info",
	Run:     statusFunc,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func statusFunc(_ *cobra.Command, _ []string) {
	user, err := RClient.GetUser()
	if err != nil {
		log.Fatal(err)
	}

	t := utils.NewTable()
	t.AppendRow(table.Row{"User ID", user.Id})
	t.AppendRow(table.Row{"Email", user.Email})
	t.AppendRow(table.Row{"First name", user.FirstName})
	t.AppendRow(table.Row{"Last name", user.LastName})
	t.Render()
}
