package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"

	"gitlab.com/chrislewispac/rmd-server/models"
)

func TestServer_GetUserFacilities(t *testing.T) {
	Convey("when you GET /auth/facilities", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/facilities"

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("with 0 facilities, you should get an empty list", func() {
				req, err := http.NewRequest(echo.GET, url, nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, res := doRequestRes(req)

				So(code, ShouldEqual, 200)
				So(res.Error.String, ShouldEqual, "")
				So(res.Msg, ShouldEqual, `Retrieved Facilities`)
				So(res.Data, ShouldNotEqual, nil)
				var data facilityData
				So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
				So(len(data.Facilities), ShouldEqual, 0)
			})

			Convey("with >0 facilities", func() {
				firstName := "test-name-A"
				firstId, err := ts.insertFacility(user.ID, firstName)
				So(err, ShouldEqual, nil)
				secondName := "test-name-B"
				secondID, err := ts.insertFacility(user.ID, secondName)
				So(err, ShouldEqual, nil)

				Convey("you should get a non-empty list", func() {
					req, err := http.NewRequest(echo.GET, url, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Retrieved Facilities`)
					So(res.Data, ShouldNotEqual, nil)
					var data facilityData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(len(data.Facilities), ShouldEqual, 2)
					So(data.Facilities[0].ID, ShouldBeIn, firstId, secondID)

					if data.Facilities[0].ID == firstId {
						So(data.Facilities[0].Name, ShouldEqual, firstName)
						So(data.Facilities[1].Name, ShouldEqual, secondName)
						So(data.Facilities[1].ID, ShouldEqual, secondID)
					} else {
						So(data.Facilities[1].ID, ShouldEqual, firstId)
						So(data.Facilities[1].Name, ShouldEqual, firstName)
						So(data.Facilities[0].Name, ShouldEqual, secondName)
						So(data.Facilities[0].ID, ShouldEqual, secondID)
					}
				})
			})
		})
	})
}

func TestServer_GetFacilityByID(t *testing.T) {
	Convey("when you GET /auth/facility/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/facility/"

		Convey("for an invalid user", func() {
			token := generateToken("invalidUserId", 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url+"123435", nil)
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

			Convey("with an non-existent facility id, you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url+"123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("with a valid facility id", func() {
				facID, err := ts.insertFacility(user.ID, "test-name")
				So(err, ShouldEqual, nil)

				Convey("and the matching user, you should get the facility", func() {
					req, err := http.NewRequest(echo.GET, url+facID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(res.Error.String, ShouldEqual, "")
					So(code, ShouldEqual, 200)
					So(res.Msg, ShouldEqual, `Successfully retrieved your facility`)
					So(res.Data, ShouldNotEqual, nil)
					var data facilityData
					So(res.Data.Unmarshal(&data), ShouldEqual, nil)
					So(len(data.Facilities), ShouldEqual, 0)
					So(data.Facility, ShouldNotEqual, nil)
					So(data.Facility.ID, ShouldEqual, facID)
				})
			})
		})
	})
}

func TestServer_DeleteFacilityByID(t *testing.T) {
	Convey("when you DELETE /auth/facility/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/facility/"

		Convey("for an invalid user", func() {
			token := generateToken("invalidUserId", 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("you should get an error", func() {
				req, err := http.NewRequest(echo.DELETE, url+"123435", nil)
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
				req, err := http.NewRequest(echo.DELETE, url+"12345", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("and a valid facillity id", func() {
				facID, err := ts.insertFacility(user.ID, "test-name")
				So(err, ShouldEqual, nil)

				Convey("the facility should be deleted", func() {
					req, err := http.NewRequest(echo.DELETE, url+facID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully deleted that facility`)
					So(res.Data, ShouldNotEqual, nil)
					var data facilityData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Facility, ShouldNotEqual, nil)
					So(data.Facility.ID, ShouldEqual, facID)

					b, err := ts.isUserFacilityDeleted(user.ID, facID)
					So(err, ShouldEqual, nil)
					So(b, ShouldEqual, true)
				})
			})
		})
	})
}

func (s *Server) isUserFacilityDeleted(userID, facilityId string) (bool, error) {
	r := s.db.QueryRowx("SELECT deleted FROM user_facilities "+
		"WHERE user_id=$1 AND facility_id=$2", userID, facilityId)
	var b bool
	err := r.Scan(&b)
	return b, err
}

func TestServer_UpdateFacilityByID(t *testing.T) {
	t.Skip("unimplemented")

	Convey("when you POST /auth/facility/update/:id", t, func() {
		ts, _, close := startTestServer(t)
		defer close()

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

				_, err := ts.insertFacility(user.ID, "test-name")
				So(err, ShouldEqual, nil)

				Convey("but the  wrong user, you should get an error", func() {
					//TODO
				})

				Convey("and the matching user, the facility should be updated", func() {
					//TODO confirm update
				})
			})
		})
	})
}

func TestServer_PostFacility(t *testing.T) {
	Convey("when you POST /auth/facility", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/facility"

		token := generateToken("invalidUserId", 10*time.Minute, false)
		tStr, err := token.SignedString(SigningKey)
		So(err, ShouldEqual, nil)

		err = ts.storeToken(tStr, 10*time.Minute)
		So(err, ShouldEqual, nil)

		Convey("with an invalid facility, you should get an error", func() {
			const invalidFacility = `{"facility": -42}`
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(invalidFacility))
			So(err, ShouldEqual, nil)

			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			So(code, ShouldEqual, 200)
			So(body.Path("error").Data(), ShouldEqual, `Cannot create a facility with a blank name.`)
		})

		Convey("with a valid facility", func() {
			validFacility := &Models.Facility{
				Name: "test-name",

				Contacts: types.JSONText("{}"),
			}
			validStr, err := json.Marshal(validFacility)
			So(err, ShouldEqual, nil)

			Convey("for an invalid user", func() {
				req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(code, ShouldEqual, 200)
				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
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

				Convey("with a valid facility, it should be inserted", func() {
					req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully added your facility`)
					So(res.Data, ShouldNotEqual, nil)
					var data facilityData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Facility, ShouldNotEqual, nil)
					// the db sets these fields
					validFacility.ID = data.Facility.ID
					validFacility.UpdatedAt = data.Facility.UpdatedAt
					validFacility.AdmissionOrders = data.Facility.AdmissionOrders
					So(data.Facility, ShouldResemble, *validFacility)
				})
			})
		})
	})
}

func TestServer_PostFacilitiesFromCsv(t *testing.T) {
	t.Skip("unimplemented")
	Convey("when you POST /auth/facilities", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/facilities"

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

			Convey("with invalid facilities, you should get an error", func() {
				//TODO
			})

			Convey("with valid facilities, they should be inserted", func() {
				//TODO
			})
		})
	})
}
