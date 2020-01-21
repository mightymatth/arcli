package cmd

import (
	"fmt"
	"strconv"

	"github.com/mightymatth/arcli/config"

	"github.com/jedib0t/go-pretty/text"

	"github.com/mightymatth/arcli/client"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:     "projects [id]",
	Args:    validProjectArgs(),
	Aliases: []string{"p", "tasks"},
	Short:   "Shows project details",
	Run:     projectFunc,
}

var myProjectsCmd = &cobra.Command{
	Use:     "my",
	Aliases: []string{"all", "show", "ls", "list"},
	Short:   "List all projects visible to user",
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := RClient.GetProjects()
		if err != nil {
			fmt.Println("Cannot fetch projects:", err)
			return
		}

		drawProjects(projects)
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(myProjectsCmd)
}

func drawProjects(projects []client.Project) {
	for _, project := range projects {
		if project.Parent == nil {
			fmt.Printf("[%v] %v\n", text.FgYellow.Sprint(project.ID),
				text.FgYellow.Sprint(project.Name))
		} else {
			fmt.Printf(" ‣ [%v] %v\n", text.FgCyan.Sprint(project.ID),
				text.FgCyan.Sprint(project.Name))
		}

	}
}

func validProjectArgs() cobra.PositionalArgs {
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
			return fmt.Errorf("project id must be integer")
		}
		return nil
	}
}

func projectFunc(_ *cobra.Command, args []string) {
	projectID, _ := strconv.ParseInt(args[0], 10, 64)
	project, err := RClient.GetProject(projectID)
	if err != nil {
		fmt.Printf("Cannot fetch project with id %v\n", projectID)
		return
	}

	fmt.Printf("[%v] %v\n", text.FgYellow.Sprint(project.ID), text.FgYellow.Sprint(project.Identifier))
	fmt.Printf("%v (%v)\n", text.FgGreen.Sprint(project.Name), project.URL())
	fmt.Printf("%v\n", project.Description)
}
