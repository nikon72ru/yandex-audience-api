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
)

func TestNewClient(t *testing.T) {
	type key string
	var token = "blah"
	//var k = key("YANDEX_AUDIENCE_TOKEN")
	Convey("new client", t, func() {
		_ = os.Setenv(tokenVariable, "")
		Convey("token from context", func() {
			c, err := NewClient(context.WithValue(context.Background(), "YANDEX_AUDIENCE_TOKEN", token))
			So(err, ShouldBeNil)
			So(c.apiURL, ShouldEqual, apiURL)
			So(c.hc, ShouldNotBeNil)
			So(c.apiVersion, ShouldEqual, "v1")
			So(c.token, ShouldEqual, token)
		})
		Convey("token from envs", func() {
			_ = os.Setenv(tokenVariable, token)
			c, err := NewClient(context.Background())
			So(err, ShouldBeNil)
			So(c.apiURL, ShouldEqual, apiURL)
			So(c.hc, ShouldNotBeNil)
			So(c.apiVersion, ShouldEqual, "v1")
			So(c.token, ShouldEqual, token)
		})
		Convey("token isn't set", func() {
			c, err := NewClient(context.Background())
			So(err, ShouldEqual, ErrTokenIsNotSet)
			So(c, ShouldBeNil)
		})
	})
}

func TestClient_Do(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	var token = "blah"
	isServerInvoked := false
	client, _ := NewClient(context.Background())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isServerInvoked = true
		if r.Header.Get("Authorization") != fmt.Sprintf("OAuth %s", token) {
			t.Fatal("header isn't set")
		}
		if r.URL.Path != fmt.Sprintf("/%s/management/path", client.apiVersion) {
			t.Fatal(r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode("resp")
	}))
	defer ts.Close()
	client.hc = ts.Client()
	client.apiURL = ts.URL
	_, err := client.Do(&http.Request{
		Method: http.MethodGet,
	}, "path")
	if err != nil {
		t.Fatal(err)
	}
	if !isServerInvoked {
		t.Fatal("server wasn't called")
	}
}

func TestClient_Close(t *testing.T) {
	_ = os.Setenv(tokenVariable, "blah")
	cl, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if cl.hc == nil {
		t.Fatal("client shouldn't be nil")
	}
	if err := cl.Close(); err != nil {
		t.Fatal(err)
	}
	if cl.hc != nil {
		t.Fatal("client should be nil after Close() func")
	}
}
