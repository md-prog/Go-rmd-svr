package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"

	"gitlab.com/chrislewispac/rmd-server/models"
	"gopkg.in/guregu/null.v3"
)

func TestServer_GetUserProviders(t *testing.T) {
	Convey("when you GET /auth/providers", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/providers"

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("with 0 providers, you should get an empty list", func() {
				req, err := http.NewRequest(echo.GET, url, nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, res := doRequestRes(req)

				So(code, ShouldEqual, 200)
				So(res.Error.String, ShouldEqual, "")
				So(res.Msg, ShouldEqual, msgSuccessOnGetUserProviders)
				So(res.Data, ShouldNotEqual, nil)
				var data providerData
				So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
				So(len(data.Providers), ShouldEqual, 0)
			})

			Convey("with >0 providers", func() {
				firstHourlyRate := 40
				firstId, err := ts.insertProvider(user.ID, firstHourlyRate)
				So(err, ShouldEqual, nil)
				secondhourlyRate := 45
				secondID, err := ts.insertProvider(user.ID, secondhourlyRate)
				So(err, ShouldEqual, nil)

				Convey("you should get a non-empty list", func() {
					req, err := http.NewRequest(echo.GET, "/auth/providers", nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, msgSuccessOnGetUserProviders)
					So(res.Data, ShouldNotEqual, nil)
					var data providerData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(len(data.Providers), ShouldEqual, 2)
					So(data.Providers[0].ID, ShouldBeIn, firstId, secondID)
					So(data.Providers[0].ID, ShouldBeIn, firstId, secondID)

					if data.Providers[0].ID == null.NewString(firstId, true) {
						So(data.Providers[0].HourlyRate, ShouldEqual, firstHourlyRate)
						So(data.Providers[1].HourlyRate, ShouldEqual, secondhourlyRate)
						So(data.Providers[1].ID, ShouldEqual, secondID)
					} else {
						So(data.Providers[1].HourlyRate, ShouldEqual, firstHourlyRate)
						So(data.Providers[0].HourlyRate, ShouldEqual, secondhourlyRate)
						So(data.Providers[0].ID, ShouldEqual, secondID)
					}
				})
			})
		})
	})
}

func TestServer_GetProviderByID(t *testing.T) {
	Convey("when you GET /auth/provider/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/provider/"

		Convey("for an invalid user", func() {
			token := generateToken("invalidUserId", 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url + "123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})
		})

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("with an non-existent provider id, you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url + "123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("with a valid provider id", func() {
				facID, err := ts.insertProvider(user.ID, 45)
				So(err, ShouldEqual, nil)

				Convey("and the matching user, you should get the provider", func() {
					req, err := http.NewRequest(echo.GET, url + facID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, body := doRequest(req)

					So(body.Path("error").Data(), ShouldEqual, nil)
					So(code, ShouldEqual, 200)
				})
			})
		})
	})
}

func TestServer_DeleteProviderByID(t *testing.T) {
	Convey("when you DELETE /auth/provider/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/provider/"

		Convey("for an invalid user", func() {
			token := generateToken("invalidUserId", 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("you should get an error", func() {
				req, err := http.NewRequest(echo.DELETE, url + "123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})
		})

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("but a non-existent id, you should get an error", func() {
				req, err := http.NewRequest(echo.DELETE, "/auth/provider/12345", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("and a valid facillity id", func() {
				facID, err := ts.insertProvider(user.ID, 45)
				So(err, ShouldEqual, nil)

				Convey("the provider should be deleted", func() {
					req, err := http.NewRequest(echo.DELETE, "/auth/provider/"+facID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)

					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully deleted that provider`)
					So(res.Data, ShouldNotEqual, nil)
					var data providerData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Provider, ShouldNotEqual, nil)
					So(data.Provider.ID, ShouldEqual, facID)

					b, err := ts.isUserProviderDeleted(user.ID, facID)
					So(err, ShouldEqual, nil)
					So(b, ShouldEqual, true)
				})
			})
		})
	})
}

func (s *Server) isUserProviderDeleted(userId, providerId string) (bool, error) {
	r := s.db.QueryRowx("SELECT deleted FROM user_providers "+
		"WHERE user_id=$1 AND provider_id=$2", userId, providerId)
	var b bool
	err := r.Scan(&b)
	return b, err
}

func TestServer_UpdateProviderByID(t *testing.T) {
	t.Skip("unimplemented")

	Convey("when you POST /auth/provider/update/:id", t, func() {
		//e := echo.New()
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/provider/update/"

		Convey("for an invalid user", func() {
			token := generateToken("invalidUserId", 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			//TODO
		})

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("but a non-existent id, you should get an error", func() {
				//TODO
			})

			Convey("and a valid id", func() {

				_, err := ts.insertProvider(user.ID, 45)
				So(err, ShouldEqual, nil)

				Convey("but the  wrong user, you should get an error", func() {
					//TODO
				})

				Convey("and the matching user, the provider should be updated", func() {
					req, err := http.NewRequest(echo.POST, url, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error, ShouldEqual, "")
				})
			})
		})
	})
}

func TestServer_PostProvider(t *testing.T) {
	Convey("when you POST /auth/provider", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/provider"

		token := generateToken("invalidUserId", 10*time.Minute, false)
		tStr, err := token.SignedString(SigningKey)
		So(err, ShouldEqual, nil)

		err = ts.storeToken(tStr, 10*time.Minute)
		So(err, ShouldEqual, nil)

		Convey("with an invalid provider, you should get an error", func() {
			const invalidProvider = `{"provider": -42}`
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(invalidProvider))
			So(err, ShouldEqual, nil)

			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			So(body.Path("error").Data(), ShouldEqual, `Cannot create a provider with a blank name.`)
			So(code, ShouldEqual, 200)
		})

		Convey("with a valid provider", func() {
			validProvider := &Models.Provider{

				HourlyRate: null.NewInt(40, true),
			}
			validStr, err := json.Marshal(validProvider)
			So(err, ShouldEqual, nil)

			Convey("for an invalid user", func() {
				req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("for a valid user", func() {
				user := new(Models.User)
				err := ts.insertUser("test@test.com", nil, user)
				So(err, ShouldEqual, nil)

				token := generateToken(user.ID, 10*time.Minute, false)
				tStr, err := token.SignedString(SigningKey)
				So(err, ShouldEqual, nil)

				err = ts.storeToken(tStr, 10*time.Minute)
				So(err, ShouldEqual, nil)

				Convey("with a valid provider, it should be inserted", func() {
					req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)

					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully added your provider`)
					So(res.Data, ShouldNotEqual, nil)
					var data providerData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Provider, ShouldNotEqual, nil)
					// the db sets these fields
					validProvider.ID = data.Provider.ID
					validProvider.UpdatedAt = data.Provider.UpdatedAt
					So(data.Provider, ShouldResemble, *validProvider)
				})
			})
		})
	})
}

func TestServer_PostProvidersFromCsv(t *testing.T) {
	t.Skip("unimplemented")
	Convey("when you POST /auth/providers", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/providers"

		token := generateToken("invalidUserId", 10*time.Minute, false)
		tStr, err := token.SignedString(SigningKey)
		So(err, ShouldEqual, nil)

		err = ts.storeToken(tStr, 10*time.Minute)
		So(err, ShouldEqual, nil)

		Convey("for an invalid user", func() {
			req, err := http.NewRequest(echo.POST, url, nil)
			So(err, ShouldEqual, nil)

			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
			So(code, ShouldEqual, 200)
		})

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("with invalid providers, you should get an error", func() {
				//TODO
			})

			Convey("with valid providers, they should be inserted", func() {
				//TODO
			})
		})
	})
}
