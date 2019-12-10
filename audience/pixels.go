package audience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//Pixel - type of segment, created by yandex pixel
type Pixel struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	UserQuantity7  int64     `json:"user_quantity_7"`
	UserQuantity30 int64     `json:"user_quantity_30"`
	UserQuantity90 int64     `json:"user_quantity_90"`
	CreateTime     time.Time `json:"create_time"`
}

//PixelsList - returns a list of existing user pixels.
func (c *Client) PixelsList() (*[]Pixel, error) {
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
	}, "pixels")
	if err != nil {
		return nil, err
	}
	defer closer(resp.Body)
	var response struct {
		Pixels []Pixel `json:"pixels"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) != 0 {
		return nil, response.Error()
	}
	return &response.Pixels, nil
}

//CreatePixel - creates a pixel with the specified parameters.
func (c *Client) CreatePixel(pixel *Pixel) error {
	jsonBody, _ := json.Marshal(pixel)
	resp, err := c.Do(&http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, "pixels")
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

//RemovePixel - deletes the specified pixel.
func (c *Client) RemovePixel(pixelID int64) error {
	resp, err := c.Do(&http.Request{
		Method: http.MethodDelete,
	}, fmt.Sprintf("pixel/%d", pixelID))
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

//UpdatePixel - changes the specified pixel.
func (c *Client) UpdatePixel(pixel *Pixel) error {
	jsonBody, _ := json.Marshal(pixel)
	resp, err := c.Do(&http.Request{
		Method: http.MethodPut,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, fmt.Sprintf("pixel/%d", pixel.ID))
	if err != nil {
		return err
	}
	defer closer(resp.Body)
	var response struct {
		Pixel Pixel `json:"pixel"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if len(response.Errors) != 0 {
		return response.Error()
	}
	*pixel = response.Pixel
	return nil
}

//UndeletePixel - recovers the deleted pixel.
func (c *Client) UndeletePixel(pixelID int64) error {
	resp, err := c.Do(&http.Request{
		Method: http.MethodPost,
	}, fmt.Sprintf("pixel/%d/undelete", pixelID))
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
