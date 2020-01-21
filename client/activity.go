package client

// Activity represents Redmine activity for time that's being tracked.
type Activity struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Activities represents a list of activities.
type Activities []Activity

type activitiesResponse struct {
	Activities Activities `json:"time_entry_activities"`
}

// GetActivities fetches all Activities that can be entered in time entry record. Project specific
// activities cannot be fetched with this method.
func (c *Client) GetActivities() (Activities, error) {
	req, err := c.getRequest("/enumerations/time_entry_activities.json", "")
	if err != nil {
		return nil, err
	}

	var response activitiesResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Activities, nil
}

// Valid checks whether activity with provided name parameter exists. If yes,
// returns its ID and true as second parameter; if not, return false as second parameter.
func (acts Activities) Valid(name string) (int64, bool) {
	theMap := make(map[string]int64)
	for _, activity := range acts {
		theMap[activity.Name] = activity.ID
	}

	activityID, exists := theMap[name]
	if !exists {
		return activityID, false
	}

	return activityID, true
}

// Names returns activity names.
func (acts Activities) Names() []string {
	theActs := make([]string, 0, len(acts))
	for _, activity := range acts {
		theActs = append(theActs, activity.Name)
	}

	return theActs
}
