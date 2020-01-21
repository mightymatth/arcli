package client

type Activity struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Activities []Activity

type ActivitiesResponse struct {
	Activities Activities `json:"time_entry_activities"`
}

func (c *Client) GetActivities() (Activities, error) {
	req, err := c.getRequest("/enumerations/time_entry_activities.json", "")
	if err != nil {
		return nil, err
	}

	var response ActivitiesResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Activities, nil
}

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

func (acts Activities) Names() []string {
	theActs := make([]string, 0, len(acts))
	for _, activity := range acts {
		theActs = append(theActs, activity.Name)
	}

	return theActs
}
