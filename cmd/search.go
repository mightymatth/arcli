package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/mightymatth/arcli/utils"
)

var searchCmd = &cobra.Command{
	Use:     "search [query]",
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(1),
	Short:   "Search Redmine",
	Run:     searchFunc,
}

var (
	searchOffset, searchLimit int
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IntVarP(&searchOffset, "offset", "o", 0, "Offset from first result")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 25, "Limit of given search results")
}

func searchFunc(_ *cobra.Command, args []string) {
	results, totalCount, err := RClient.GetSearchResults(args[0], searchOffset, searchLimit)
	if err != nil {
		fmt.Println("Search failed:", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Printf("Found %d results. Showing results from %d. to %d.\n",
		totalCount, searchOffset+1, searchOffset+len(results))

	t := utils.NewTable()
	t.AppendHeader(table.Row{"Resource ID", "Title", "URL"})
	for _, result := range results {
		t.AppendRow(table.Row{result.Id, result.Title, result.Url})
	}
	t.Render()
}
