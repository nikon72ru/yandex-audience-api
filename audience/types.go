package audience

import "time"

type BaseSegment struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
	Owner      string    `json:"owner"`
}

type PixelSegment struct {
	BaseSegment
	PixelId                int    `json:"pixel_id"`
	PeriodLength           int    `json:"period_length"`
	TimesQuantity          int    `json:"times_quantity"`
	TimesQuantityOperation string `json:"times_quantity_operation"`
	UtmSource              string `json:"utm_source"`
	UtmContent             string `json:"utm_content"`
	UtmCampaign            string `json:"utm_campaign"`
	UtmTerm                string `json:"utm_term"`
	UtmMedium              string `json:"utm_medium"`
}

type Point struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type PolygonGeoSegment struct {
	BaseSegment
	GeoSegmentType string    `json:"geo_segment_type"`
	TimesQuantity  int       `json:"times_quantity"`
	PeriodLength   int       `json:"period_length"`
	Polygons       [][]Point `json:"polygons"`
}

type AppMetricaSegment struct {
	BaseSegment
	AppMetricaSegmentType string `json:"app_metrica_segment_type"`
	AppMetricaSegmentId   int    `json:"app_metrica_segment_id"`
}

type CircleGeoSegment struct {
	BaseSegment
	GeoSegmentType string `json:"geo_segment_type"`
	TimesQuantity  int    `json:"times_quantity"`
	PeriodLength   int    `json:"period_length"`
	Radius         int    `json:"radius"`
	Points         []Point
}

type UploadingSegment struct {
	BaseSegment
	Hashed      bool   `json:"hashed"`
	ContentType string `json:"content_type"`
}

type LookalikeSegment struct {
	BaseSegment
	LookalikeLink              int  `json:"lookalike_link"`
	LookalikeValue             int  `json:"lookalike_value"`
	MaintainDeviceDistribution bool `json:"maintain_device_distribution"`
	MaintainGeoDistribution    bool `json:"maintain_geo_distribution"`
}

type MetrikaSegment struct {
	BaseSegment
	MetrikaSegmentType string `json:"metrika_segment_type"`
	MetrikaSegmentId   int    `json:"metrika_segment_id"`
}

type Segment interface{}
