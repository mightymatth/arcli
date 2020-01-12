package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type TimeEntry struct {
	Id        int64     `json:"id"`
	Project   Entity    `json:"project"`
	Issue     EntityId  `json:"issue"`
	User      Entity    `json:"user"`
	Activity  Entity    `json:"activity"`
	Hours     float64   `json:"hours"`
	Comments  string    `json:"comments"`
	SpentOn   DateTime  `json:"spent_on"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

type TimeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

type TimeEntryResponse struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

func (c *Client) GetTimeEntries(queryParams string) ([]TimeEntry, error) {
	req, err := c.getRequest("/time_entries.json", queryParams)
	if err != nil {
		return nil, err
	}

	var response TimeEntriesResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.TimeEntries, nil
}

type TimeEntryBody struct {
	TimeEntry TimeEntryPost `json:"time_entry"`
}

type TimeEntryPost struct {
	IssueId    int      `json:"issue_id,omitempty"`
	ProjectId  int      `json:"project_id,omitempty"`
	SpentOn    DateTime `json:"spent_on"`
	Hours      float32  `json:"hours"`
	ActivityId int      `json:"activity_id"`
	Comments   string   `json:"comments"`
}

func (c *Client) AddTimeEntry(entry TimeEntryPost) (*TimeEntry, error) {
	req, err := c.postRequest("/time_entries.json", TimeEntryBody{TimeEntry: entry})
	if err != nil {
		return nil, err
	}

	var response TimeEntryResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.TimeEntry, nil
}

func (c *Client) DeleteTimeEntry(id int) error {
	req, err := c.deleteRequest(fmt.Sprintf("/time_entries/%v.json", id))
	if err != nil {
		return err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("there is no time entry with id %v", id)
	default:
		return fmt.Errorf("status code %v", resp.StatusCode)
	}
}

type DateTime struct {
	time.Time
}

func NewDateTime(time time.Time) *DateTime {
	return &DateTime{Time: time}
}

const DateTimeFormat = "2006-01-02"

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
