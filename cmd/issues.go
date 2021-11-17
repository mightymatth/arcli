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

func newIssuesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "issues [id]",
		Args:    validIssueArgs(),
		Aliases: []string{"i", "tasks", "show"},
		Short:   "Shows issue details",
		Run:     issueFunc,
	}

	c.AddCommand(newMyIssuesCmd())
	c.AddCommand(newMyRelatedIssuesCmd())
	c.AddCommand(newMyWatchedIssuesCmd())

	return c
}

func validIssueArgs() cobra.PositionalArgs {
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

func issueFunc(_ *cobra.Command, args []string) {
	issueID, _ := strconv.ParseInt(args[0], 10, 64)
	issue, err := RClient.GetIssue(issueID)
	if err != nil {
		fmt.Printf("Cannot fetch issue with id %v\n", issueID)
		return
	}
	project := client.Project{ID: issue.Project.ID, Name: issue.Project.Name}
	fmt.Printf("[%v] %v (%v)\n", text.FgYellow.Sprint(project.ID), text.FgYellow.Sprint(project.Name), project.URL())
	fmt.Printf("  [%v] %v (%v)\n", text.FgGreen.Sprint(issue.ID), text.FgGreen.Sprint(issue.Subject), issue.URL())
	fmt.Printf("%v\n", issue.Description)
}

func newMyIssuesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "my",
		Aliases: []string{"assigned", "list", "ls"},
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

	return c
}

func newMyRelatedIssuesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "related",
		Aliases: []string{"rel"},
		Short:   "List all related issues assigned to the user (including groups)",
		Run: func(cmd *cobra.Command, args []string) {
			issues, err := RClient.GetMyRelatedIssues()
			if err != nil {
				fmt.Println("Cannot fetch my related issues:", err)
				return
			}

			drawIssues(issues)
		},
	}

	return c
}

func newMyWatchedIssuesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "watched",
		Aliases: []string{"w"},
		Short:   "List all issues watched by the user",
		Run: func(cmd *cobra.Command, args []string) {
			issues, err := RClient.GetMyWatchedIssues()
			if err != nil {
				log.Println("Cannot fetch watched issues:", err)
				return
			}

			drawIssues(issues)
		},
	}

	return c
}

func drawIssues(issues []client.Issue) {
	t := utils.NewTable()
	t.AppendHeader(table.Row{"ID", "Project", "Subject", "URL"})
	for _, issue := range issues {
		t.AppendRow(table.Row{issue.ID, issue.Project.Name, issue.Subject, issue.URL()})
	}

	t.Render()
}
