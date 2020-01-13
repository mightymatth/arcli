package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/now"

	"github.com/mightymatth/arcli/config"

	"github.com/mightymatth/arcli/client"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
)

var timeEntriesCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"l", "entries"},
	Short:   "Time entries on projects and issues.",
}

var timeEntriesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "all"},
	Short:   "List user time entries.",
	Run:     timeEntriesListFunc,
}

var timeEntriesIssueCmd = &cobra.Command{
	Use:     "issue [id]",
	Args:    ValidIssueArgs(),
	Aliases: []string{"i"},
	Short:   "Add time entry to issue.",
	Run:     timeEntriesAddFunc(false),
}

var timeEntriesProjectCmd = &cobra.Command{
	Use:     "project [id]",
	Args:    ValidProjectArgs(),
	Aliases: []string{"p"},
	Short:   "Add time entry to project.",
	Run:     timeEntriesAddFunc(true),
}

var timeEntriesDeleteCmd = &cobra.Command{
	Use:     "delete [id...]",
	Args:    ValidDeleteTimeEntryArgs(),
	Aliases: []string{"remove", "rm", "del"},
	Short:   "Delete time entry.",
	Run:     timeEntriesDeleteFunc,
}

var (
	limit    int
	spentOn  string
	hours    float32
	activity string
	comments string
)

func init() {
	rootCmd.AddCommand(timeEntriesCmd)

	timeEntriesCmd.AddCommand(timeEntriesListCmd)
	timeEntriesCmd.AddCommand(timeEntriesIssueCmd)
	timeEntriesCmd.AddCommand(timeEntriesProjectCmd)
	timeEntriesCmd.AddCommand(timeEntriesDeleteCmd)

	timeEntriesListCmd.Flags().IntVarP(&limit, "limit", "l", 10,
		"Limit number of results")

	for _, cmd := range []*cobra.Command{timeEntriesIssueCmd, timeEntriesProjectCmd} {
		cmd.Flags().StringVarP(&spentOn, "date", "d", "today",
			"The date the time was spent ('today', 'yesterday', '2020-01-15')")
		cmd.Flags().Float32VarP(&hours, "hours", "t", 0,
			"The number of spent hours")
		cmd.Flags().StringVarP(&activity, "activity", "a", "",
			"The name of activity for spent time (this overrides default config value)")
		cmd.Flags().StringVarP(&comments, "message", "m", "",
			"Short comment")
		_ = cmd.MarkFlagRequired("hours")
	}
}

var timeNow = now.EndOfDay()

func timeEntriesListFunc(cmd *cobra.Command, _ []string) {
	limit := cmd.Flags().Lookup("limit").Value.String()
	queryParams := fmt.Sprintf("limit=%s&user_id=me", limit)
	logs, err := RClient.GetTimeEntries(queryParams)
	if err != nil {
		fmt.Println("Cannot get time entries")
		return
	}

	t := utils.NewTable()
	t.AppendHeader(table.Row{"ID", "Project", "Issue ID",
		"Activity", "Hours", "Spent on", "Comment"})
	for _, log := range logs {
		t.AppendRow(table.Row{log.Id, log.Project.Name, log.Issue.Id,
			log.Activity.Name, log.Hours, RelativeDateString(log.SpentOn),
			log.Comments})
	}

	t.Render()
}

func RelativeDateString(dateTime client.DateTime) string {
	durationDays := int(timeNow.Sub(dateTime.Time).Hours() / 24)
	date := dateTime.Time.Format(client.DateTimeFormat)

	switch {
	case durationDays < 0:
		return date
	case durationDays == 0:
		return fmt.Sprintf("today (%v)", date)
	case durationDays == 1:
		return fmt.Sprintf("yesterday (%v)", date)
	default:
		return fmt.Sprintf("%v days ago (%v)", durationDays, date)
	}
}

func timeEntriesAddFunc(isProject bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		id, _ := strconv.ParseInt(args[0], 10, 64)

		activities, err := RClient.GetActivities()
		if err != nil {
			fmt.Println("Cannot get time entry activities")
			return
		}

		if activity == "" {
			activity = config.Defaults()[string(config.Activity)]
			if activity == "" {
				fmt.Println("Provide activity either by flag or setting default.")
				return
			}
		}

		activityId, exists := activities.Valid(activity)
		if !exists {
			fmt.Printf("Invalid activity (allowed ones: [%v])",
				printWithDelimiter(activities.Names()))
			return
		}

		switch spentOn {
		case "today":
			spentOn = timeNow.Format(client.DateTimeFormat)
		case "yesterday":
			spentOn = timeNow.AddDate(0, 0, -1).Format(client.DateTimeFormat)
		default:
			_, err = time.Parse(client.DateTimeFormat, spentOn)
			if err != nil {
				fmt.Printf("Invalid date format (use '%v' instead)\n",
					client.DateTimeFormat)
				return
			}
		}
		spentOnTime, _ := time.Parse(client.DateTimeFormat, spentOn)

		var entryPost *client.TimeEntryPost
		if isProject {
			entryPost = &client.TimeEntryPost{
				ProjectId:  int(id),
				SpentOn:    *client.NewDateTime(spentOnTime),
				Hours:      hours,
				ActivityId: int(activityId),
				Comments:   comments,
			}
		} else {
			entryPost = &client.TimeEntryPost{
				IssueId:    int(id),
				SpentOn:    *client.NewDateTime(spentOnTime),
				Hours:      hours,
				ActivityId: int(activityId),
				Comments:   comments,
			}
		}

		_, err = RClient.AddTimeEntry(*entryPost)
		if err != nil {
			fmt.Printf("Cannot create time entry: %v\n", err)
			return
		}

		fmt.Println("Time entry created!")
	}
}

func ValidDeleteTimeEntryArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := cobra.MinimumNArgs(1)(cmd, args)
		if err != nil {
			return err
		}

		for _, arg := range args {
			_, err = strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return fmt.Errorf("time entry id must be integer, but given %v", arg)
			}
		}

		return nil
	}
}

func timeEntriesDeleteFunc(_ *cobra.Command, args []string) {
	for _, arg := range args {
		entryId, _ := strconv.ParseInt(arg, 10, 64)

		err := RClient.DeleteTimeEntry(int(entryId))
		if err != nil {
			fmt.Println("Cannot delete time entry:", err)
			return
		}

		fmt.Printf("Time entry with id %v successfully deleted.\n", arg)
	}
}
