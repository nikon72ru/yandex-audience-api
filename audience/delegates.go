package audience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//Delegate - a user who has been granted access to the current user account.
type Delegate struct {
	UserLogin string    `json:"user_login"`
	Perm      string    `json:"perm"`
	CreatedAt time.Time `json:"created_at"`
	Comment   string    `json:"comment"`
}

//DelegatesList - returns a list of representatives who have been granted access to the current user account.
func (c *Client) DelegatesList() ([]*Delegate, error) {
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
	}, "delegates")
	if err != nil {
		return nil, err
	}
	defer closer(resp.Body)
	var response struct {
		Delegates []*Delegate `json:"delegates"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) != 0 {
		return nil, response.Error()
	}
	return response.Delegates, nil
}

//CreateDelegate - adds the user login to the list of representatives for the current account.
func (c *Client) CreateDelegate(delegate *Delegate) error {
	requestStruct := struct {
		Delegate *Delegate `json:"delegate"`
		APIError
	}{Delegate: delegate}
	jsonBody, _ := json.Marshal(requestStruct)
	resp, err := c.Do(&http.Request{
		Method: http.MethodPut,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, "delegate")
	if err != nil {
		return err
	}
	defer closer(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&requestStruct); err != nil {
		return err
	}
	if len(requestStruct.Errors) != 0 {
		return requestStruct.Error()
	}
	return nil
}

//RemoveDelegate - deletes the user login from the list of representatives for the current account.
func (c *Client) RemoveDelegate(userLogin string) error {
	resp, err := c.Do(&http.Request{
		Method: http.MethodDelete,
	}, fmt.Sprintf("delegate?user_login=%s", userLogin))
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
