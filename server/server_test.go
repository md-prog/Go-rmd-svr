package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/jinzhu/configor"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/chrislewispac/rmd-server/models"
	"io/ioutil"
)

const (
	testEmail = "chrislewispac+rmdtest@gmail.com"
)

// A TestConfig wraps a Config in order to use the same settings but with
// distinct keys.
type TestConfig struct {
	Test Config
}

// startTestServer returns a new *empty* test server instance, its address, and
// a close function to stop serving.
func startTestServer(t *testing.T) (*Server, string, func()) {
	c := &TestConfig{}
	if _, err := os.Stat("../config.test.yml"); err == nil {
		t.Log("loading config.test.yml...")
		configor.Load(c, "../config.test.yml")
	} else {
		t.Log("no config.test.yml found")
		configor.Load(c)
	}

	s, err := NewServer(&c.Test)
	if err != nil {
		t.Fatal("failed to create server:", err)
	}

	if err := s.truncate("users", "facilities", "contracts", "people", "emrs"); err != nil {
		t.Fatal("failed to truncate tables:", err)
	}

	e := echo.New()
	s.InitRoutes(e)

	httpSrv := httptest.NewServer(e)
	return s, httpSrv.URL, func() {
		httpSrv.Close()
		s.Close()
	}
}

// truncate executes a cascading truncate against each table.
func (s *Server) truncate(tables ...string) error {
	for _, table := range tables {
		_, err := s.db.Exec("TRUNCATE " + table + " CASCADE;")
		if err != nil {
			return err
		}
	}
	return nil
}

// doRequest executes a request in the context of a goconvey test and returns
// the status code and gabs Container body if successful.
func doRequest(req *http.Request) (int, *gabs.Container) {
	resp, err := http.DefaultClient.Do(req)
	So(err, ShouldEqual, nil)
	b, err := ioutil.ReadAll(resp.Body)
	So(err, ShouldEqual, nil)
	So(resp.Body.Close(), ShouldEqual, nil)
	body, err := gabs.ParseJSON(b)
	So(err, ShouldEqual, nil)

	return resp.StatusCode, body
}

// doRequest executes a request in the context of a goconvey test and returns
// the status code and Models.Res body if successful.
func doRequestRes(req *http.Request) (int, *Models.Res) {
	resp, err := http.DefaultClient.Do(req)
	So(err, ShouldEqual, nil)
	b, err := ioutil.ReadAll(resp.Body)
	So(err, ShouldEqual, nil)
	So(resp.Body.Close(), ShouldEqual, nil)
	var res Models.Res
	So(json.Unmarshal(b, &res), ShouldEqual, nil)

	return resp.StatusCode, &res
}

// doRequest executes a request in the context of a goconvey test and returns
// the status code and gabs Container body if successful.
func doRequestRaw(req *http.Request) (int, []byte) {
	resp, err := http.DefaultClient.Do(req)
	So(err, ShouldEqual, nil)
	b, err := ioutil.ReadAll(resp.Body)
	So(err, ShouldEqual, nil)
	So(resp.Body.Close(), ShouldEqual, nil)

	return resp.StatusCode, b
}
