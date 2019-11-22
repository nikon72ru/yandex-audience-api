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
	apiUrl        = "https://api-audience.yandex.ru/"
)

//content_types constants
const (
	idfaGain = "idfa_gaid"
	clientId = "client_id"
	mac      = "mac"
	crm      = "crm"
)

type Client struct {
	apiVersion string
	token      string
	hc         *http.Client
}

type YandexAudience struct {
	token string
}

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

func (c *Client) Close() error {
	c.hc = nil
	return nil
}

func (c *Client) SegmentsList(pixel ...int) (*[]map[string]interface{}, error) {
	requestPath := fmt.Sprintf("%s%s/management/segments", apiUrl, c.apiVersion)
	if len(pixel) > 0 {
		requestPath += fmt.Sprintf("?pixel=%d", pixel[0])
	}
	requestUrl, err := url.Parse(requestPath)
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodGet,
		URL:    requestUrl,
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

func (c *Client) CreateSegmentFromFile(name, filename, contentType string) (*UploadingSegment, error) {
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
	requestUrl, err := url.Parse(fmt.Sprintf("%s%s/management/segments/upload_file", apiUrl, c.apiVersion))
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodPost,
		URL:    requestUrl,
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
		ApiError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return nil, err
	}
	if respStruct.Segment.Id != 0 {
		respStruct.Segment.Name = name
		respStruct.Segment.ContentType = contentType
		return &respStruct.Segment, nil
	} else if len(respStruct.Errors) != 0 {
		return nil, respStruct.Error()
	} else {
		return nil, errors.New("unexpected answer format")
	}
}

func (c *Client) SaveUploadedSegment(segment *UploadingSegment) error {
	requestUrl, err := url.Parse(fmt.Sprintf("%s%s/management/segment/%d/confirm?", apiUrl, c.apiVersion, segment.Id))
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
		URL:    requestUrl,
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
		ApiError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return err
	}
	if respStruct.Segment.Id != 0 {
		segment = &respStruct.Segment
		return nil
	} else if len(respStruct.Errors) != 0 {
		return respStruct.Error()
	} else {
		return errors.New("unexpected answer format")
	}
}

func (c *Client) RemoveSegment(id int64) (bool, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s%s/management/segment/%d", apiUrl, c.apiVersion, id))
	if err != nil {
		return false, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodDelete,
		URL:    requestUrl,
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
		ApiError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return false, err
	}
	if len(respStruct.Errors) != 0 {
		return false, respStruct.Error()
	} else {
		return true, nil
	}
}
