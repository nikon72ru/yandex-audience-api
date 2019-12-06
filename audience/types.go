package audience

import (
	"encoding/json"
	"fmt"
	"time"
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
	GeoSegmentType string    `json:"geo_segment_type"`
	TimesQuantity  int       `json:"times_quantity"`
	PeriodLength   int       `json:"period_length"`
	Polygons       [][]Point `json:"polygons"`
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
	GeoSegmentType string `json:"geo_segment_type"`
	TimesQuantity  int    `json:"times_quantity"`
	PeriodLength   int    `json:"period_length"`
	Radius         int    `json:"radius"`
	Points         []Point
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

//Segments - a group of segments
type Segments struct {
	MetrikaSegments    []MetrikaSegment
	LookalikeSegments  []LookalikeSegment
	UploadingSegments  []UploadingSegment
	CircleGeoSegments  []CircleGeoSegment
	AppMetricaSegments []AppMetricaSegment
	PolygonGeoSegments []PolygonGeoSegment
	PixelSegments      []PixelSegment
	UnknownSegmtns     []BaseSegment
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
