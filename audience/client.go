package audience

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Errors section
var (
	ErrTokenIsNotSet = errors.New("yandex audience token isn't set")
	ErrNotDeleted    = errors.New("not deleted")
	ErrNotCreated    = errors.New("not created")
)

//constants
const (
	tokenVariable = "YANDEX_AUDIENCE_TOKEN"
	apiURL        = "https://api-audience.yandex.ru/"
)

//Client - a client of yandex audience API
type Client struct {
	token      string
	apiVersion string
	apiURL     string
	hc         *http.Client
}

//NewClient - create a new client to work with API
func NewClient(ctx context.Context) (*Client, error) {
	var client Client
	//Get token from context or env
	if os.Getenv(tokenVariable) == "" {
		if tok, ok := ctx.Value(tokenVariable).(string); ok && tok != "" {
			client.token = tok
		} else {
			return nil, ErrTokenIsNotSet
		}
	} else {
		client.token = os.Getenv(tokenVariable)
	}
	//Creating http client
	client.hc = &http.Client{}
	//Set api version: now available only v1
	client.apiVersion = "v1"
	//Set api URL
	client.apiURL = apiURL
	return &client, nil
}

//Do - append authorization header with token and call simple http.Client Do method
func (c *Client) Do(req *http.Request, path string) (*http.Response, error) {
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("OAuth %s", c.token)},
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s/management/%s", c.apiURL, c.apiVersion, path))
	if err != nil {
		return nil, err
	}
	req.URL = u
	return c.hc.Do(req)
}

//Close - close the client (all requests after will return errors)
func (c *Client) Close() error {
	c.hc = nil
	return nil
}

//Error - format describing return errors
type Error struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
	Location  string `json:"location"`
}

//APIError - API returned error
type APIError struct {
	Errors  []Error `json:"errors"`
	Code    int     `json:"code"`
	Message string  `json:"message"`
}

func (e *APIError) Error() error {
	err, _ := json.Marshal(e.Errors)
	return fmt.Errorf("%d: %s ([%s])", e.Code, e.Message, err)
}

func closer(p io.Closer) {
	if err := p.Close(); err != nil {
		log.Printf("can't close: %s", err.Error())
	}
}
