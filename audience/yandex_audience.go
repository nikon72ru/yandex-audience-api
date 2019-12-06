package audience

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// Errors section
var (
	ErrTokenIsNotSet = errors.New("yandex audience token isn't set")
)

//constants
const (
	tokenVariable = "YANDEX_AUDIENCE_TOKEN"
	apiURL        = "https://api-audience.yandex.ru/"
)

//content_types constants
const (
	IdfaGain = "idfa_gaid"
	ClientID = "client_id"
	Mac      = "mac"
	Crm      = "crm"
)

//Client - a client of yandex audience API
type Client struct {
	apiVersion string
	token      string
	hc         *http.Client
}

//YandexAudience - ...
type YandexAudience struct {
	token string
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
	return &client, nil
}

//Close - close the client (all requests after will return errors)
func (c *Client) Close() error {
	c.hc = nil
	return nil
}

//SegmentsList - returns a list of existing segments available to the user.
func (c *Client) SegmentsList(pixel ...int) (*[]map[string]interface{}, error) {
	requestPath := fmt.Sprintf("%s%s/management/segments", apiURL, c.apiVersion)
	if len(pixel) > 0 {
		requestPath += fmt.Sprintf("?pixel=%d", pixel[0])
	}
	requestURL, err := url.Parse(requestPath)
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodGet,
		URL:    requestURL,
		Header: http.Header{"Authorization": {fmt.Sprintf("OAuth %s", c.token)}},
	})
	if err != nil {
		return nil, err
	}
	var rawMap struct {
		Segments []map[string]interface{} `json:"segments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawMap); err != nil {
		return nil, err
	}
	//TODO: Apply a method to separate segments by type and break them into structures
	//Problem: Problem: Different formats returns as one array. Bad fields-description in API documentation
	return &rawMap.Segments, nil
}

//CreateFileSegment - creates a segment from a data file. The file must have at least 1000 entries.
func (c *Client) CreateFileSegment(segment *UploadingSegment, filename string) (*UploadingSegment, error) {
	var f *os.File
	var err error
	if f, err = os.Open(filename); err != nil {
		return nil, err
	}
	rp, wp := io.Pipe()
	mpw := multipart.NewWriter(wp)
	errorChan := make(chan error, 1)
	go func() {
		var part io.Writer
		var err error
		defer wp.Close()
		defer f.Close()
		if part, err = mpw.CreateFormField("key"); err != nil {
			return
		}
		if _, err = part.Write([]byte("KEY")); err != nil {
			return
		}

		if part, err = mpw.CreateFormFile("file", filename); err != nil {
			errorChan <- err
			return
		}
		if _, err = io.Copy(part, f); err != nil {
			errorChan <- err
		}
		if err = mpw.Close(); err != nil {
			errorChan <- err
		}
		errorChan <- nil
	}()
	requestURL, err := url.Parse(fmt.Sprintf("%s%s/management/segments/upload_file", apiURL, c.apiVersion))
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodPost,
		URL:    requestURL,
		Header: http.Header{
			"Authorization": {fmt.Sprintf("OAuth %s", c.token)},
			"Content-Type":  {mpw.FormDataContentType()},
		},
		Body: rp,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := <-errorChan; err != nil {
		return nil, err
	}
	var respStruct struct {
		Segment UploadingSegment `json:"segment"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return nil, err
	}
	if respStruct.Segment.ID != 0 {
		return &respStruct.Segment, nil
	} else if len(respStruct.Errors) != 0 {
		return nil, respStruct.Error()
	} else {
		return nil, errors.New("unexpected answer format")
	}
}

//SaveUploadedSegment - saves a segment created from a data file.
func (c *Client) SaveUploadedSegment(segment *UploadingSegment) error {
	requestURL, err := url.Parse(fmt.Sprintf("%s%s/management/segment/%d/confirm?", apiURL, c.apiVersion, segment.ID))
	if err != nil {
		return err
	}
	jsonBody, err := json.Marshal(struct {
		Segment *UploadingSegment `json:"segment"`
	}{segment})
	if err != nil {
		return err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodPost,
		URL:    requestURL,
		Header: http.Header{
			"Authorization": {fmt.Sprintf("OAuth %s", c.token)},
		},
		Body: ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var respStruct struct {
		Segment UploadingSegment `json:"segment"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return err
	}
	if respStruct.Segment.ID != 0 {
		segment = &respStruct.Segment
		return nil
	} else if len(respStruct.Errors) != 0 {
		return respStruct.Error()
	} else {
		return errors.New("unexpected answer format")
	}
}

//RemoveSegment - deletes the specified segment.
func (c *Client) RemoveSegment(id int64) (bool, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s%s/management/segment/%d", apiURL, c.apiVersion, id))
	if err != nil {
		return false, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodDelete,
		URL:    requestURL,
		Header: http.Header{
			"Authorization": {fmt.Sprintf("OAuth %s", c.token)},
		},
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var respStruct struct {
		Success bool `json:"success"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return false, err
	}
	if len(respStruct.Errors) != 0 {
		return false, respStruct.Error()
	}
	return true, nil
}

//CreatePixelSegment - creates a segment of type "pixel" with the specified parameters.
//If different conditions are used when creating a segment (for example, several labels are specified),
//then a user who satisfies all the specified conditions at the same time will get into the segment.
func (c *Client) CreatePixelSegment(segment *PixelSegment) error {
	requestURL, err := url.Parse(fmt.Sprintf("%s%s/management/segments/create_pixel?", apiURL, c.apiVersion))
	if err != nil {
		return nil
	}
	jsonBody, err := json.Marshal(struct {
		Segment *PixelSegment `json:"segment"`
	}{segment})
	if err != nil {
		return nil
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodPost,
		URL:    requestURL,
		Header: http.Header{
			"Authorization": {fmt.Sprintf("OAuth %s", c.token)},
		},
		Body: ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var respStruct struct {
		Segment PixelSegment `json:"segment"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return nil
	}
	if respStruct.Segment.ID != 0 {
		segment = &respStruct.Segment
		return nil
	} else if len(respStruct.Errors) != 0 {
		return respStruct.Error()
	} else {
		return errors.New("unexpected answer format")
	}
}

//CreateLookalikeSegment - creates a “lookalike” type segment with the specified parameters.
func (c *Client) CreateLookalikeSegment(segment *LookalikeSegment) error {
	requestURL, err := url.Parse(fmt.Sprintf("%s%s/management/segments/create_lookalike?", apiURL, c.apiVersion))
	if err != nil {
		return nil
	}
	jsonBody, err := json.Marshal(struct {
		Segment *LookalikeSegment `json:"segment"`
	}{segment})
	if err != nil {
		return nil
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodPost,
		URL:    requestURL,
		Header: http.Header{
			"Authorization": {fmt.Sprintf("OAuth %s", c.token)},
		},
		Body: ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var respStruct struct {
		Segment LookalikeSegment `json:"segment"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return nil
	}
	if respStruct.Segment.ID != 0 {
		segment = &respStruct.Segment
		return nil
	}
	if len(respStruct.Errors) != 0 {
		return respStruct.Error()
	}
	return errors.New("unexpected answer format")
}
