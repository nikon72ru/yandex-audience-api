package audience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

//Yandex audience const values
const (
	IdfaGain  = "idfa_gaid"
	ClientID  = "client_id"
	Mac       = "mac"
	Crm       = "crm"
	GoalID    = "goal_id"
	SegmentID = "segment_id"
	CounterID = "counter_id"
)

//BaseSegment defines the fields of the base segment (these fields exist in each type of segment)
type BaseSegment struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
	Owner      string    `json:"owner"`
}

//PixelSegment - a segment created by pixel.
type PixelSegment struct {
	BaseSegment
	PixelID                int    `json:"pixel_id"`
	PeriodLength           int    `json:"period_length"`
	TimesQuantity          int    `json:"times_quantity"`
	TimesQuantityOperation string `json:"times_quantity_operation"`
	UtmSource              string `json:"utm_source"`
	UtmContent             string `json:"utm_content"`
	UtmCampaign            string `json:"utm_campaign"`
	UtmTerm                string `json:"utm_term"`
	UtmMedium              string `json:"utm_medium"`
}

//PolygonGeoSegment - a segment based on geolocation data for polygons.
type PolygonGeoSegment struct {
	BaseSegment
	GeoSegmentType string   `json:"geo_segment_type"`
	TimesQuantity  int      `json:"times_quantity"`
	PeriodLength   int      `json:"period_length"`
	Polygons       []Points `json:"polygons"`
}

//Points - a group of points
type Points struct {
	Points []Point `json:"points"`
}

//Point - point's coordinates
type Point struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

//AppMetricaSegment - a segment imported from AppMetrica.
type AppMetricaSegment struct {
	BaseSegment
	AppMetricaSegmentType string `json:"app_metrica_segment_type"`
	AppMetricaSegmentID   int    `json:"app_metrica_segment_id"`
}

//CircleGeoSegment - segment based on circumferential geolocation data.
type CircleGeoSegment struct {
	BaseSegment
	GeoSegmentType string  `json:"geo_segment_type"`
	TimesQuantity  int     `json:"times_quantity"`
	PeriodLength   int     `json:"period_length"`
	Radius         int     `json:"radius"`
	Points         []Point `json:"points"`
}

//UploadingSegment - a segment created from a user data file.
type UploadingSegment struct {
	BaseSegment
	Hashed      bool   `json:"hashed"`
	ContentType string `json:"content_type"`
}

//LookalikeSegment - a segment from users who are “similar” to another segment of the client (Look-alike technology).
type LookalikeSegment struct {
	BaseSegment
	LookalikeLink              int64 `json:"lookalike_link"`
	LookalikeValue             int64 `json:"lookalike_value"`
	MaintainDeviceDistribution bool  `json:"maintain_device_distribution"`
	MaintainGeoDistribution    bool  `json:"maintain_geo_distribution"`
}

//MetrikaSegment - a segment imported from Yandex.Metrica.
type MetrikaSegment struct {
	BaseSegment
	MetrikaSegmentType string `json:"metrika_segment_type"`
	MetrikaSegmentID   int    `json:"metrika_segment_id"`
}

//SegmentsList - returns a list of existing segments available to the user.
func (c *Client) SegmentsList(pixel ...int) ([]map[string]interface{}, error) {
	requestPath := "segments"
	if len(pixel) > 0 {
		requestPath += fmt.Sprintf("?pixel=%d", pixel[0])
	}
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
	}, requestPath)
	if err != nil {
		return nil, err
	}
	var response struct {
		Segments []map[string]interface{} `json:"segments"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) != 0 {
		return nil, response.Error()
	}
	//TODO: Apply a method to separate segments by type and break them into structures
	//Problem: Problem: Different formats returns as one array. Bad fields-description in API documentation
	return response.Segments, nil
}

//CreateFileSegment - creates a segment from a data file. The file must have at least 1000 entries.
func (c *Client) CreateFileSegment(segment *UploadingSegment, filename string) error {
	var f *os.File
	var err error
	if f, err = os.Open(filename); err != nil {
		return err
	}
	defer closer(f)
	return c.CreateReaderSegment(segment, f, false)
}

//CreateCSVSegment - creates a segment from a csv data file. The file must have at least 1000 entries.
func (c *Client) CreateCSVSegment(segment *UploadingSegment, filename string) error {
	var f *os.File
	var err error
	if f, err = os.Open(filename); err != nil {
		return err
	}
	defer closer(f)
	return c.CreateReaderSegment(segment, f, true)
}

//CreateReaderSegment - creates a segment from a reader. The reader must have at least 1000 entries.
func (c *Client) CreateReaderSegment(segment *UploadingSegment, reader io.Reader, isCSV bool) error {
	rp, wp := io.Pipe()
	mpw := multipart.NewWriter(wp)
	errorChan := make(chan error, 1)
	go func() {
		var part io.Writer
		var err error
		defer closer(wp)
		if part, err = mpw.CreateFormFile("file", segment.Name); err != nil {
			errorChan <- err
			return
		}
		if _, err = io.Copy(part, reader); err != nil {
			errorChan <- err
		}
		if err = mpw.Close(); err != nil {
			errorChan <- err
		}
		errorChan <- nil
	}()
	URLPath := "upload_file"
	if isCSV {
		URLPath = "upload_csv_file"
	}
	req := http.Request{
		Method: http.MethodPost,
		Header: http.Header{
			"Content-Type": {mpw.FormDataContentType()},
		},
		Body: rp,
	}
	resp, err := c.Do(&req, fmt.Sprintf("segments/%s", URLPath))
	if err != nil {
		return err
	}
	defer closer(resp.Body)
	if err := <-errorChan; err != nil {
		return err
	}
	requestStruct := struct {
		Segment *UploadingSegment `json:"segment"`
		APIError
	}{Segment: segment}
	if err := json.NewDecoder(resp.Body).Decode(&requestStruct); err != nil {
		return err
	}
	if len(requestStruct.Errors) != 0 {
		return requestStruct.Error()
	}
	if segment.ID == 0 {
		return ErrNotCreated
	}
	return nil
}

//SaveUploadedSegment - saves a segment created from a data file.
func (c *Client) SaveUploadedSegment(segment *UploadingSegment) error {
	requestStruct := struct {
		Segment *UploadingSegment `json:"segment"`
		APIError
	}{Segment: segment}
	jsonBody, err := json.Marshal(requestStruct)
	if err != nil {
		return err
	}
	resp, err := c.Do(&http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, fmt.Sprintf("segment/%d/confirm?", segment.ID))
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

//RemoveSegment - deletes the specified segment.
func (c *Client) RemoveSegment(id int64) error {
	resp, err := c.Do(&http.Request{
		Method: http.MethodDelete,
	}, fmt.Sprintf("segment/%d", id))
	if err != nil {
		return err
	}
	defer closer(resp.Body)
	var respStruct struct {
		Success bool `json:"success"`
		APIError
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return err
	}
	if len(respStruct.Errors) != 0 {
		return respStruct.Error()
	}
	if !respStruct.Success {
		return ErrNotDeleted
	}
	return nil
}

//CreatePixelSegment - creates a segment of type "pixel" with the specified parameters.
//If different conditions are used when creating a segment (for example, several labels are specified),
//then a user who satisfies all the specified conditions at the same time will get into the segment.
func (c *Client) CreatePixelSegment(segment *PixelSegment) error {
	return c.createSegment(segment, "create_pixel")
}

//CreateLookalikeSegment - creates a “lookalike” type segment with the specified parameters.
func (c *Client) CreateLookalikeSegment(segment *LookalikeSegment) error {
	return c.createSegment(segment, "create_lookalike")
}

//CreateMetrikaSegment - creates a segment from a metric with the specified parameters.
func (c *Client) CreateMetrikaSegment(segment *MetrikaSegment) error {
	return c.createSegment(segment, "create_metrika")
}

//CreateAppMetrikaSegment - creates a segment from AppMetrica with the specified parameters.
func (c *Client) CreateAppMetrikaSegment(segment *AppMetricaSegment) error {
	return c.createSegment(segment, "create_appmetrica")
}

//CreateCircleGeoSegment - creates a segment based on geolocation data with the “circle” type.
func (c *Client) CreateCircleGeoSegment(segment *CircleGeoSegment) error {
	return c.createSegment(segment, "create_geo")
}

//CreatePolygonGeoSegment - creates a segment based on geolocation data with the “polygons” type.
func (c *Client) CreatePolygonGeoSegment(segment *PolygonGeoSegment) error {
	return c.createSegment(segment, "create_geo_polygon")
}

func (c *Client) createSegment(segment interface{}, URLPath string) error {
	requestStruct := struct {
		Segment interface{} `json:"segment"`
		APIError
	}{Segment: segment}
	jsonBody, err := json.Marshal(requestStruct)
	if err != nil {
		return err
	}
	resp, err := c.Do(&http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, "segments/"+URLPath)
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

//UpdateSegment - changes the specified segment.
func (c *Client) UpdateSegment(ID int64, segment interface{}) error {
	requestStruct := struct {
		Segment interface{} `json:"segment"`
		APIError
	}{Segment: segment}
	jsonBody, err := json.Marshal(requestStruct)
	if err != nil {
		return err
	}
	resp, err := c.Do(&http.Request{
		Method: http.MethodPut,
		Body:   ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
	}, fmt.Sprintf("segment/%d", ID))
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

//ReprocessSegment - starts a forced recount of a segment.
//Quotas for using the method: 2 requests per segment and 20 requests for user_login in the last 24 hours.
func (c *Client) ReprocessSegment(segmentID int64) error {
	resp, err := c.Do(&http.Request{
		Method: http.MethodPut,
	}, fmt.Sprintf("segment/%d/reprocess", segmentID))
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
		return ErrNotReprocessed
	}
	return nil
}
