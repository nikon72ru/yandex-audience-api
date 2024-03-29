package audience

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
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
	points := Points{Points: []Point{{
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
	}}}
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
					Polygons: []Points{
						points,
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
				var d struct {
					Segment AppMetricaSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				d.Segment.CreateTime = createdTime
				_ = json.NewEncoder(w).Encode(d)
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
				var d struct {
					Segment CircleGeoSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				d.Segment.CreateTime = createdTime
				_ = json.NewEncoder(w).Encode(d)
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
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create file segment", t, func() {
		Convey("check file integrity", func(c C) {
			filepath := "../test-files/macs_for_uploads.csv"
			f, err := os.Open(filepath)
			So(err, ShouldBeNil)
			data, err := ioutil.ReadAll(f)
			So(err, ShouldBeNil)
			_ = f.Close()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				err := r.ParseMultipartForm(32 << 20)
				c.So(err, ShouldBeNil)
				files := r.MultipartForm.File["file"]
				c.So(len(files), ShouldEqual, 1)
				file := files[0]
				f, err := file.Open()
				c.So(err, ShouldBeNil)
				receivedData, err := ioutil.ReadAll(f)
				c.So(err, ShouldBeNil)
				c.So(receivedData, ShouldResemble, data)
				c.So(r.URL.Path, ShouldEndWith, "segments/upload_file")
				_ = json.NewEncoder(w).Encode(struct {
					Segment UploadingSegment `json:"segment"`
				}{UploadingSegment{
					BaseSegment: BaseSegment{ID: 12},
				}})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			var segment = UploadingSegment{
				BaseSegment: BaseSegment{Name: "file segment"},
			}
			err = client.CreateFileSegment(&segment, filepath)
			So(err, ShouldBeNil)
			So(segment.ID, ShouldEqual, 12)
			So(isServerInvoked, ShouldBeTrue)
		})
	})

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
				var d struct {
					Segment LookalikeSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				d.Segment.CreateTime = createdTime
				_ = json.NewEncoder(w).Encode(d)
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
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create csv segment", t, func() {
		Convey("check file integrity", func(c C) {
			filepath := "../test-files/macs_for_uploads.csv"
			f, err := os.Open(filepath)
			So(err, ShouldBeNil)
			data, err := ioutil.ReadAll(f)
			So(err, ShouldBeNil)
			_ = f.Close()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				err := r.ParseMultipartForm(32 << 20)
				c.So(err, ShouldBeNil)
				files := r.MultipartForm.File["file"]
				c.So(len(files), ShouldEqual, 1)
				file := files[0]
				f, err := file.Open()
				c.So(err, ShouldBeNil)
				receivedData, err := ioutil.ReadAll(f)
				c.So(err, ShouldBeNil)
				c.So(receivedData, ShouldResemble, data)
				c.So(r.URL.Path, ShouldEndWith, "segments/upload_csv_file")
				_ = json.NewEncoder(w).Encode(struct {
					Segment UploadingSegment `json:"segment"`
				}{UploadingSegment{
					BaseSegment: BaseSegment{ID: 12},
				}})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			var segment = UploadingSegment{
				BaseSegment: BaseSegment{Name: "segment from csv file"},
			}
			err = client.CreateCSVSegment(&segment, filepath)
			So(err, ShouldBeNil)
			So(segment.ID, ShouldEqual, 12)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
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
				var d struct {
					Segment MetrikaSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				d.Segment.CreateTime = createdTime
				c.So(r.URL.Path, ShouldEndWith, "create_metrika")
				_ = json.NewEncoder(w).Encode(d)
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
				var d struct {
					Segment PixelSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				d.Segment.CreateTime = createdTime
				c.So(r.URL.Path, ShouldEndWith, "create_pixel")
				_ = json.NewEncoder(w).Encode(d)
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
	points := Points{Points: []Point{{
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
	}}}
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
				Polygons:       []Points{points},
			}
			var createdTime = time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d struct {
					Segment PolygonGeoSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				d.Segment.CreateTime = createdTime
				c.So(r.URL.Path, ShouldEndWith, "create_geo_polygon")
				_ = json.NewEncoder(w).Encode(d)
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
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create reader segment", t, func() {
		Convey("check file integrity", func(c C) {
			f, err := os.Open("../test-files/macs_for_uploads.csv")
			So(err, ShouldBeNil)
			data, err := ioutil.ReadAll(f)
			So(err, ShouldBeNil)
			_ = f.Close()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				err := r.ParseMultipartForm(32 << 20)
				c.So(err, ShouldBeNil)
				files := r.MultipartForm.File["file"]
				c.So(len(files), ShouldEqual, 1)
				file := files[0]
				f, err := file.Open()
				c.So(err, ShouldBeNil)
				receivedData, err := ioutil.ReadAll(f)
				c.So(err, ShouldBeNil)
				c.So(receivedData, ShouldResemble, data)
				c.So(r.URL.Path, ShouldEndWith, "segments/upload_file")
				_ = json.NewEncoder(w).Encode(struct {
					Segment UploadingSegment `json:"segment"`
				}{UploadingSegment{
					BaseSegment: BaseSegment{ID: 12},
				}})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			var segment = UploadingSegment{
				BaseSegment: BaseSegment{Name: "test segment"},
			}
			err = client.CreateReaderSegment(&segment, bytes.NewBuffer(data), false)
			So(err, ShouldBeNil)
			So(segment.ID, ShouldEqual, 12)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
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
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("save uploaded segment", t, func() {
		Convey("simple case", func(c C) {
			segment := UploadingSegment{
				BaseSegment: BaseSegment{
					ID:   12,
					Name: "segmentname",
				},
				Hashed:      true,
				ContentType: Mac,
			}
			creationTime := time.Now()
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var s struct {
					Segment UploadingSegment `json:"segment"`
				}
				err = json.NewDecoder(r.Body).Decode(&s)
				c.So(err, ShouldBeNil)
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/confirm", segment.ID))
				c.So(r.Method, ShouldEqual, http.MethodPost)
				c.So(s.Segment, ShouldResemble, segment)
				s.Segment.CreateTime = creationTime
				_ = json.NewEncoder(w).Encode(struct {
					Segment *UploadingSegment `json:"segment"`
				}{&s.Segment})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.SaveUploadedSegment(&segment); err != nil {
				So(err, ShouldBeNil)
			}
			So(isServerInvoked, ShouldBeTrue)
			So(segment.CreateTime, ShouldEqual, creationTime)
		})
		Convey("api return error", func(c C) {
			segmentID := int64(13)
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
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/confirm", segmentID))
				c.So(r.Method, ShouldEqual, http.MethodPost)
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.SaveUploadedSegment(&UploadingSegment{BaseSegment: BaseSegment{ID: segmentID}})
			So(err, ShouldResemble, data.Error())
		})
	})
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
				var d struct {
					Segment AppMetricaSegment `json:"segment"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Segment, ShouldResemble, data)
				isServerInvoked = true
				_ = json.NewEncoder(w).Encode(d)
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
