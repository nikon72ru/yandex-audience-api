package audience

import (
	"context"
	"encoding/json"
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
			var data = []Pixel{
				{ID: 2, Name: "name", UserQuantity7: 7, UserQuantity30: 30, UserQuantity90: 90},
				{ID: 4, Name: "name4", UserQuantity7: 14, UserQuantity30: 60, UserQuantity90: 180},
			}
			//set CreatedAt separately because time format
			data[1].CreateTime, _ = time.Parse(time.RFC3339, "2007-01-02T15:04:05Z07:00")
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
			So(*pixels, ShouldResemble, data)
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
			So(len(*pixels), ShouldEqual, 0)
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
	t.Fatal("not implemented")
}

func TestClient_RemovePixel(t *testing.T) {
	t.Fatal("not implemented")
}

func TestClient_UndeletePixel(t *testing.T) {
	t.Fatal("not implemented")
}

func TestClient_UpdatePixel(t *testing.T) {
	t.Fatal("not implemented")
}
