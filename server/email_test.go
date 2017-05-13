package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

//TestEmailReceiver tests
func TestEmailReceiver(t *testing.T) {
	_, addr, close := startTestServer(t)
	defer close()

	Convey("When sparkpost posts to /email_forwarding", t, func() {
		url := addr + "/email_forwarding"
		Convey("and the handle has an associated user", func() {
			file, err := ioutil.ReadFile("./test_files/email_test.json")
			if err != nil {
				fmt.Printf("File error: %v\n", err)
				os.Exit(1)
			}
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(string(file)))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			error := body.Path("error").Data()
			So(error, ShouldEqual, nil)
			So(code, ShouldEqual, 200)
		})
		Convey("and the handle does not have an associated user", func() {
			file, err := ioutil.ReadFile("./test_files/email_test_bad.json")
			if err != nil {
				fmt.Printf("File error: %v\n", err)
				os.Exit(1)
			}
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(string(file)))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			error := body.Path("error").Data()
			So(error, ShouldEqual, nil)
			So(code, ShouldEqual, 200)
		})

	})

}
