package Controllers

import (
	_ "fmt"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	userJSON  = `{"name":"Jon Snow","email":"jon@labstack.com"}`
	userEmail = `{"email":"jon@labstack.com"}`
)

func TestLogin(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("When you POST /login", t, func() {
		Convey("And the Body is Empty", func() {
			Convey("You should get the right error message", func() {
				e := echo.New()
				req, err := http.NewRequest(echo.POST, "/login", strings.NewReader(""))
				if err == nil {
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					c := e.NewContext(req, rec)
					LoginCtrl(c)
					So(rec.Code, ShouldEqual, 500)
					expectedResMsg := `"There was nothing submitted"`
					So(rec.Body.String(), ShouldEqual, expectedResMsg)
				}
			})
		})
	})
}

func TestPWReset(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("When you POST /pw_reset", t, func() {
		Convey("And the email exists", func() {
			Convey("It should send an email with a password link", func() {
				e := echo.New()
				req, err := http.NewRequest(echo.POST, "/pw_reset", strings.NewReader(""))
				if err == nil {
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					c := e.NewContext(req, rec)
					ForgotPasswordCtrl(c)
					So(rec.Code, ShouldEqual, 200)
					expectedResMsg := `"password reset email sent"`
					So(rec.Body.String(), ShouldEqual, expectedResMsg)
				}
			})
		})
		Convey("And the email does not exist", func() {
			Convey("It should send an email with a password link", func() {
			})
		})
	})
}
