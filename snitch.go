package deadmanssnitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Snitch represents the details of a snitch
type Snitch struct {

	// The snitch's identifying token.
	Token string `json:"token,omitempty"`

	// API URL to retrieve data about this specific Snitch.
	Href string `json:"href,omitempty"`

	// The name of the snitch.
	Name string `json:"name,omitempty"`

	// The list of keyword tags for this snitch.
	Tags []string `json:"tags,omitempty"`

	// The status of the snitch. It could be:
	// "pending"	The snitch is new and your job has not yet checked in.
	// "healthy"	Your job has checked in since the beginning of the last period.
	// "failed"	Your job has not checked in since the beginning of the last period. (At least one alert has been sent.)
	// "errored"	Your job has reported that is has errored. (At least one alert has been sent.) Error Notices are only available on some plans.
	// "paused"	The snitch has been paused and will not worry about your failing job until your job checks-in again after you fix it.
	Status string `json:"status,omitempty"`

	// Any user-supplied notes about this snitch.
	Notes string `json:"notes,omitempty"`

	// The last time your job checked in healthy, as an ISO 8601 datetime with millisecond precision. The timezone is always UTC. If your job has not checked in healthy yet, this will be null.
	CheckedInAt string `json:"checked_in_at,omitempty"`

	// The url your job should hit to check-in.
	CheckInURL string `json:"check_in_url,omitempty"`

	// The size of the period window. If your job does not check-in during an entire period, you will be notified and the snitch status will show up as "failed". The interval can be "15_minute", "30_minute", "hourly", "daily", "weekly", or "monthly".
	Interval string `json:"interval,omitempty"`

	// The type of alerts the snitch will use. basic will have a static deadline that it will expect to hear from it by, while smart will learn when your snitch checks in, moving the deadline closer so you can be alerted sooner.
	AlertType string `json:"alert_type,omitempty"`

	// When the snitch was created, as an ISO 8601 datetime with millisecond precision. The timezone is always UTC.
	CreatedAt string `json:"created_at,omitempty"`
}

// ListSnitches returns a list of snitches with the provided `filters`.  An empty filter will result in all snitches.
func (c *Client) ListSnitches(filters []string) (*[]Snitch, error) {
	snitchList := []Snitch{}

	tagList := strings.Join(filters, ",")

	body, err := c.do("GET", fmt.Sprintf("snitches?tags=%s", tagList), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &snitchList)
	if err != nil {
		return nil, err
	}

	return &snitchList, nil
}

// CheckIn calls the check-in url for the snitch
func (c *Client) CheckIn(token string) error {
	CheckInURL := fmt.Sprintf("https://nosnch.in/%s", token)

	resp, err := http.Get(CheckInURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading the API check-in response: %v", err)
	}

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	return fmt.Errorf("Error checking in GET %s, HTTP %d. %s", CheckInURL, resp.StatusCode, body)
}

// GetSnitch returns a single snitch
func (c *Client) GetSnitch(token string) (*Snitch, error) {

	snitch := Snitch{}

	body, err := c.do("GET", fmt.Sprintf("snitches/%s", token), nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &snitch)
	if err != nil {
		return nil, err
	}

	return &snitch, nil
}

// CreateSnitch creates a new snitch
func (c *Client) CreateSnitch(snitch *Snitch) (*Snitch, error) {
	snitchData, err := json.Marshal(snitch)
	if err != nil {
		return nil, err
	}

	body, err := c.do("POST", "snitches", snitchData)
	if err != nil {
		return nil, err
	}

	newSnitch := Snitch{}
	err = json.Unmarshal(body, &newSnitch)
	if err != nil {
		return nil, err
	}

	return &newSnitch, nil
}

// UpdateSnitch updates the snitch identified by `token`
// The `updatedSnitch` parameter accepts a Snitch object in which you may
// provide only the attributes you wish to change. Empty fields
// in the object will not be touched.
func (c *Client) UpdateSnitch(token string, updatedSnitch *Snitch) (*Snitch, error) {
	snitchData, err := json.Marshal(updatedSnitch)
	if err != nil {
		return nil, err
	}

	body, err := c.do("PATCH", fmt.Sprintf("snitches/%s", token), snitchData)
	if err != nil {
		return nil, err
	}

	newSnitch := Snitch{}
	err = json.Unmarshal(body, &newSnitch)
	if err != nil {
		return nil, err
	}

	return &newSnitch, nil
}

// AddTags adds the given tags to the snitch, leaving existing tags unchanged
func (c *Client) AddTags(token string, newTags []string) error {
	newTagData, err := json.Marshal(newTags)
	if err != nil {
		return err
	}

	_, err = c.do("POST", fmt.Sprintf("snitches/%s/tags", token), newTagData)
	if err != nil {
		return err
	}

	return nil
}

// RemoveTags removes the given tags from the snitch
func (c *Client) RemoveTags(token string, rmTags []string) error {
	for _, tag := range rmTags {
		_, err := c.do("DELETE", fmt.Sprintf("snitches/%s/tags/%s", token, tag), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// PauseSnitch pauses a snitch
func (c *Client) PauseSnitch(token string) error {

	_, err := c.do("POST", fmt.Sprintf("snitches/%s/pause", token), nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSnitch deletes a snitch
func (c *Client) DeleteSnitch(token string) error {
	_, err := c.do("DELETE", fmt.Sprintf("snitches/%s", token), nil)
	if err != nil {
		return err
	}

	return nil
}
