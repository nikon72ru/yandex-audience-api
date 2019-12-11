package audience

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestClient_SegmentsList(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("segments list", t, func() {
		Convey("simple case", func(c C) {
			var data = []interface{}{
				PixelSegment{
					BaseSegment: BaseSegment{
						ID:         30,
						Name:       "pixel segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					PixelID:                130,
					PeriodLength:           90,
					TimesQuantity:          10,
					TimesQuantityOperation: "eq",
					UtmSource:              "utm_source",
					UtmContent:             "utm_content",
					UtmCampaign:            "utm_campaign",
					UtmTerm:                "utm_term",
					UtmMedium:              "utm_medium",
				},
				LookalikeSegment{
					BaseSegment: BaseSegment{
						ID:         40,
						Name:       "lookalike segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					LookalikeLink:              12,
					LookalikeValue:             1,
					MaintainDeviceDistribution: true,
					MaintainGeoDistribution:    false,
				},
				MetrikaSegment{
					BaseSegment: BaseSegment{
						ID:         50,
						Name:       "metrika_segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					MetrikaSegmentType: GoalID,
					MetrikaSegmentID:   123,
				},
				AppMetricaSegment{
					BaseSegment: BaseSegment{
						ID:         60,
						Name:       "app metrika segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					AppMetricaSegmentType: "api_key",
					AppMetricaSegmentID:   124,
				},
				CircleGeoSegment{
					BaseSegment: BaseSegment{
						ID:         70,
						Name:       "circle geo segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					GeoSegmentType: "condition",
					TimesQuantity:  3,
					PeriodLength:   30,
					Radius:         300,
					Points: []Point{
						{
							Latitude:    65,
							Longitude:   67,
							Description: "center",
						},
					},
				},
				PolygonGeoSegment{
					BaseSegment: BaseSegment{
						ID:         80,
						Name:       "polygon geo segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					GeoSegmentType: "work",
					TimesQuantity:  3,
					PeriodLength:   13,
					Polygons: [][]Point{
						{
							{
								Latitude:    65,
								Longitude:   65,
								Description: "left bottom",
							}, {
								Latitude:    66,
								Longitude:   65,
								Description: "right bottom",
							}, {
								Latitude:    65,
								Longitude:   66,
								Description: "left up",
							}, {
								Latitude:    66,
								Longitude:   66,
								Description: "right up",
							},
						},
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "segments")
				_ = json.NewEncoder(w).Encode(struct {
					Segments []interface{} `json:"segments"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			segments, err := client.SegmentsList()
			if err != nil {
				t.Fatal(err)
			}
			marshaled, _ := json.Marshal(data)
			var dataMap []map[string]interface{}
			_ = json.Unmarshal(marshaled, &dataMap)
			So(segments, ShouldResemble, dataMap)
		})
		Convey("zero results", func(c C) {
			var data = make([]*interface{}, 0)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "segments")
				_ = json.NewEncoder(w).Encode(struct {
					Segments []*interface{} `json:"segments"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			pixels, err := client.SegmentsList()
			if err != nil {
				t.Fatal(err)
			}
			So(len(pixels), ShouldEqual, 0)
		})
		Convey("api return error", func(c C) {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "segments")
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			_, err := client.SegmentsList()
			So(err, ShouldResemble, data.Error())
		})
		Convey("with pixel", func(c C) {
			pixelID := 123
			var data = []interface{}{
				PixelSegment{
					BaseSegment: BaseSegment{
						ID:         30,
						Name:       "pixel segment",
						Status:     "uploaded",
						CreateTime: time.Time{},
						Owner:      "lva",
					},
					PixelID:                130,
					PeriodLength:           90,
					TimesQuantity:          10,
					TimesQuantityOperation: "eq",
					UtmSource:              "utm_source",
					UtmContent:             "utm_content",
					UtmCampaign:            "utm_campaign",
					UtmTerm:                "utm_term",
					UtmMedium:              "utm_medium",
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "segments")
				_ = r.ParseForm()
				c.So(r.FormValue("pixel"), ShouldEqual, strconv.Itoa(pixelID))
				_ = json.NewEncoder(w).Encode(struct {
					Segments []interface{} `json:"segments"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			segments, err := client.SegmentsList(pixelID)
			if err != nil {
				t.Fatal(err)
			}
			marshaled, _ := json.Marshal(data)
			var dataMap []map[string]interface{}
			_ = json.Unmarshal(marshaled, &dataMap)
			So(segments, ShouldResemble, dataMap)
		})
	})
}

func TestClient_CreateAppMetrikaSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create app metrika segment", t, func() {
		Convey("simple case", func(c C) {
			var data = AppMetricaSegment{
				BaseSegment: BaseSegment{
					ID:         10,
					Name:       "app metrika segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				AppMetricaSegmentType: "segment_id",
				AppMetricaSegmentID:   30,
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d AppMetricaSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				d.CreateTime = createdTime
				_ = json.NewEncoder(w).Encode(struct {
					Segment AppMetricaSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreateAppMetrikaSegment(&data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(data.CreateTime, ShouldEqual, createdTime)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateAppMetrikaSegment(&AppMetricaSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreateCircleGeoSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create circle geo segment", t, func() {
		Convey("simple case", func(c C) {
			var data = CircleGeoSegment{
				BaseSegment: BaseSegment{
					ID:         70,
					Name:       "circle geo segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				GeoSegmentType: "condition",
				TimesQuantity:  3,
				PeriodLength:   30,
				Radius:         300,
				Points: []Point{
					{
						Latitude:    65,
						Longitude:   67,
						Description: "center",
					},
				},
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "segments/create_geo")
				var d CircleGeoSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				d.CreateTime = createdTime
				_ = json.NewEncoder(w).Encode(struct {
					Segment CircleGeoSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreateCircleGeoSegment(&data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(data.CreateTime, ShouldEqual, createdTime)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateCircleGeoSegment(&CircleGeoSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreateFileSegment(t *testing.T) {
	t.Fatal("implement me")
}

func TestClient_CreateLookalikeSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create lookalike segment", t, func() {
		Convey("simple case", func(c C) {
			var data = LookalikeSegment{
				BaseSegment: BaseSegment{
					ID:         40,
					Name:       "lookalike segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				LookalikeLink:              12,
				LookalikeValue:             1,
				MaintainDeviceDistribution: true,
				MaintainGeoDistribution:    false,
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "segments/create_lookalike")
				var d LookalikeSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				d.CreateTime = createdTime
				_ = json.NewEncoder(w).Encode(struct {
					Segment LookalikeSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreateLookalikeSegment(&data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(data.CreateTime, ShouldEqual, createdTime)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateLookalikeSegment(&LookalikeSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreateCSVSegment(t *testing.T) {
	t.Fatal("implement me!")
}

func TestClient_CreateMetrikaSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create metrika segment", t, func() {
		Convey("simple case", func(c C) {
			var data = MetrikaSegment{
				BaseSegment: BaseSegment{
					ID:         50,
					Name:       "metrika_segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				MetrikaSegmentType: GoalID,
				MetrikaSegmentID:   123,
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d MetrikaSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				d.CreateTime = createdTime
				c.So(r.URL.Path, ShouldEndWith, "create_metrika")
				_ = json.NewEncoder(w).Encode(struct {
					Segment MetrikaSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreateMetrikaSegment(&data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(data.CreateTime, ShouldEqual, createdTime)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateMetrikaSegment(&MetrikaSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreatePixelSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create pixel segment", t, func() {
		Convey("simple case", func(c C) {
			var data = PixelSegment{
				BaseSegment: BaseSegment{
					ID:         30,
					Name:       "pixel segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				PixelID:                130,
				PeriodLength:           90,
				TimesQuantity:          10,
				TimesQuantityOperation: "eq",
				UtmSource:              "utm_source",
				UtmContent:             "utm_content",
				UtmCampaign:            "utm_campaign",
				UtmTerm:                "utm_term",
				UtmMedium:              "utm_medium",
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d PixelSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				d.CreateTime = createdTime
				c.So(r.URL.Path, ShouldEndWith, "create_pixel")
				_ = json.NewEncoder(w).Encode(struct {
					Segment PixelSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreatePixelSegment(&data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(data.CreateTime, ShouldEqual, createdTime)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreatePixelSegment(&PixelSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreatePolygonGeoSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create polygon geo segment", t, func() {
		Convey("simple case", func(c C) {
			var data = PolygonGeoSegment{
				BaseSegment: BaseSegment{
					ID:         80,
					Name:       "polygon geo segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				GeoSegmentType: "work",
				TimesQuantity:  3,
				PeriodLength:   13,
				Polygons: [][]Point{
					{
						{
							Latitude:    65,
							Longitude:   65,
							Description: "left bottom",
						}, {
							Latitude:    66,
							Longitude:   65,
							Description: "right bottom",
						}, {
							Latitude:    65,
							Longitude:   66,
							Description: "left up",
						}, {
							Latitude:    66,
							Longitude:   66,
							Description: "right up",
						},
					},
				},
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d PolygonGeoSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				d.CreateTime = createdTime
				c.So(r.URL.Path, ShouldEndWith, "create_geo_polygon")
				_ = json.NewEncoder(w).Encode(struct {
					Segment PolygonGeoSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreatePolygonGeoSegment(&data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(data.CreateTime, ShouldEqual, createdTime)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateAppMetrikaSegment(&AppMetricaSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreateReaderSegment(t *testing.T) {
	t.Fatal("implement me!")
}

func TestClient_RemoveSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var segmentID = int64(142)
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("remove segment", t, func() {
		Convey("simple case", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodDelete)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{true})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.RemoveSegment(segmentID); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
		})
		Convey("api return error", func(c C) {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodDelete)
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemoveSegment(segmentID)
			So(err, ShouldResemble, data.Error())
		})
		Convey("api return false", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodDelete)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{false})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemoveSegment(segmentID)
			So(err, ShouldEqual, ErrNotDeleted)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
}

func TestClient_ReprocessSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var segmentID = int64(142)
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("reprocess segment", t, func() {
		Convey("simple case", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/reprocess", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodPut)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{true})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.ReprocessSegment(segmentID); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
		})
		Convey("api return error", func(c C) {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/reprocess", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodPut)
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.ReprocessSegment(segmentID)
			So(err, ShouldResemble, data.Error())
		})
		Convey("api return false", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/reprocess", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodPut)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{false})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.ReprocessSegment(segmentID)
			So(err, ShouldEqual, ErrNotReprocessed)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
}

func TestClient_SaveUploadedSegment(t *testing.T) {
	t.Fatal("implement me!")
}

func TestClient_UpdateSegment(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("update segment", t, func() {
		Convey("simple case", func(c C) {
			var data = AppMetricaSegment{
				BaseSegment: BaseSegment{
					ID:         10,
					Name:       "app metrika segment",
					Status:     "uploaded",
					CreateTime: time.Time{},
					Owner:      "lva",
				},
				AppMetricaSegmentType: "segment_id",
				AppMetricaSegmentID:   30,
			}
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d AppMetricaSegment
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				isServerInvoked = true
				_ = json.NewEncoder(w).Encode(struct {
					Segment AppMetricaSegment `json:"segment"`
				}{d})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.UpdateSegment(data.ID, &data); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
		})
		Convey("api return error", func() {
			var data = APIError{
				Errors: []Error{{
					ErrorType: "backend_error",
					Message:   "simple error",
					Location:  "right here",
				}},
				Code:    503,
				Message: "simple error",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateAppMetrikaSegment(&AppMetricaSegment{})
			So(err, ShouldResemble, data.Error())
		})
	})
}
