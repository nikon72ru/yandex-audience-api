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

func TestClient_GrantsList(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var segmentID = int64(142)
	Convey("grants list", t, func() {
		Convey("simple case", func(c C) {
			var data = []*Grant{{UserLogin: "guest", Comment: "1'st comment"}, {UserLogin: "editor", Comment: "2'nd comment"}}
			//set CreatedAt separately because time format
			data[0].CreatedAt, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
			data[1].CreatedAt, _ = time.Parse(time.RFC3339, "2007-01-02T15:04:05Z07:00")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("/segment/%d/grants", segmentID))
				_ = json.NewEncoder(w).Encode(struct {
					Grants []*Grant `json:"grants"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			grants, err := client.GrantsList(segmentID)
			if err != nil {
				t.Fatal(err)
			}
			So(grants, ShouldResemble, data)
		})
		Convey("zero results", func(c C) {
			var data = make([]Grant, 0)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("/segment/%d/grants", segmentID))
				_ = json.NewEncoder(w).Encode(struct {
					Grants []Grant `json:"grants"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			grants, err := client.GrantsList(segmentID)
			if err != nil {
				t.Fatal(err)
			}
			So(len(grants), ShouldEqual, 0)
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
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("/segment/%d/grants", segmentID))
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			_, err := client.GrantsList(segmentID)
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreateGrant(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var segmentID = int64(142)
	Convey("create grant", t, func() {
		Convey("simple case", func(c C) {
			var data = Grant{UserLogin: "guest", Comment: "1'st comment"}
			data.CreatedAt, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d Grant
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d, ShouldResemble, data)
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/grant", segmentID))
				isServerInvoked = true
				_ = json.NewEncoder(w).Encode(struct{}{})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreateGrant(segmentID, &data); err != nil {
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
				c.So(r.URL.Path, ShouldEndWith, fmt.Sprintf("segment/%d/grant", segmentID))
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.CreateGrant(segmentID, &Grant{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_RemoveGrant(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var login = "my_login"
	var segmentID = int64(142)
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("remove grant", t, func() {
		Convey("simple case", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, login)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{true})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.RemoveGrant(segmentID, login); err != nil {
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
				c.So(r.URL.Path, ShouldEndWith, login)
				_ = json.NewEncoder(w).Encode(data)
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemoveGrant(segmentID, login)
			So(err, ShouldResemble, data.Error())
		})
		Convey("api return false", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.Path, ShouldEndWith, login)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{false})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemoveGrant(segmentID, login)
			So(err, ShouldEqual, ErrNotDeleted)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
}
