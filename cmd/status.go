package cmd

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"

	tm "github.com/buger/goterm"
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
	user := "Loading user..."
	var today, yesterday, thisWeek, lastWeek, thisMonth, lastMonth string

	refresh := make(chan refreshData)

	var g errgroup.Group
	g.Go(asyncUserResult(&user, refresh))
	g.Go(asyncPeriodResult(spentOnToday, &today, refresh))
	g.Go(asyncPeriodResult(spentOnYesterday, &yesterday, refresh))
	g.Go(asyncPeriodResult(spentOnThisWeek, &thisWeek, refresh))
	g.Go(asyncPeriodResult(spentOnLastWeek, &lastWeek, refresh))
	g.Go(asyncPeriodResult(spentOnThisMonth, &thisMonth, refresh))
	g.Go(asyncPeriodResult(spentOnLastMonth, &lastMonth, refresh))

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
		for refreshData := range refresh {
			refreshData.update()
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

func asyncUserResult(dest *string, refresh chan<- refreshData) func() error {
	return func() error {
		u, err := RClient.GetUser()
		var result string
		if err == nil {
			result = fmt.Sprintf("[%d] %s %s (%s)", u.ID, u.FirstName, u.LastName, u.Email)
		} else {
			result = "Cannot fetch user."
		}
		refresh <- refreshData{Dest: dest, Value: result}

		return err
	}
}

func asyncPeriodResult(t timeSpentOn, dest *string, refresh chan<- refreshData) func() error {
	return func() error {
		data, err := getDataForPeriod(t)
		var result string
		if err == nil {
			result = data
		} else {
			result = "ERR"
		}
		refresh <- refreshData{Dest: dest, Value: result}

		return err
	}
}

func getDataForPeriod(spentOn timeSpentOn) (string, error) {
	entries, err := RClient.GetTimeEntries(fmt.Sprintf("spent_on=%s&user_id=me&limit=200", spentOn))
	if err != nil {
		return "", fmt.Errorf("cannot get period data (%v): %v", spentOn, err)
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

	return fmt.Sprintf("%-7s %-7s %-8d %-8d",
		formatFloat(hoursSum), formatFloat(hoursAvg),
		issueCount, projectCount), nil
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

type refreshData struct {
	Dest  *string
	Value string
}

func (rd *refreshData) update() {
	*rd.Dest = rd.Value
}
