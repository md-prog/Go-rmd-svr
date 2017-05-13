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


func TestServer_GetPersonByID(t *testing.T) {
	Convey("when you GET /auth/person/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/person/"

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

			Convey("with an non-existent person id, you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url + "123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("with a valid person id", func() {
				facID, err := ts.insertPersonIntoDb(user.ID, &Models.Person{
					FirstName:	"test-fname",
					LastName:	"test-lname",
					Email:		"test-email",
				})
				So(err, ShouldEqual, nil)

				Convey("and the matching user, you should get the person", func() {
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


func TestServer_UpdatePersonByID(t *testing.T) {
	t.Skip("unimplemented")

	Convey("when you POST /auth/person/update/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/person/update/"

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

				_, err := ts.insertPersonIntoDb(user.ID, &Models.Person{
					FirstName:	"test-fname",
					LastName:	"test-lname",
					Email:		"test-email",
				})
				So(err, ShouldEqual, nil)

				Convey("but the  wrong user, you should get an error", func() {
					//TODO
				})

				Convey("and the matching user, the person should be updated", func() {
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

func TestServer_GetUserPersons(t *testing.T) {
	Convey("when you GET /auth/persons", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/persons"

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("with 0 persons, you should get an empty list", func() {
				req, err := http.NewRequest(echo.GET, url, nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, res := doRequestRes(req)

				So(code, ShouldEqual, 200)
				So(res.Error, ShouldEqual, "")
				So(res.Msg, ShouldEqual, msgSuccessOnGetUserPersons)
				So(res.Data, ShouldNotEqual, nil)
				var data personData
				So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
				So(len(data.People), ShouldEqual, 0)
			})

			Convey("with >0 persons", func() {
				firstPerson := &Models.Person{
					FirstName:	"test-fname-A",
					LastName:	"test-lname-A",
					Email:		"testA@test.com",
				}
				firstId, err := ts.insertPersonIntoDb(user.ID, firstPerson)
				So(err, ShouldEqual, nil)

				secondPerson := &Models.Person{
					FirstName:	"test-fname-B",
					LastName:	"test-lname-B",
					Email:		"testB@test.com",
				}
				secondId, err := ts.insertPersonIntoDb(user.ID, secondPerson)
				So(err, ShouldEqual, nil)

				Convey("you should get a non-empty list", func() {
					req, err := http.NewRequest(echo.GET, url, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, msgSuccessOnGetUserPersons)
					So(res.Data, ShouldNotEqual, nil)
					var data personData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(len(data.People), ShouldEqual, 2)
					So(data.People[0].ID, ShouldBeIn, firstId, secondId)
					So(data.People[0].ID, ShouldBeIn, firstId, secondId)

					if data.People[0].ID == firstId {
						So(data.People[0].FirstName, ShouldEqual, firstPerson.FirstName)
						So(data.People[0].LastName, ShouldEqual, firstPerson.LastName)
						So(data.People[0].Email, ShouldEqual, firstPerson.Email)
						So(data.People[1].FirstName, ShouldEqual, secondPerson.FirstName)
						So(data.People[1].LastName, ShouldEqual, secondPerson.LastName)
						So(data.People[1].Email, ShouldEqual, secondPerson.Email)
						So(data.People[1].ID, ShouldEqual, secondId)
					} else {
						So(data.People[1].FirstName, ShouldEqual, firstPerson.FirstName)
						So(data.People[1].LastName, ShouldEqual, firstPerson.LastName)
						So(data.People[1].Email, ShouldEqual, firstPerson.Email)
						So(data.People[0].FirstName, ShouldEqual, secondPerson.FirstName)
						So(data.People[0].LastName, ShouldEqual, secondPerson.LastName)
						So(data.People[0].Email, ShouldEqual, secondPerson.Email)
						So(data.People[0].ID, ShouldEqual, secondId)
					}
				})
			})
		})
	})
}


func TestServer_DeletePersonByID(t *testing.T) {
	Convey("when you DELETE /auth/person/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/person/"

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
				req, err := http.NewRequest(echo.DELETE, "/auth/person/12345", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
				So(code, ShouldEqual, 200)
			})

			Convey("and a valid person id", func() {
				facID, err := ts.insertPersonIntoDb(user.ID, &Models.Person{
					FirstName:	"test-fname",
					LastName:	"test-lname",
					Email:		"test@test.com",
				})
				So(err, ShouldEqual, nil)

				Convey("the person should be deleted", func() {
					req, err := http.NewRequest(echo.DELETE, "/auth/person/"+facID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully deleted that person`)
					So(res.Data, ShouldNotEqual, nil)
					var data personData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Person, ShouldNotEqual, nil)
					So(data.Person.ID, ShouldEqual, facID)

					b, err := ts.isUserPersonDeleted(user.ID, facID)
					So(err, ShouldEqual, nil)
					So(b, ShouldEqual, true)
				})
			})
		})
	})
}

func (s *Server) isUserPersonDeleted(userId, personId string) (bool, error) {
	r := s.db.QueryRowx("SELECT deleted FROM user_contacts "+
		"WHERE user_id=$1 AND person_id=$2", userId, personId)
	var b bool
	err := r.Scan(&b)
	return b, err
}


func TestServer_PostPerson(t *testing.T) {
	Convey("when you POST /auth/person", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/person"

		token := generateToken("invalidUserId", 10*time.Minute, false)
		tStr, err := token.SignedString(SigningKey)
		So(err, ShouldEqual, nil)

		err = ts.storeToken(tStr, 10*time.Minute)
		So(err, ShouldEqual, nil)

		Convey("with an invalid person, you should get an error", func() {
			const invalidPerson = `{"person": -42}`
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(invalidPerson))
			So(err, ShouldEqual, nil)

			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			So(body.Path("error").Data(), ShouldEqual, `Cannot create a person with a blank name.`)
			So(code, ShouldEqual, 200)
		})

		Convey("with a valid person", func() {
			validPerson := &Models.Person{
				FirstName: "test-name",
				LastName: "test-name",

				Contacts: types.JSONText("{}"),
			}
			validStr, err := json.Marshal(validPerson)
			So(err, ShouldEqual, nil)

			Convey("for an invalid user", func() {
				req, err := http.NewRequest(echo.POST, "/auth/person", bytes.NewReader(validStr))
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

				Convey("with a valid person, it should be inserted", func() {
					req, err := http.NewRequest(echo.POST, "/auth/person", bytes.NewReader(validStr))
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)

					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully added your person`)
					So(res.Data, ShouldNotEqual, nil)
					var data personData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Person, ShouldNotEqual, nil)
					// the db sets these fields
					validPerson.ID = data.Person.ID
					validPerson.UpdatedAt = data.Person.UpdatedAt
					So(data.Person, ShouldResemble, *validPerson)
				})
			})
		})
	})
}

func TestServer_PostPersonsFromCsv(t *testing.T) {
	t.Skip("unimplemented")
	Convey("when you POST /auth/persons", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/persons"

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

			Convey("with invalid persons, you should get an error", func() {
				//TODO
			})

			Convey("with valid persons, they should be inserted", func() {
				//TODO
			})
		})
	})
}
