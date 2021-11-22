package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/client"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"strings"
)

func newStatusCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "status",
		Aliases: []string{"me"},
		Short:   "Overall account info",
		Long: `Shows user info and statistics of several periods showing: sum of tracked time hours,
average hours per tracked time, number of issues and number of projects.`,
		Run: statusFunc,
	}

	return c
}

func statusFunc(_ *cobra.Command, _ []string) {
	var user client.User
	var today, yesterday, thisWeek, lastWeek, thisMonth, lastMonth periodData

	var g errgroup.Group
	g.Go(asyncUserResult(&user))
	g.Go(asyncPeriodResult(spentOnToday, &today))
	g.Go(asyncPeriodResult(spentOnYesterday, &yesterday))
	g.Go(asyncPeriodResult(spentOnThisWeek, &thisWeek))
	g.Go(asyncPeriodResult(spentOnLastWeek, &lastWeek))
	g.Go(asyncPeriodResult(spentOnThisMonth, &thisMonth))
	g.Go(asyncPeriodResult(spentOnLastMonth, &lastMonth))

	err := g.Wait()

	if err != nil {
		fmt.Println("Failed to get status:", err)
		return
	}

	fmt.Printf("[%d] %s %s (%s)\n", user.ID, user.FirstName, user.LastName, user.Email)

	t := utils.NewTable()
	t.AppendHeader(table.Row{"PERIOD", "HOURS", "H/LOG", "# of I", "# of P"})
	appendRow(t, "Today", today)
	appendRow(t, "Yesterday", yesterday)
	appendRow(t, "This Week", thisWeek)
	appendRow(t, "Last Week", lastWeek)
	appendRow(t, "This Month", thisMonth)
	appendRow(t, "Last Month", lastMonth)

	t.Render()
}

func asyncUserResult(dest *client.User) func() error {
	return func() error {
		user, err := RClient.GetUser()
		if err != nil {
			return err
		}

		*dest = *user

		return nil
	}
}

func asyncPeriodResult(t timeSpentOn, dest *periodData) func() error {
	return func() error {
		data, err := getDataForPeriod(t)
		if err != nil {
			return err
		}

		*dest = data

		return nil
	}
}

func getDataForPeriod(spentOn timeSpentOn) (periodData, error) {
	entries, err := RClient.GetTimeEntries(fmt.Sprintf("spent_on=%s&user_id=me&limit=200", spentOn))
	if err != nil {
		return periodData{}, fmt.Errorf("cannot get period data (%v): %v", spentOn, err)
	}

	var hoursSum float64
	issues := make(map[int64]struct{})
	projects := make(map[int64]struct{})

	for _, entry := range entries {
		hoursSum += entry.Hours
		issues[entry.Issue.ID] = struct{}{}
		projects[entry.Project.ID] = struct{}{}
	}
	delete(issues, 0) // time tracked on projects

	issueCount := len(issues)
	projectCount := len(projects)
	var hoursAvg float64
	if len(entries) != 0 {
		hoursAvg = hoursSum / float64(len(entries))
	}

	return periodData{
		hoursSum:     hoursSum,
		hoursAvg:     hoursAvg,
		issueCount:   issueCount,
		projectCount: projectCount,
	}, nil
}

func appendRow(t table.Writer, period string, data periodData) {
	t.AppendRow(table.Row{
		period, formatFloat(data.hoursSum), formatFloat(data.hoursAvg),
		data.issueCount, data.projectCount,
	})
}

func formatFloat(num float64) string {
	s := fmt.Sprintf("%.1f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

type timeSpentOn string

const (
	spentOnToday     timeSpentOn = "t"
	spentOnYesterday timeSpentOn = "ld"
	spentOnThisWeek  timeSpentOn = "w"
	spentOnLastWeek  timeSpentOn = "lw"
	spentOnThisMonth timeSpentOn = "m"
	spentOnLastMonth timeSpentOn = "lm"
)

type periodData struct {
	hoursSum     float64
	hoursAvg     float64
	issueCount   int
	projectCount int
}
