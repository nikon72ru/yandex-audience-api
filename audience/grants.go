package audience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//Grant - segment management permissions information.
type Grant struct {
	UserLogin string    `json:"user_login"`
	CreatedAt time.Time `json:"created_at"`
	Comment   string    `json:"comment"`
}

//GrantsList - returns information about segment management permissions.
func (c *Client) GrantsList(segmentID int64) (*[]Grant, error) {
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
	}, fmt.Sprintf("segment/%d/grants", segmentID))
	if err != nil {
		return nil, err
	}
	defer closer(resp.Body)
	var response struct {
		Grants []Grant `json:"grants"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) != 0 {
		return nil, response.Error()
	}
	return &response.Grants, nil
}

//CreateGrant - creates permission to manage a segment.
func (c *Client) CreateGrant(segmentID int64, grant *Grant) error {
	reqURL, _ := url.Parse(fmt.Sprintf("%s%s/management/segment/%d/grants", apiURL, c.apiVersion, segmentID))
	jsonBody, _ := json.Marshal(grant)
	resp, err := c.Do(&http.Request{
		Method: http.MethodPut,
		URL:    reqURL,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, fmt.Sprintf("segment/%d/grant", segmentID))
	if err != nil {
		return err
	}
	defer closer(resp.Body)
	var response struct {
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if len(response.Errors) != 0 {
		return response.Error()
	}
	return nil
}

//RemoveGrant - removes permission to manage a segment.
func (c *Client) RemoveGrant(segmentID int64, userLogin string) error {
	resp, err := c.Do(&http.Request{
		Method: http.MethodDelete,
	}, fmt.Sprintf("segment/%d/grant/%s", segmentID, userLogin))
	if err != nil {
		return err
	}
	defer closer(resp.Body)
	var response struct {
		Success bool `json:"success"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if len(response.Errors) != 0 {
		return response.Error()
	}
	if !response.Success {
		return ErrNotDeleted
	}
	return nil
}
