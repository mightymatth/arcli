package cmd

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"

	tm "github.com/buger/goterm"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"me"},
	Short:   "Overall account info",
	Long: `Shows user info and statistics of several periods showing: sum of tracked time hours,
average hours per tracked time, number of issues and number of projects.`,
	Run: statusFunc,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func statusFunc(_ *cobra.Command, _ []string) {
	user := "Loading user..."
	var today, yesterday, thisWeek, lastWeek, thisMonth, lastMonth string

	refresh := make(chan struct{}, 7)

	var g errgroup.Group
	g.Go(asyncUserResult(&user, refresh))
	g.Go(asyncPeriodResult(SpentOnToday, &today, refresh))
	g.Go(asyncPeriodResult(SpentOnYesterday, &yesterday, refresh))
	g.Go(asyncPeriodResult(SpentOnThisWeek, &thisWeek, refresh))
	g.Go(asyncPeriodResult(SpentOnLastWeek, &lastWeek, refresh))
	g.Go(asyncPeriodResult(SpentOnThisMonth, &thisMonth, refresh))
	g.Go(asyncPeriodResult(SpentOnLastMonth, &lastMonth, refresh))

	drawScreen := func() {
		_, _ = tm.Println(user)
		_, _ = tm.Println("PERIOD      ", fmt.Sprintf("%-7s %-7s %-8s %-8s",
			"HOURS", "H/LOG", "# of I", "# of P"))
		_, _ = tm.Println("Today       ", today)
		_, _ = tm.Println("Yesterday   ", yesterday)
		_, _ = tm.Println("This Week   ", thisWeek)
		_, _ = tm.Println("Last Week   ", lastWeek)
		_, _ = tm.Println("This Month  ", thisMonth)
		_, _ = tm.Println("Last Month  ", lastMonth)
		tm.Flush()
		tm.MoveCursorUp(8)
	}

	var writing sync.WaitGroup
	writing.Add(1)
	go func() {
		drawScreen()
		for range refresh {
			drawScreen()
		}
		writing.Done()
	}()

	err := g.Wait()
	close(refresh)
	writing.Wait()

	if err != nil {
		fmt.Println("Failed to get status:", err)
		return
	}
}

func asyncUserResult(dest *string, refresh chan<- struct{}) func() error {
	return func() error {
		defer func() { refresh <- struct{}{} }()

		u, err := RClient.GetUser()
		if err == nil {
			*dest = fmt.Sprintf("[%d] %s %s (%s)", u.Id, u.FirstName, u.LastName, u.Email)
		} else {
			*dest = "Cannot fetch user."
		}

		return err
	}
}

func asyncPeriodResult(t TimeSpentOn, dest *string, refresh chan<- struct{}) func() error {
	return func() error {
		defer func() { refresh <- struct{}{} }()

		result, err := getDataForPeriod(t)
		if err == nil {
			*dest = result
		} else {
			*dest = "ERR"
		}

		return err
	}
}

func getDataForPeriod(spentOn TimeSpentOn) (string, error) {
	entries, err := RClient.GetTimeEntries(fmt.Sprintf("spent_on=%s&user_id=me&limit=200", spentOn))
	if err != nil {
		return "", fmt.Errorf("cannot get period data (%v): %v", spentOn, err)
	}

	var hoursSum float64
	issues := make(map[int64]struct{})
	projects := make(map[int64]struct{})

	for _, entry := range entries {
		hoursSum += entry.Hours
		issues[entry.Issue.Id] = struct{}{}
		projects[entry.Project.Id] = struct{}{}
	}
	delete(issues, 0) // time tracked on projects

	issueCount := len(issues)
	projectCount := len(projects)
	var hoursAvg float64
	if len(entries) != 0 {
		hoursAvg = hoursSum / float64(len(entries))
	}

	return fmt.Sprintf("%-7s %-7s %-8d %-8d",
		formatFloat(hoursSum), formatFloat(hoursAvg),
		issueCount, projectCount), nil
}

func formatFloat(num float64) string {
	s := fmt.Sprintf("%.1f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

type TimeSpentOn string

const (
	SpentOnToday     TimeSpentOn = "t"
	SpentOnYesterday TimeSpentOn = "ld"
	SpentOnThisWeek  TimeSpentOn = "w"
	SpentOnLastWeek  TimeSpentOn = "lw"
	SpentOnThisMonth TimeSpentOn = "m"
	SpentOnLastMonth TimeSpentOn = "lm"
)
