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

func TestClient_DelegatesList(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("delegates list", t, func() {
		Convey("simple case", func() {
			var data = []*Delegate{{UserLogin: "guest", Perm: View, Comment: "1'st comment"}, {UserLogin: "editor", Perm: Edit, Comment: "2'nd comment"}}
			//set CreatedAt separately because time format
			data[0].CreatedAt, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
			data[1].CreatedAt, _ = time.Parse(time.RFC3339, "2007-01-02T15:04:05Z07:00")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(struct {
					Delegates []*Delegate `json:"delegates"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			delegates, err := client.DelegatesList()
			if err != nil {
				t.Fatal(err)
			}
			So(delegates, ShouldResemble, data)
		})
		Convey("zero results", func() {
			var data = make([]Delegate, 0)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(struct {
					Delegates []Delegate `json:"delegates"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			delegates, err := client.DelegatesList()
			if err != nil {
				t.Fatal(err)
			}
			So(len(delegates), ShouldEqual, 0)
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
			_, err := client.AccountsList()
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_CreateDelegate(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("create delegate", t, func() {
		Convey("simple case", func(c C) {
			var data = Delegate{UserLogin: "guest", Perm: View, Comment: "1'st comment"}
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var d struct {
					Delegate Delegate `json:"delegate"`
				}
				_ = json.NewDecoder(r.Body).Decode(&d)
				c.So(d.Delegate, ShouldResemble, data)
				isServerInvoked = true
				_ = json.NewEncoder(w).Encode(struct{}{})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.CreateDelegate(&data); err != nil {
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
			err := client.CreateDelegate(&Delegate{})
			So(err, ShouldResemble, data.Error())
		})
	})
}

func TestClient_RemoveDelegate(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var login = "my_login"
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("remove delegate", t, func() {
		Convey("simple case", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.String(), ShouldEndWith, login)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{true})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			if err := client.RemoveDelegate(login); err != nil {
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
			err := client.RemoveDelegate(login)
			So(err, ShouldResemble, data.Error())
		})
		Convey("api return false", func(c C) {
			isServerInvoked := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isServerInvoked = true
				c.So(r.URL.String(), ShouldEndWith, login)
				_ = json.NewEncoder(w).Encode(struct {
					Success bool `json:"success"`
				}{false})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			err := client.RemoveDelegate(login)
			So(err, ShouldEqual, ErrNotDeleted)
			So(isServerInvoked, ShouldBeTrue)
		})
	})
}
