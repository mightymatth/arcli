package cmd

import (
	"fmt"
	"log"

	"github.com/mightymatth/arcli/client"
	"github.com/spf13/cobra"
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
	user, err := client.RClient.GetUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-15s%-15v\n", "User ID", user.Id)
	fmt.Printf("%-15s%-15v\n", "Email", user.Email)
	fmt.Printf("%-15s%-15v\n", "First name", user.FirstName)
	fmt.Printf("%-15s%-15v\n", "Last name", user.LastName)
}
