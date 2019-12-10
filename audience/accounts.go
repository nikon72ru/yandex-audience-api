package audience

import (
	"encoding/json"
	"net/http"
	"time"
)

//perms
const (
	View = "view"
	Edit = "edit"
)

//Account - account represented by the current user.
type Account struct {
	UserLogin string    `json:"user_login"`
	Perm      string    `json:"perm"`
	CreatedAt time.Time `json:"created_at"`
}

//AccountsList - returns a list of accounts the current user is a representative of.
func (c *Client) AccountsList() (*[]Account, error) {
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
	}, "accounts")
	if err != nil {
		return nil, err
	}
	defer closer(resp.Body)
	var response struct {
		Accounts []Account `json:"accounts"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) != 0 {
		return nil, response.Error()
	}
	return &response.Accounts, nil
}
