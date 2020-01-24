package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"
)

// TimeEntry represents Redmine time entry model.
type TimeEntry struct {
	ID        int64     `json:"id"`
	Project   entity    `json:"project"`
	Issue     entityID  `json:"issue"`
	User      entity    `json:"user"`
	Activity  entity    `json:"activity"`
	Hours     float64   `json:"hours"`
	Comments  string    `json:"comments"`
	SpentOn   DateTime  `json:"spent_on"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

// PrintTable prints table in suitable format.
func (te TimeEntry) PrintTable() {
	t := utils.NewTable()
	t.AppendHeader(table.Row{"Entry ID", "Project Name", "Issue ID", "Hours", "Activity", "Comment", "Spent On"})
	t.AppendRow(table.Row{fmt.Sprint(te.ID), te.Project.Name, te.Issue.String(),
		fmt.Sprint(te.Hours), te.Activity.Name, te.Comments, te.SpentOn.Format(DayDateFormat)})
	t.Render()
}

type timeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

type timeEntryResponse struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

// GetTimeEntries fetches time entries for requested queryParams
func (c *Client) GetTimeEntries(queryParams string) ([]TimeEntry, error) {
	req, err := c.getRequest("/time_entries.json", queryParams)
	if err != nil {
		return nil, err
	}

	var response timeEntriesResponse
	res, err := c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusOK:
		return response.TimeEntries, nil
	default:
		return nil, fmt.Errorf("cannot get time entries (status %v)", res.StatusCode)
	}
}

// GetTimeEntry fetches time entry for given ID.
func (c *Client) GetTimeEntry(id int) (*TimeEntry, error) {
	req, err := c.getRequest(fmt.Sprintf("/time_entries/%d.json", id), "")
	if err != nil {
		return nil, err
	}

	var response timeEntryResponse
	res, err := c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusOK:
		return &response.TimeEntry, nil
	default:
		return nil, fmt.Errorf("cannot get time entry (status %v)", res.StatusCode)
	}
}

type timeEntryBody struct {
	TimeEntry TimeEntryPost `json:"time_entry"`
}

// TimeEntryPost represents data which should be placed to request body
// while creating a new time entry.
type TimeEntryPost struct {
	IssueID    int      `json:"issue_id,omitempty"`
	ProjectID  int      `json:"project_id,omitempty"`
	SpentOn    DateTime `json:"spent_on,omitempty"`
	Hours      float32  `json:"hours,omitempty"`
	ActivityID int      `json:"activity_id,omitempty"`
	Comments   string   `json:"comments,omitempty"`
}

// AddTimeEntry adds new time entry.
func (c *Client) AddTimeEntry(entry TimeEntryPost) (*TimeEntry, error) {
	req, err := c.postRequest("/time_entries.json", timeEntryBody{TimeEntry: entry})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		var teRes timeEntryResponse
		err = json.NewDecoder(resp.Body).Decode(&teRes)
		if err != nil {
			return nil, err
		}
		return &teRes.TimeEntry, nil
	case http.StatusUnprocessableEntity:
		var errRes error422Response
		err = json.NewDecoder(resp.Body).Decode(&errRes)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(utils.PrintWithDelimiter(errRes.Errors))
	default:
		return nil, fmt.Errorf("status %v", resp.StatusCode)
	}
}

// UpdateTimeEntry adds new time entry.
func (c *Client) UpdateTimeEntry(id int, entry TimeEntryPost) error {
	req, err := c.putRequest(fmt.Sprintf("/time_entries/%d.json", id), timeEntryBody{TimeEntry: entry})
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnprocessableEntity:
		var errRes error422Response
		err = json.NewDecoder(resp.Body).Decode(&errRes)
		if err != nil {
			return err
		}
		return fmt.Errorf(utils.PrintWithDelimiter(errRes.Errors))
	default:
		return fmt.Errorf("status %v", resp.StatusCode)
	}
}

// DeleteTimeEntry deletes time entry with requested ID.
func (c *Client) DeleteTimeEntry(id int) error {
	req, err := c.deleteRequest(fmt.Sprintf("/time_entries/%v.json", id))
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("there is no time entry with id %v", id)
	default:
		return fmt.Errorf("status code %v", resp.StatusCode)
	}
}

// DateTime custom representation of date.
type DateTime struct {
	time.Time
}

// NewDateTime creates new DateTime for specific time.Time.
func NewDateTime(time time.Time) *DateTime {
	return &DateTime{Time: time}
}

// DateTimeFormat represents date format
const DateTimeFormat = "2006-01-02"

// DayDateFormat represents day and date format
const DayDateFormat = "Mon, 2006-01-02"

// UnmarshalJSON override
func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	parsedTime, err := time.Parse(`"`+DateTimeFormat+`"`, string(data))
	if err != nil {
		return err
	}

	*t = DateTime{parsedTime}

	return nil
}

// MarshalJSON override
func (t DateTime) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(DateTimeFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, DateTimeFormat)
	b = append(b, '"')
	return b, nil
}
