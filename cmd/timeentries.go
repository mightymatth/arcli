package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mightymatth/arcli/config"

	"github.com/mightymatth/arcli/client"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
)

var timeEntriesCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"entries"},
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
	Run:     timeEntriesIssueFunc,
}

var timeEntriesDeleteCmd = &cobra.Command{
	Use:     "delete [id]",
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

	timeEntriesListCmd.Flags().IntVarP(&limit, "limit", "l", 10,
		"Limit number of results")

	timeEntriesIssueCmd.Flags().StringVarP(&spentOn, "date", "d", "today",
		"The date the time was spent ('today', 'yesterday', '2020-01-15')")
	timeEntriesIssueCmd.Flags().Float32VarP(&hours, "hours", "t", 0,
		"The number of spent hours")
	timeEntriesIssueCmd.Flags().StringVarP(&activity, "activity", "a", "",
		"The name of activity for spent time (this overrides default config value)")
	timeEntriesIssueCmd.Flags().StringVarP(&comments, "message", "m", "",
		"Short comment")
	_ = timeEntriesIssueCmd.MarkFlagRequired("hours")

	timeEntriesCmd.AddCommand(timeEntriesDeleteCmd)
}

var timeNow = time.Now()

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
		"Activity", "Hours", "Spent on"})
	for _, log := range logs {
		t.AppendRow(table.Row{log.Id, log.Project.Name, log.Issue.Id,
			log.Activity.Name, log.Hours, RelativeDateString(log.SpentOn)})
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

func timeEntriesIssueFunc(_ *cobra.Command, args []string) {
	issueId, _ := strconv.ParseInt(args[0], 10, 64)

	activities, err := RClient.GetActivities()
	if err != nil {
		fmt.Println("cannot get time entry activities")
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
		fmt.Printf("invalid activity (allowed ones: [%v])", printWithDelimiter(activities.Names()))
		return
	}

	switch spentOn {
	case "today":
		spentOn = timeNow.Format(client.DateTimeFormat)
	case "yesterday":
		spentOn = timeNow.AddDate(0, 0, -1).Format(client.DateTimeFormat)
	default:
		_, err = time.Parse(`"`+client.DateTimeFormat+`"`, spentOn)
		if err != nil {
			fmt.Printf("invalid date format (use '%v' instead)",
				client.DateTimeFormat)
		}
	}
	spentOnTime, _ := time.Parse(client.DateTimeFormat, spentOn)

	entryPost := client.TimeEntryPost{
		IssueId:    int(issueId),
		SpentOn:    *client.NewDateTime(spentOnTime),
		Hours:      hours,
		ActivityId: int(activityId),
		Comments:   comments,
	}

	_, err = RClient.AddTimeEntry(entryPost)
	if err != nil {
		fmt.Printf("Cannot create time entry: %v\n", err)
		return
	}

	fmt.Println("Time entry created!")
}

func timeEntriesDeleteFunc(_ *cobra.Command, args []string) {
	entryId, _ := strconv.ParseInt(args[0], 10, 64)

	err := RClient.DeleteTimeEntry(int(entryId))
	if err != nil {
		fmt.Println("Cannot delete time entry:", err)
		return
	}

	fmt.Println("Time entry successfully deleted.")
}
