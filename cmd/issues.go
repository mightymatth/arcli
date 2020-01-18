package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mightymatth/arcli/config"

	"github.com/jedib0t/go-pretty/text"

	"github.com/mightymatth/arcli/client"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:     "issues [id]",
	Args:    ValidIssueArgs(),
	Aliases: []string{"i", "tasks", "show"},
	Short:   "Shows issue details",
	Run:     IssueFunc,
}

var myIssuesCmd = &cobra.Command{
	Use:     "my",
	Aliases: []string{"assigned", "all", "list", "ls"},
	Short:   "List all issues assigned to the user",
	Run: func(cmd *cobra.Command, args []string) {
		issues, err := RClient.GetMyIssues()
		if err != nil {
			fmt.Println("Cannot fetch my issues:", err)
			return
		}

		drawIssues(issues)
	},
}

var myWatchedIssuesCmd = &cobra.Command{
	Use:   "watched",
	Short: "List all issues watched by user",
	Run: func(cmd *cobra.Command, args []string) {
		issues, err := RClient.GetMyWatchedIssues()
		if err != nil {
			log.Println("Cannot fetch watched issues:", err)
			return
		}

		drawIssues(issues)
	},
}

func init() {
	rootCmd.AddCommand(issuesCmd)
	issuesCmd.AddCommand(myIssuesCmd)
	issuesCmd.AddCommand(myWatchedIssuesCmd)
}

func drawIssues(issues []client.Issue) {
	t := utils.NewTable()
	t.AppendHeader(table.Row{"ID", "Project", "Subject", "URL"})
	for _, issue := range issues {
		t.AppendRow(table.Row{issue.Id, issue.Project.Name, issue.Subject, issue.URL()})
	}

	t.Render()
}

func ValidIssueArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}

		val, found := config.GetAlias(args[0])
		if found {
			args[0] = val
			return nil
		}

		_, err = strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("issue id must be integer")
		}
		return nil
	}
}

func IssueFunc(_ *cobra.Command, args []string) {
	issueId, _ := strconv.ParseInt(args[0], 10, 64)
	issue, err := RClient.GetIssue(issueId)
	if err != nil {
		fmt.Printf("Cannot fetch issue with id %v\n", issueId)
		return
	}
	project := client.Project{Id: issue.Project.Id, Name: issue.Project.Name}
	fmt.Printf("[%v] %v (%v)\n", text.FgYellow.Sprint(project.Id), text.FgYellow.Sprint(project.Name), project.URL())
	fmt.Printf("  [%v] %v (%v)\n", text.FgGreen.Sprint(issue.Id), text.FgGreen.Sprint(issue.Subject), issue.URL())
	fmt.Printf("%v\n", issue.Description)
}
