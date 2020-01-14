package cmd

import (
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/mightymatth/arcli/client"

	"github.com/spf13/cobra"

	tm "github.com/buger/goterm"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"me"},
	Short:   "Overall account info",
	Run:     statusFunc,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

type asyncResult struct {
	Result float64
	Err    error
}

func statusFunc(_ *cobra.Command, _ []string) {

	user := "Loading user..."
	var today, yesterday, thisWeek, lastWeek, thisMonth, lastMonth string

	userResult := make(chan *client.User)
	todayResult := make(chan asyncResult)
	yesterdayResult := make(chan asyncResult)
	thisWeekResult := make(chan asyncResult)
	lastWeekResult := make(chan asyncResult)
	thisMonthResult := make(chan asyncResult)
	lastMonthResult := make(chan asyncResult)
	asyncResults := []chan asyncResult{todayResult, yesterdayResult, thisWeekResult,
		lastWeekResult, thisMonthResult, lastMonthResult}
	refresh := make(chan struct{})

	var g errgroup.Group
	g.Go(func() error {
		u, err := RClient.GetUser()
		if err == nil {
			userResult <- u
		} else {
			userResult <- nil
		}
		return err
	})
	g.Go(asyncHoursResult(SpentOnToday, todayResult))
	g.Go(asyncHoursResult(SpentOnYesterday, yesterdayResult))
	g.Go(asyncHoursResult(SpentOnThisWeek, thisWeekResult))
	g.Go(asyncHoursResult(SpentOnLastWeek, lastWeekResult))
	g.Go(asyncHoursResult(SpentOnThisMonth, thisMonthResult))
	g.Go(asyncHoursResult(SpentOnLastMonth, lastMonthResult))

	go func() {
		saveResult := func(r asyncResult, dest *string) {
			if r.Err != nil {
				*dest = "ERR"
			} else {
				*dest = fmt.Sprintf("%v", r.Result)
			}
		}
		for i := 0; i < 7; i++ {
			select {
			case u := <-userResult:
				if u == nil {
					user = "Cannot fetch user."
				} else {
					user = fmt.Sprintf("[%d] %s %s (%s)", u.Id, u.FirstName, u.LastName, u.Email)
				}
			case r := <-todayResult:
				saveResult(r, &today)
			case r := <-yesterdayResult:
				saveResult(r, &yesterday)
			case r := <-thisWeekResult:
				saveResult(r, &thisWeek)
			case r := <-lastWeekResult:
				saveResult(r, &lastWeek)
			case r := <-thisMonthResult:
				saveResult(r, &thisMonth)
			case r := <-lastMonthResult:
				saveResult(r, &lastMonth)
			}
			refresh <- struct{}{}
		}
		close(refresh)
	}()

	drawScreen := func() {
		_, _ = tm.Println(user)
		_, _ = tm.Println("PERIOD      ", "HOURS")
		_, _ = tm.Println("Today       ", today)
		_, _ = tm.Println("Yesterday   ", yesterday)
		_, _ = tm.Println("This Week   ", thisWeek)
		_, _ = tm.Println("Last Week   ", lastWeek)
		_, _ = tm.Println("This Month  ", thisMonth)
		_, _ = tm.Println("This Month  ", lastMonth)
		tm.Flush()
		tm.MoveCursor(1, -8)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		//tm.Clear()
		drawScreen()
		for range refresh {
			drawScreen()
		}
		wg.Done()
	}()

	err := g.Wait()
	close(userResult)
	for _, res := range asyncResults {
		close(res)
	}

	if err != nil {
		fmt.Println("Failed to get status:", err)
		return
	}

	wg.Wait()
}

func asyncHoursResult(t TimeSpentOn, ch chan<- asyncResult) func() error {
	return func() error {
		hours, err := getHoursForPeriod(t)
		if err == nil {
			ch <- asyncResult{Result: hours, Err: nil}
		} else {
			ch <- asyncResult{Result: hours, Err: err}
		}
		return err
	}
}

func getHoursForPeriod(spentOn TimeSpentOn) (float64, error) {
	entries, err := RClient.GetTimeEntries(fmt.Sprintf("spent_on=%s&user_id=me&limit=200", spentOn))
	if err != nil {
		return 0, fmt.Errorf("cannot get last month hours: %v", err)
	}

	return calculateHours(entries), nil
}

func calculateHours(entries []client.TimeEntry) float64 {
	var sum float64
	for _, entry := range entries {
		sum += entry.Hours
	}

	return sum
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
