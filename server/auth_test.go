package server

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/chrislewispac/rmd-server/models"
)

/*
func TestLogin(t *testing.T) {
	Convey("When you POST /login", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()
		Convey("And the Body is Empty", func() {
			Convey("You should get the right error message", func() {
				req, err := http.NewRequest(echo.POST, addr+"/login", strings.NewReader(""))
				So(err, ShouldEqual, nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(code, ShouldEqual, 200)
				error := body.Path("error").Data().(string)
				So(error, ShouldEqual, `No user was found`)
			})
		})

		Convey("And the user does not exist", func() {
			Convey("You should get the right error message", func() {
				req, err := http.NewRequest(echo.POST, addr+"/login", strings.NewReader(`{"email": "test@doesntexist.com", "password": "t"}`))
				So(err, ShouldEqual, nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				error := body.Path("error").Data().(string)
				pw := body.Path("data.user.password").Data().(string)
				token := body.Path("data.user.token").Data()
				So(error, ShouldEqual, "User Not Found")
				So(code, ShouldEqual, 200)
				So(pw, ShouldEqual, "")
				So(token, ShouldEqual, nil)
			})
		})

		Convey("And the user exists", func() {
			const email = "test@test.com"
			const pass = "test"
			hashedPw, err := ts.hashPass(pass)
			So(err, ShouldEqual, nil)
			err = ts.insertUser(email, hashedPw, new(Models.User))
			So(err, ShouldEqual, nil)

			Convey("And the password is incorrect", func() {
				body := `{"email": "` + email + `", "password": "wrong"}`
				Convey("You should get the right error message", func() {
					req, err := http.NewRequest(echo.POST, addr+"/login", strings.NewReader(body))
					So(err, ShouldEqual, nil)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, body := doRequest(req)

					error := body.Path("error").Data().(string)
					pw := body.Path("data.user.password").Data().(string)
					token := body.Path("data.user.token").Data()
					So(error, ShouldEqual, "Wrong Password")
					So(code, ShouldEqual, 200)
					So(pw, ShouldEqual, "")
					So(token, ShouldEqual, nil)
				})
			})

			Convey("And the password is correct", func() {
				body := `{"email": "` + email + `", "password": "` + pass + `"}`
				Convey("You should receive a response which includes a token", func() {
					req, err := http.NewRequest(echo.POST, addr+"/login", strings.NewReader(body))
					So(err, ShouldEqual, nil)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, body := doRequest(req)

					error := body.Path("error").Data()
					expectedResMsg := `Successfully Logged In`
					So(error, ShouldEqual, nil)
					So(code, ShouldEqual, 200)
					So(body.Path("msg").Data(), ShouldEqual, expectedResMsg)
					So(body.Path("data.user.password").Data(), ShouldEqual, "")
					So(body.Path("data.user.token").Data(), ShouldNotEqual, nil)
				})
			})
		})
	})
}
*/
func TestForgotPassword(t *testing.T) {
	Convey("When you POST to /forgot_password", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		Convey("and the email does not exist, you should get an error", func() {
			req, err := http.NewRequest(echo.POST, addr+"/forgot_password", strings.NewReader(`{"email": "doesnotexist@statrecruit.com"}`))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			error := body.Path("error").Data().(string)
			So(error, ShouldEqual, `That email is not associated with a user.`)
			So(code, ShouldEqual, 200)
		})

		Convey("and the email exists, you should get email confirmation", func(convey C) {
			err := ts.insertUser(testEmail, nil, new(Models.User))
			So(err, ShouldEqual, nil)

			req, err := http.NewRequest(echo.POST, addr+"/forgot_password", strings.NewReader(`{"email": "`+testEmail+`"}`))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			expectedResMsg := `Please check your email for a link to reset your password.`
			error := body.Path("error").Data()
			token := body.Path("token").Data()
			msg := body.Path("msg").Data().(string)
			pw := body.Path("data.user.password").Data().(string)
			So(error, ShouldEqual, nil)
			So(code, ShouldEqual, 200)
			So(token, ShouldEqual, nil)
			So(msg, ShouldEqual, expectedResMsg)
			So(pw, ShouldEqual, "")
		})

	})

}

// func TestResetPassword(t *testing.T) {
// 	Convey("When you POST to /reset_password", t, func() {
// 		e := echo.New()
// 		loadConfig()
// 		db := startTestServer(t)
//
// 		Convey("and the token is valid but passwords do not match, you should get an error", func() {
//
// 		})
//
// 		Convey("and the token has expired, you should get an error", func() {
// 			req, err := http.NewRequest(echo.POST, "/reset_password", strings.NewReader(`{"password": "anything", "password_confirm": "anything"}`))
// 			if err == nil {
// 				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 				rec := httptest.NewRecorder()
// 				c := e.NewContext(req, rec)
// 				ForgotPassword(db, Models.Redis)(c)
// 				expectedResMsg := `That email is not associated with a user.`
// 				body := MustParseJSON(rec.Body.String())
// 				error := body.Path("error").Data().(string)
// 				So(error, ShouldEqual, expectedResMsg)
// 				So(rec.Code, ShouldEqual, 200)
// 			}
// 		})
//
// 		Convey("and the token is valid and passwords match, response should be without err and state password has been updated", func() {
// 			req, err := http.NewRequest(echo.POST, "/forgot_password", strings.NewReader(`{"email": "`+testEmail+`"}`))
// 			if err == nil {
// 				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 				rec := httptest.NewRecorder()
// 				c := e.NewContext(req, rec)
// 				ForgotPassword(db, Models.Redis)(c)
// 				expectedResMsg := `Please check your email for a link to reset your password. This link will expire in 30 minutes!`
// 				body := MustParseJSON(rec.Body.String())
// 				error := body.Path("error").Data()
// 				token := body.Path("token").Data()
// 				msg := body.Path("msg").Data().(string)
// 				pw := body.Path("data.user.password").Data().(string)
// 				So(error, ShouldEqual, nil)
// 				So(rec.Code, ShouldEqual, 200)
// 				So(token, ShouldEqual, nil)
// 				So(msg, ShouldEqual, expectedResMsg)
// 				So(pw, ShouldEqual, "")
// 			}
// 		})
//
// 	})
//
// }

func TestRegister(t *testing.T) {
	ts, addr, close := startTestServer(t)
	defer close()
	Convey("When you post to /register", t, func() {
		Convey("and the email doesn't already exist you should get a response with a token", func() {
			u1 := uuid.NewV4()
			em := fmt.Sprintf("%s@test.com", u1)
			jsonObj := gabs.New()
			jsonObj.SetP(em, "email")
			jsonObj.SetP("test", "password")
			req, err := http.NewRequest(echo.POST, addr+"/register", strings.NewReader(jsonObj.String()))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			expectedResMsg := `Successfully Registered`
			error := body.Path("error").Data()
			msg := body.Path("msg").Data().(string)
			pw := body.Path("data.user.password").Data().(string)
			So(error, ShouldEqual, nil)
			So(code, ShouldEqual, 200)
			So(msg, ShouldEqual, expectedResMsg)
			So(pw, ShouldEqual, "")

		})
		Convey("and the email already exists you should get an err", func() {
			const email = "test@test.com"
			err := ts.insertUser(email, nil, new(Models.User))
			So(err, ShouldEqual, nil)

			jsonObj := gabs.New()
			jsonObj.SetP(email, "email")
			jsonObj.SetP("test", "password")
			req, err := http.NewRequest(echo.POST, addr+"/register", strings.NewReader(jsonObj.String()))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			expectedMsg := `Account already exists with email test@test.com`
			error := body.Path("error").Data()
			pw := body.Path("data.user.password").Data().(string)
			So(error, ShouldEqual, expectedMsg)
			So(code, ShouldEqual, 200)
			So(pw, ShouldEqual, "")
		})

	})

}

func TestLogout(t *testing.T) {
	ts, addr, close := startTestServer(t)
	defer close()

	Convey("When you post to /logout", t, func() {
		url := addr + "/auth/logout"

		Convey("with a missing token, you should get an error", func() {
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(""))
			So(err, ShouldEqual, nil)

			code, _ := doRequest(req)

			So(code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("with an invalid token, you should get an error", func() {
			const tStr = "ASSDDF-invalid-token-TEST_ADSF"

			req, err := http.NewRequest(echo.POST, url, strings.NewReader(""))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))

			code, _ := doRequest(req)

			So(code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("with a valid token, that is in storage, you should get a success msg and be logged out", func() {
			tokenWrong := generateToken("test", 10*time.Minute, false)
			tStr, err := tokenWrong.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			req, err := http.NewRequest(echo.POST, url, strings.NewReader(""))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))

			code, body := doRequest(req)

			error := body.Path("error").Data()
			msg := body.Path("msg").Data().(string)
			pw := body.Path("data.user.password").Data().(string)
			token := body.Path("data.user.token").Data()
			So(error, ShouldEqual, nil)
			So(msg, ShouldEqual, "Logged Out")
			So(code, ShouldEqual, 200)
			So(pw, ShouldEqual, "")
			So(token, ShouldEqual, "")

		})

		Convey("with a valid token, that is not in storage, you should get err", func() {
			token := generateToken("notInStorage", 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			req, err := http.NewRequest(echo.POST, url, strings.NewReader(""))
			So(err, ShouldEqual, nil)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))

			code, _ := doRequest(req)

			So(code, ShouldEqual, http.StatusUnauthorized)
		})

	})

}
