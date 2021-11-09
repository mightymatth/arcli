package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jinzhu/now"
	"github.com/mightymatth/arcli/client"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
)

var spentOnMonth string

const cellHorizontalSpaces = 3

// CalendarCell represents the calendar's cells.
type CalendarCell struct {
	Day   int
	Hours float64
}

func newViewCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "view",
		Aliases: []string{"v"},
		Short:   "Shows data in different views",
	}

	c.AddCommand(newTimeEntriesCalendarCmd())

	return c
}

// newTimeEntriesCalendarCmd add the calendar command to the list of available commands.
func newTimeEntriesCalendarCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "calendar",
		Aliases: []string{"c", "cal"},
		Short:   "List user's time entries in a calendar view",
		RunE:    timeEntriesCalendarFunc,
	}

	c.Flags().StringVarP(&spentOnMonth, "month", "m", "current",
		"The month the times was spent ('current', '2020-01')")

	return c
}

// timeEntriesCalendarFunc is the function that is called when the command calendar is ran.
func timeEntriesCalendarFunc(cmd *cobra.Command, _ []string) error {
	// Check if the month format is correct.
	var re = regexp.MustCompile(`current|[\d]{4}-[\d]{2}`)

	if !re.MatchString(spentOnMonth) {
		return errors.New("the format is not correct. The following formats are supported: \"current\", \"2021-11\"")
	}

	fmt.Print("Loading data...\r")

	// Define the days of th week.
	var daysOfWeek = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	// Get the date depending on the user's choice.
	var date *now.Now

	if spentOnMonth == "current" {
		date = now.With(time.Now())
	} else {
		var dateParsed, _ = time.Parse("2006-01-02", fmt.Sprintf("%s-01", spentOnMonth))
		date = now.With(dateParsed)
	}

	// Get the wanted date.
	var formattedDate = date.Format("January 2006")

	// Get the totals days in the month.
	var _, _, totalDaysInMonth = date.EndOfMonth().Date()

	// Get the weekday of the first day of the month.
	var weekday = int(date.BeginningOfMonth().Weekday())

	// Calculate the total weeks in the month.
	var totalWeeksInMonth int = int(math.Ceil(float64(totalDaysInMonth+weekday) / 7))

	// Get the time entries for the month.
	var spentOnFrom = date.BeginningOfMonth().Format("2006-01-02")
	var spentOnTo = date.EndOfMonth().Format("2006-01-02")

	var timeEntries, err = RClient.GetTimeEntries(fmt.Sprintf("from=%s&to=%s&user_id=me&limit=200", spentOnFrom, spentOnTo))

	if err != nil {
		fmt.Println("Cannot get time entries:", err)
		os.Exit(1)
	}

	// Calculate the total of hours logged per day and store it in an slice.
	var weeks = make([][]CalendarCell, totalWeeksInMonth)
	var dayIndex = 1
	var biggestLoggedHours = 0

	for i := 0; i < totalWeeksInMonth; i++ {
		weeks[i] = make([]CalendarCell, 7)

		for u := 0; u < 7; u++ {
			if (i == 0 && u < weekday) || dayIndex > totalDaysInMonth {
				weeks[i][u] = CalendarCell{
					Day: -1,
				}
			} else {
				var hours float64 = 0.0

				// Loop over all time entries and count the number of hours
				// for the current day.
				for _, timeEntry := range timeEntries {
					if timeEntry.SpentOn.Day() == dayIndex {
						hours += timeEntry.Hours
					}
				}

				weeks[i][u] = CalendarCell{
					Day:   dayIndex,
					Hours: hours,
				}

				// Store the biggest logged hours to be able
				// to display the calendar's cells width.
				biggestLoggedHours = int(math.Max(float64(len(formatFloatTwoFloating(hours))), float64(biggestLoggedHours)))

				dayIndex += 1
			}
		}
	}

	var cellWidth = (cellHorizontalSpaces * 2) + biggestLoggedHours

	// Calculate the spaces needed to filled the empty space.
	var dateHeaderSpacesNeeded = ((cellWidth * 7) + 7) - len(formattedDate) - 2

	// Write the date on top of the calendar.
	timeEntriesCalendarPrintSeparator("+", "-", cellWidth)
	fmt.Printf("| %s%s|\n", tm.Color(formattedDate, tm.CYAN), strings.Repeat(" ", dateHeaderSpacesNeeded))

	// Show the days of the week.
	timeEntriesCalendarPrintSeparator("+", "-", cellWidth)

	for _, day := range daysOfWeek {
		fmt.Printf("| %s%s", tm.Color(day, tm.CYAN), strings.Repeat(" ", cellWidth-4))
	}

	fmt.Print("|\n")
	timeEntriesCalendarPrintSeparator("+", "-", cellWidth)

	// Show the calendar days cells.
	for _, daysCells := range weeks {
		for _, dayCell := range daysCells {
			if dayCell.Day == -1 {
				fmt.Printf("|%s", strings.Repeat(" ", cellWidth))
			} else {
				var spaces = 3

				if dayCell.Day < 10 {
					spaces = 2
				}

				fmt.Printf("| %s%s", tm.Color(strconv.Itoa(dayCell.Day), tm.CYAN), strings.Repeat(" ", cellWidth-spaces))
			}
		}

		fmt.Println("|")
		timeEntriesCalendarPrintSeparator("|", " ", cellWidth)

		// Display hours.
		for _, dayCell := range daysCells {
			var timeLogged = formatFloatTwoFloating(dayCell.Hours)

			if dayCell.Day == -1 || dayCell.Hours == 0 {
				timeLogged = strings.Repeat(" ", biggestLoggedHours)
			}

			var (
				spacesLeft  = cellHorizontalSpaces
				spacesRight = cellHorizontalSpaces + (biggestLoggedHours - len(timeLogged))
			)

			if len(timeLogged) != biggestLoggedHours && (spacesLeft+spacesRight)%2 == 0 {
				var spacesDifference = (spacesRight - spacesLeft) / 2
				spacesLeft += spacesDifference
				spacesRight -= spacesDifference
			}

			fmt.Printf("|%s%s%s", strings.Repeat(" ", spacesLeft), tm.Color(tm.Bold(timeLogged), tm.GREEN), strings.Repeat(" ", spacesRight))
		}

		fmt.Println("|")

		timeEntriesCalendarPrintSeparator("+", "-", cellWidth)
	}

	// Print a summary table of the time entries grouped by projects.
	if len(timeEntries) > 0 {
		timeEntriesPrintSummary(timeEntries)
	}

	return nil
}

// timeEntriesPrintSummary print a summary table of the time entries by projects.
func timeEntriesPrintSummary(timeEntries []client.TimeEntry) {
	fmt.Println()

	// Reverse the time entries slice.
	reversedtimeEntries := []client.TimeEntry{}

	for i := range timeEntries {
		timeEntry := timeEntries[len(timeEntries)-1-i]
		reversedtimeEntries = append(reversedtimeEntries, timeEntry)
	}

	// Store information by projects.
	var projectsInfo = make(map[string]float64)

	for _, timeEntry := range reversedtimeEntries {
		projectsInfo[timeEntry.Project.Name] += timeEntry.Hours
	}

	// Sort the keys to be able to print the projects alphabetically.
	projectsInfoKeys := make([]string, 0, len(projectsInfo))

	for key := range projectsInfo {
		projectsInfoKeys = append(projectsInfoKeys, key)
	}

	sort.Strings(projectsInfoKeys)

	// Print the summary table.
	t := utils.NewTable()
	t.SetTitle("Summary")
	t.SetStyle(table.StyleDefault)
	t.AppendHeader(table.Row{"Project Name", "Total Hours"})

	for _, projectName := range projectsInfoKeys {
		t.AppendRow(table.Row{projectName, projectsInfo[projectName]})
	}

	t.Render()
}

func formatFloatTwoFloating(num float64) string {
	s := fmt.Sprintf("%.2f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

// timeEntriesCalendarPrintSeparator print a row separator in the calendar.
func timeEntriesCalendarPrintSeparator(border, separator string, cellWidth int) {
	var tableWidth = (cellWidth * 7) + 8

	for i := 0; i < tableWidth; i++ {
		if i == 0 || i%(cellWidth+1) == 0 {
			fmt.Print(border)
		} else {
			fmt.Print(separator)
		}
	}

	fmt.Println()
}
