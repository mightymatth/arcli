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

var (
	limit    int
	spentOn  string
	hours    float32
	activity string
	comments string
)

var timeNow = now.EndOfDay()

func newTimeEntriesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "log",
		Aliases: []string{"l", "entries"},
		Short:   "Time entries on projects and issues",
	}

	c.AddCommand(newTimeEntriesListCmd())
	c.AddCommand(newTimeEntriesIssueCmd())
	c.AddCommand(newTimeEntriesProjectCmd())
	c.AddCommand(newTimeEntriesUpdateCmd())
	c.AddCommand(newTimeEntriesDeleteCmd())

	return c
}

func newTimeEntriesListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "all"},
		Short:   "List user time entries",
		Run:     timeEntriesListFunc,
	}

	c.Flags().IntVarP(&limit, "limit", "l", 10,
		"Limit number of results")

	return c
}

func timeEntriesListFunc(cmd *cobra.Command, _ []string) {
	limit := cmd.Flags().Lookup("limit").Value.String()
	queryParams := fmt.Sprintf("limit=%s&user_id=me", limit)
	logs, err := RClient.GetTimeEntries(queryParams)
	if err != nil {
		fmt.Println("Cannot get time entries:", err)
		return
	}

	t := utils.NewTable()
	t.AppendHeader(table.Row{"ID", "Project", "Issue ID",
		"Activity", "Hours", "Spent on", "Comment"})
	for _, log := range logs {
		t.AppendRow(table.Row{log.ID, log.Project.Name, log.Issue.String(),
			log.Activity.Name, log.Hours, relativeDateString(log.SpentOn),
			log.Comments})
	}

	t.Render()
}

func newTimeEntriesIssueCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "issue [id]",
		Args:    validIssueArgs(),
		Aliases: []string{"i"},
		Short:   "Add time entry to issue.",
		Run:     timeEntriesAddFunc(false),
	}

	c.Flags().StringVarP(&spentOn, "date", "d", "today",
		"The date the time was spent ('today', 'yesterday', '2020-01-15')")
	c.Flags().Float32VarP(&hours, "hours", "t", 0,
		"The number of spent hours")
	c.Flags().StringVarP(&activity, "activity", "a", "",
		"The name of activity for spent time (this overrides default config value)")
	c.Flags().StringVarP(&comments, "message", "m", "",
		"Short comment")
	_ = c.MarkFlagRequired("hours")

	return c
}

func newTimeEntriesProjectCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "project [id]",
		Args:    validProjectArgs(),
		Aliases: []string{"p"},
		Short:   "Add time entry to project",
		Run:     timeEntriesAddFunc(true),
	}

	c.Flags().StringVarP(&spentOn, "date", "d", "today",
		"The date the time was spent ('today', 'yesterday', '2020-01-15')")
	c.Flags().Float32VarP(&hours, "hours", "t", 0,
		"The number of spent hours")
	c.Flags().StringVarP(&activity, "activity", "a", "",
		"The name of activity for spent time (this overrides default config value)")
	c.Flags().StringVarP(&comments, "message", "m", "",
		"Short comment")
	_ = c.MarkFlagRequired("hours")

	return c
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

		activityID, exists := activities.Valid(activity)
		if !exists {
			fmt.Printf("Invalid activity (allowed ones: [%v])",
				utils.PrintWithDelimiter(activities.Names()))
			return
		}

		spentOnTime, err := spentOnParse(spentOn)
		if err != nil {
			fmt.Printf("Cannot parse date value: %v", err)
			return
		}

		var entryPost *client.TimeEntryPost
		if isProject {
			entryPost = &client.TimeEntryPost{
				ProjectID:  int(id),
				SpentOn:    *client.NewDateTime(*spentOnTime),
				Hours:      hours,
				ActivityID: int(activityID),
				Comments:   comments,
			}
		} else {
			entryPost = &client.TimeEntryPost{
				IssueID:    int(id),
				SpentOn:    *client.NewDateTime(*spentOnTime),
				Hours:      hours,
				ActivityID: int(activityID),
				Comments:   comments,
			}
		}

		entry, err := RClient.AddTimeEntry(*entryPost)
		if err != nil {
			fmt.Printf("Cannot create time entry: %v\n", err)
			return
		}

		fmt.Println("Time entry created!")
		entry.PrintTable()
	}
}

func newTimeEntriesUpdateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "update [id]",
		Args:    validTimeEntryArgs(),
		Aliases: []string{"u", "edit", "modify"},
		Short:   "Update time entry",
		Run:     timeEntriesUpdateFunc(),
	}

	c.Flags().StringVarP(&spentOn, "date", "d", "today",
		"The date the time was spent ('today', 'yesterday', '2020-01-15')")
	c.Flags().Float32VarP(&hours, "hours", "t", 0,
		"The number of spent hours")
	c.Flags().StringVarP(&activity, "activity", "a", "",
		"The name of activity for spent time (this overrides default config value)")
	c.Flags().StringVarP(&comments, "message", "m", "",
		"Short comment")

	return c
}

func timeEntriesUpdateFunc() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		entryID, _ := strconv.ParseInt(args[0], 10, 64)
		var entryUpdate client.TimeEntryPost

		if activity != "" {
			activities, err := RClient.GetActivities()
			if err != nil {
				fmt.Println("Cannot get time entry activities")
				return
			}

			activityID, exists := activities.Valid(activity)
			if !exists {
				fmt.Printf("Invalid activity (allowed ones: [%v])",
					utils.PrintWithDelimiter(activities.Names()))
				return
			}

			entryUpdate.ActivityID = int(activityID)
		}

		if spentOn != "" {
			spentOnTime, err := spentOnParse(spentOn)
			if err != nil {
				fmt.Printf("Cannot parse date value: %v", err)
				return
			}

			entryUpdate.SpentOn = *client.NewDateTime(*spentOnTime)
		}

		if cmd.Flags().Changed("message") {
			if comments == "" {
				comments = " "
			}
			entryUpdate.Comments = comments
		}

		if hours != 0 {
			entryUpdate.Hours = hours
		}

		err := RClient.UpdateTimeEntry(int(entryID), entryUpdate)
		if err != nil {
			fmt.Printf("Cannot update time entry: %v\n", err)
			return
		}
		fmt.Println("Time entry updated!")

		updatedEntry, err := RClient.GetTimeEntry(int(entryID))
		if err != nil {
			fmt.Printf("Time entry with ID %d cannot be fetched: %v", entryID, err)
		}
		updatedEntry.PrintTable()
	}
}

func newTimeEntriesDeleteCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "delete [id...]",
		Args:    validTimeEntryArgs(),
		Aliases: []string{"remove", "rm", "del"},
		Short:   "Delete time entry",
		Run:     timeEntriesDeleteFunc,
	}

	return c
}

func relativeDateString(dateTime client.DateTime) string {
	durationDays := int(timeNow.Sub(dateTime.Time).Hours() / 24)
	date := dateTime.Time.Format(client.DayDateFormat)

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

func validTimeEntryArgs() cobra.PositionalArgs {
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
		entryID, _ := strconv.ParseInt(arg, 10, 64)

		err := RClient.DeleteTimeEntry(int(entryID))
		if err != nil {
			fmt.Println("Cannot delete time entry:", err)
			return
		}

		fmt.Printf("Time entry with id %v successfully deleted.\n", arg)
	}
}

func spentOnModify(spentOn string) (string, error) {
	var modified string
	switch spentOn {
	case "today":
		modified = timeNow.Format(client.DateTimeFormat)
	case "yesterday":
		modified = timeNow.AddDate(0, 0, -1).Format(client.DateTimeFormat)
	default:
		_, err := time.Parse(client.DateTimeFormat, spentOn)
		if err != nil {
			return "", fmt.Errorf("invalid date format (use '%v' instead)",
				client.DateTimeFormat)
		}
		modified = spentOn
	}

	return modified, nil
}

func spentOnParse(spentOn string) (*time.Time, error) {
	modified, err := spentOnModify(spentOn)
	if err != nil {
		return nil, err
	}
	spentOnTime, _ := time.Parse(client.DateTimeFormat, modified)

	return &spentOnTime, nil
}
