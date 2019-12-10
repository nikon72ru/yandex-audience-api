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

func TestClient_AccountsList(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Convey("account list", t, func() {
		Convey("simple case", func() {
			var data = []Account{{UserLogin: "guest", Perm: View}, {UserLogin: "editor", Perm: Edit}}
			//set CreatedAt separately because time format
			data[0].CreatedAt, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
			data[1].CreatedAt, _ = time.Parse(time.RFC3339, "2007-01-02T15:04:05Z07:00")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(struct {
					Accounts []Account `json:"accounts"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			accounts, err := client.AccountsList()
			if err != nil {
				t.Fatal(err)
			}
			So(*accounts, ShouldResemble, data)
		})
		Convey("zero results", func() {
			var data = make([]Account, 0)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(struct {
					Accounts []Account `json:"accounts"`
				}{data})
			}))
			defer ts.Close()
			client.hc = ts.Client()
			client.apiURL = ts.URL
			accounts, err := client.AccountsList()
			if err != nil {
				t.Fatal(err)
			}
			So(len(*accounts), ShouldEqual, 0)
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
