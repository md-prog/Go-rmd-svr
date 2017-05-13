package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/guregu/null.v3"

	"github.com/jmoiron/sqlx/types"
	"gitlab.com/chrislewispac/rmd-server/models"
)

func TestServer_GetUserContracts(t *testing.T) {
	Convey("when you GET /auth/contracts", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/contracts"

		Convey("for a valid user", func() {
			user := new(Models.User)
			err := ts.insertUser("test@test.com", nil, user)
			So(err, ShouldEqual, nil)

			token := generateToken(user.ID, 10*time.Minute, false)
			tStr, err := token.SignedString(SigningKey)
			So(err, ShouldEqual, nil)

			err = ts.storeToken(tStr, 10*time.Minute)
			So(err, ShouldEqual, nil)

			Convey("with 0 contracts, you should get an empty list", func() {
				req, err := http.NewRequest(echo.GET, url, nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, res := doRequestRes(req)

				So(code, ShouldEqual, 200)
				So(res.Error.String, ShouldEqual, "")
				So(res.Msg, ShouldEqual, `Successfully retrieved your contracts`)
				So(res.Data, ShouldNotEqual, nil)
				var data contractData
				So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
				So(len(data.Contracts), ShouldEqual, 0)
			})

			Convey("with >0 contracts", func() {
				now := time.Now().In(time.UTC).Round(time.Minute)
				firstName := "test-name-A"
				firstId, firstUC, err := ts.insertContract(user.ID, &Models.Contract{
					Name:        firstName,
					StartDate:   null.TimeFrom(now),
					RenewalDate: null.TimeFrom(now.AddDate(1, 0, 0)),
				})
				So(err, ShouldEqual, nil)

				facilityA, err := ts.insertFacility(user.ID, "facility-a")
				So(err, ShouldEqual, nil)
				facilityB, err := ts.insertFacility(user.ID, "facility-b")
				So(err, ShouldEqual, nil)

				then := time.Now().AddDate(0, 1, 0).In(time.UTC).Round(time.Minute)
				secondName := "test-name-B"
				secondID, secondUC, err := ts.insertContract(user.ID, &Models.Contract{
					Name:        secondName,
					StartDate:   null.TimeFrom(then),
					RenewalDate: null.TimeFrom(then.AddDate(1, 0, 0)),
					Facilities:  types.NullJSONText{Valid: true, JSONText: []byte(`["` + facilityA + `","` + facilityB + `"]`)},
				})
				So(err, ShouldEqual, nil)

				Convey("you should get a non-empty list", func() {
					req, err := http.NewRequest(echo.GET, url, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully retrieved your contracts`)
					So(res.Data, ShouldNotEqual, nil)
					var data contractData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(len(data.Contracts), ShouldEqual, 2)
					So(data.Contracts[0].ID, ShouldBeIn, firstId, secondID)
					So(data.Contracts[0].ID, ShouldBeIn, firstId, secondID)

					var facilities []string
					if data.Contracts[0].ID == firstId {
						So(data.Contracts[0].UsersContractsID, ShouldEqual, firstUC)
						So(data.Contracts[0].Name, ShouldEqual, firstName)
						So(data.Contracts[1].Name, ShouldEqual, secondName)
						So(data.Contracts[1].ID, ShouldEqual, secondID)
						So(data.Contracts[1].UsersContractsID, ShouldEqual, secondUC)
						fStr := data.Contracts[1].Facilities.String()
						So(len(fStr), ShouldNotEqual, 0)
						So(fStr, ShouldNotEqual, "{}")

						So(data.Contracts[1].Facilities.Unmarshal(&facilities), ShouldEqual, nil)
					} else {
						So(data.Contracts[1].ID, ShouldEqual, firstId)
						So(data.Contracts[1].UsersContractsID, ShouldEqual, firstUC)
						So(data.Contracts[1].Name, ShouldEqual, firstName)
						So(data.Contracts[0].Name, ShouldEqual, secondName)
						So(data.Contracts[0].ID, ShouldEqual, secondID)
						So(data.Contracts[0].UsersContractsID, ShouldEqual, secondUC)
						fStr := data.Contracts[0].Facilities.String()
						So(len(fStr), ShouldNotEqual, 0)
						So(fStr, ShouldNotEqual, "{}")

						So(data.Contracts[0].Facilities.Unmarshal(&facilities), ShouldEqual, nil)
					}
					So(len(facilities), ShouldEqual, 2)
					if facilities[0] == facilityA {
						So(facilities[1], ShouldEqual, facilityB)
					} else {
						So(facilities[1], ShouldEqual, facilityA)
						So(facilities[0], ShouldEqual, facilityB)
					}
				})
			})
		})
	})
}

func TestServer_GetContractByID(t *testing.T) {
	Convey("when you GET /auth/contract/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/contract/"

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

				So(body.Path("error").Data(), ShouldEqual, `Something Went Wrong!`)
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

			Convey("with an non-existent contract id, you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url+"123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, `Something Went Wrong!`)
				So(code, ShouldEqual, 200)
			})

			Convey("with a valid contract id", func() {
				facilityA, err := ts.insertFacility(user.ID, "facility-a")
				So(err, ShouldEqual, nil)
				facilityB, err := ts.insertFacility(user.ID, "facility-b")
				So(err, ShouldEqual, nil)

				now := time.Now().In(time.UTC).Round(time.Minute)
				ctID, ucID, err := ts.insertContract(user.ID, &Models.Contract{
					Name:        "test-name",
					StartDate:   null.TimeFrom(now),
					RenewalDate: null.TimeFrom(now.AddDate(1, 0, 0)),
					Facilities:  types.NullJSONText{Valid: true, JSONText: []byte(`["` + facilityA + `","` + facilityB + `"]`)},
				})
				So(err, ShouldEqual, nil)

				Convey("and the matching user, you should get the contract", func() {
					req, err := http.NewRequest(echo.GET, url+ctID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderAcceptEncoding, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(res.Error.String, ShouldEqual, "")
					So(code, ShouldEqual, 200)
					So(res.Msg, ShouldEqual, `Successfully retrieved your contract`)
					So(res.Data, ShouldNotEqual, nil)
					var data contractData
					So(res.Data.Unmarshal(&data), ShouldEqual, nil)
					So(data.Contract, ShouldNotEqual, nil)
					So(data.Contract.ID, ShouldEqual, ctID)
					So(data.Contract.UsersContractsID, ShouldEqual, ucID)

					fStr := data.Contract.Facilities.String()
					So(len(fStr), ShouldNotEqual, 0)
					So(fStr, ShouldNotEqual, "{}")

					var facilities []string
					So(data.Contract.Facilities.Unmarshal(&facilities), ShouldEqual, nil)

					So(len(facilities), ShouldEqual, 2)
					if facilities[0] == facilityA {
						So(facilities[1], ShouldEqual, facilityB)
					} else {
						So(facilities[1], ShouldEqual, facilityA)
						So(facilities[0], ShouldEqual, facilityB)
					}
				})
			})
		})
	})
}

//TODO dump
func (s *Server) insertContractFacility(userID, facilityName, contractId string) (string, error) {
	const insertFacility = `INSERT INTO facilities(name) VALUES ($1) RETURNING id;`

	const insertUserFacilityRelation = `INSERT INTO user_facilities(user_id, facility_id) VALUES ($1, $2) RETURNING id;`

	const insertContractFacilities = `INSERT INTO contract_facilities(user_facilities_id, contract_id) VALUES ($1, $2);`

	var fID string
	return fID, s.tx(func(tx *sqlx.Tx) error {
		err := tx.Get(&fID, insertFacility, facilityName)
		if err != nil {
			return err
		}

		var ufID string
		err = tx.Get(&ufID, insertUserFacilityRelation, userID, fID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(insertContractFacilities, ufID, contractId)
		return err
	})
}

func TestServer_DeleteContractByID(t *testing.T) {
	Convey("when you DELETE /auth/contract/:id", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/contract/"

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

				So(body.Path("error").Data(), ShouldEqual, `Something Went Wrong!`)
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

				So(body.Path("error").Data(), ShouldEqual, `Something Went Wrong!`)
				So(code, ShouldEqual, 200)
			})

			Convey("and a valid contract id", func() {
				now := time.Now().In(time.UTC).Round(time.Minute)
				ctID, _, err := ts.insertContract(user.ID, &Models.Contract{
					Name:        "test-name",
					StartDate:   null.TimeFrom(now),
					RenewalDate: null.TimeFrom(now.AddDate(1, 0, 0)),
				})
				So(err, ShouldEqual, nil)

				Convey("the contract should be deleted", func() {
					req, err := http.NewRequest(echo.DELETE, url+ctID, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully deleted that contract`)
					So(res.Data, ShouldNotEqual, nil)
					var data contractData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Contract, ShouldNotEqual, nil)
					So(data.Contract.ID, ShouldEqual, ctID)

					del, err := ts.isUserContractDeleted(user.ID, ctID)
					So(err, ShouldEqual, nil)
					So(del, ShouldEqual, true)
				})
			})
		})
	})
}

func (s *Server) isUserContractDeleted(userID, contractId string) (bool, error) {
	r := s.db.QueryRowx("SELECT deleted FROM user_contracts "+
		"WHERE user_id=$1 AND contract_id=$2", userID, contractId)
	var b bool
	err := r.Scan(&b)
	return b, err
}

func TestServer_UpdateContractByID(t *testing.T) {
	t.Skip("unimplemented")

	Convey("when you POST /auth/contract/update/:id", t, func() {
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
				//TODO insert contract

				Convey("but the  wrong user, you should get an error", func() {
					//TODO
				})

				Convey("and the matching user, the contract should be updated", func() {
					//TODO confirm update
				})
			})
		})
	})
}

func TestServer_PostContract(t *testing.T) {
	Convey("when you POST /auth/contract", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/contract"

		token := generateToken("invalidUserId", 10*time.Minute, false)
		tStr, err := token.SignedString(SigningKey)
		So(err, ShouldEqual, nil)

		err = ts.storeToken(tStr, 10*time.Minute)
		So(err, ShouldEqual, nil)

		Convey("with an invalid contract, you should get an error", func() {
			const invalidContract = `{"contract": -42}`
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(invalidContract))
			So(err, ShouldEqual, nil)

			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			code, body := doRequest(req)

			So(body.Path("error").Data(), ShouldEqual, `There was an error parsing the name field`)
			So(code, ShouldEqual, 200)
		})

		Convey("with a valid contract", func() {
			now := time.Now().In(time.UTC).Round(time.Minute)
			validContract := &Models.Contract{
				Name:        "test-name",
				StartDate:   null.TimeFrom(now),
				RenewalDate: null.TimeFrom(now.AddDate(1, 0, 0)),
			}
			validStr, err := json.Marshal(validContract)
			So(err, ShouldEqual, nil)

			Convey("for an invalid user", func() {
				req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(body.Path("error").Data(), ShouldEqual, `Something Went Wrong!`)
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

				Convey("with a valid contract, it should be inserted", func() {
					req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, `Successfully added your contract`)
					So(res.Data, ShouldNotEqual, nil)
					var data contractData
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.Contract, ShouldNotEqual, nil)
					// the db sets these fields
					validContract.ID = data.Contract.ID
					validContract.UsersContractsID = data.Contract.UsersContractsID
					validContract.UpdatedAt = data.Contract.UpdatedAt
					// ignore any garbage JSONText, e.g. {} or [null]
					So(data.Contract.Facilities.Valid, ShouldEqual, false)
					data.Contract.Facilities = types.NullJSONText{}
					So(data.Contract.Contacts.Valid, ShouldEqual, false)
					data.Contract.Contacts = types.NullJSONText{}
					So(data.Contract.StatusChanges.Valid, ShouldEqual, false)
					data.Contract.StatusChanges = types.NullJSONText{}

					So(*data.Contract, ShouldResemble, *validContract)
				})
			})
		})
	})
}

func TestServer_PostContractsFromCsv(t *testing.T) {
	t.Skip("unimplemented")
	Convey("when you POST /auth/contracts", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/contracts"

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

			So(body.Path("error").Data(), ShouldEqual, `Something Went Wrong!`)
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

			Convey("with invalid contracts, you should get an error", func() {
				//TODO
			})

			Convey("with valid contracts, they should be inserted", func() {
				//TODO
			})
		})
	})
}
