package audience

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestClient_PixelsList(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("pixels list", t, func() {
		Convey("simple case", func(c C) {
			var data = []*Pixel{
				{ID: 2, Name: "name", UserQuantity7: 7, UserQuantity30: 30, UserQuantity90: 90},
				{ID: 4, Name: "name4", UserQuantity7: 14, UserQuantity30: 60, UserQuantity90: 180},
			}
			//set CreatedAt separately because time format
			data[1].CreateTime, _ = time.Parse(time.RFC3339, "2007-01-02T15:04:05Z07:00")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "pixels")
				_ = json.NewEncoder(w).Encode(struct {
					Pixels []*Pixel `json:"pixels"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			pixels, err := client.PixelsList()
			if err != nil {
				t.Fatal(err)
			}
			So(pixels, ShouldResemble, data)
		})
		Convey("zero results", func(c C) {
			var data = make([]Pixel, 0)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, "pixels")
				_ = json.NewEncoder(w).Encode(struct {
					Pixels []Pixel `json:"pixels"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			pixels, err := client.PixelsList()
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
				c.So(r.URL.Path, ShouldEndWith, "pixels")
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			_, err := client.PixelsList()
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreatePixel(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create pixel", t, func() {
		Convey("simple case", func(c C) {
			var data = Pixel{Name: "pixelname"}
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d struct {
					Pixel Pixel `json:"pixel"`
					APIError
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Pixel, ShouldResemble, data)
				c.So(r.URL.Path, ShouldEndWith, "pixels")
				isServerInvoked = true
				_ = json.NewEncoder(w).Encode(struct{}{})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreatePixel(&data); err != nil {
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
				c.So(r.URL.Path, ShouldEndWith, "pixels")
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreatePixel(&Pixel{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_RemovePixel(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var pixelID = int64(142)
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("remove pixel", t, func() {
		Convey("simple case", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d", pixelID))
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{true})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.RemovePixel(pixelID); err != nil {
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
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d", pixelID))
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemovePixel(pixelID)
			So(err, ShouldResemble, data.Error())
		})
		Convey("api return false", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d", pixelID))
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{false})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemovePixel(pixelID)
			So(err, ShouldEqual, ErrNotDeleted)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
}

func TestClient_UndeletePixel(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var pixelID = int64(142)
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("undelete pixel", t, func() {
		Convey("simple case", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d/undelete", pixelID))
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{true})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.UndeletePixel(pixelID); err != nil {
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
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d/undelete", pixelID))
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.UndeletePixel(pixelID)
			So(err, ShouldResemble, data.Error())
		})
		Convey("api return false", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d/undelete", pixelID))
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{false})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.UndeletePixel(pixelID)
			So(err, ShouldEqual, ErrNotRestored)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
}

func TestClient_UpdatePixel(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	pixelID := int64(132)
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("update pixel", t, func() {
		Convey("simple case", func(c C) {
			var data = Pixel{ID: pixelID, Name: "pixelname"}
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d struct {
					Pixel Pixel `json:"pixel"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Pixel, ShouldResemble, data)
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d", pixelID))
				isServerInvoked = true
				_ = json.NewEncoder(w).Encode(struct{}{})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.UpdatePixel(&data); err != nil {
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
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("pixel/%d", pixelID))
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.UpdatePixel(&Pixel{ID: pixelID})
			So(err, ShouldResemble, data.Error())
		})
	})
}
