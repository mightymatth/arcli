package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"

	"golang.org/x/sync/errgroup"

	"github.com/mightymatth/arcli/client"

	"github.com/spf13/cobra"
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

func statusFunc(_ *cobra.Command, _ []string) {
	var g errgroup.Group
	user := make(chan *client.User, 1)
	today := make(chan float64, 1)
	yesterday := make(chan float64, 1)
	thisWeek := make(chan float64, 1)
	lastWeek := make(chan float64, 1)
	thisMonth := make(chan float64, 1)
	lastMonth := make(chan float64, 1)

	g.Go(func() error {
		u, err := RClient.GetUser()
		if err == nil {
			user <- u
		}
		return err
	})

	g.Go(func() error {
		hours, err := getHoursForPeriod(SpentOnToday)
		if err == nil {
			today <- hours
		}
		return err
	})

	g.Go(func() error {
		hours, err := getHoursForPeriod(SpentOnYesterday)
		if err == nil {
			yesterday <- hours
		}
		return err
	})

	g.Go(func() error {
		hours, err := getHoursForPeriod(SpentOnThisWeek)
		if err == nil {
			thisWeek <- hours
		}
		return err
	})

	g.Go(func() error {
		hours, err := getHoursForPeriod(SpentOnLastWeek)
		if err == nil {
			lastWeek <- hours
		}
		return err
	})

	g.Go(func() error {
		hours, err := getHoursForPeriod(SpentOnThisMonth)
		if err == nil {
			thisMonth <- hours
		}
		return err
	})

	g.Go(func() error {
		hours, err := getHoursForPeriod(SpentOnLastMonth)
		if err == nil {
			lastMonth <- hours
		}
		return err
	})

	err := g.Wait()
	close(user)
	close(today)
	close(yesterday)
	close(thisWeek)
	close(lastWeek)
	close(thisMonth)
	close(lastMonth)

	if err != nil {
		fmt.Println("Failed to get status:", err)
		return
	}

	u := <-user
	fmt.Printf("[%d] %s %s (%s)\n", u.Id, u.FirstName, u.LastName, u.Email)

	t := utils.NewTable()
	t.AppendRow(table.Row{"Period", "Hours"})
	t.AppendRow(table.Row{"Today", <-today})
	t.AppendRow(table.Row{"Yesterday", <-yesterday})
	t.AppendRow(table.Row{"This week", <-thisWeek})
	t.AppendRow(table.Row{"This month", <-thisMonth})
	t.AppendRow(table.Row{"Last week", <-lastWeek})
	t.AppendRow(table.Row{"Last month", <-lastMonth})
	t.Render()
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
