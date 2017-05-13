package server

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"errors"

	log "github.com/Sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"

	"github.com/Jeffail/gabs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	Models "gitlab.com/chrislewispac/rmd-server/models"

	"github.com/jmoiron/sqlx"
	"strings"
)

// GetPersonByID exported
// Get profile information of a contact of the current user
func (s *Server) GetPersonByID(c echo.Context) error {

	id := userID(c)
	contact_id := c.Param("id")

	person, err := s.getPersonFromDb(id, contact_id)

	if err != nil {
		createPersonErrorResponse(person, nil, errMsgSomethingWrong)
	}

	anon := struct {
		Person Models.Person `json:"person"`
	}{
		*person,
	}

	data, _ := json.Marshal(anon)

	res := Models.NewResponse()
	res.Msg = msgSuccessOnGetPersonByID
	res.Data = data

	return c.JSON(http.StatusOK, res)
}

// UpdatePersonByID exported
func (s *Server) UpdatePersonByID(c echo.Context) error {
	//x := c.Get("user").(*jwt.Token)
	//claims := x.Claims.(jwt.MapClaims)
	//id := claims["id"].(string)

	// parse the request body and check the validity
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
	}

	p0, err := convertJsonToPerson(b)

	if(err != nil) {
		handleErr(err)
		return c.JSON(http.StatusOK, createPersonErrorResponse(p0, nil, err.Error()))
	}

	res := Models.NewResponse()
	res.Msg = "Successfully updated contact #1"
	// TODO: implementation
	return c.JSON(http.StatusOK, res)
}


// GetUserPersons exported
// Get the full list of contacts of the current user
func (s *Server) GetUserPersons(c echo.Context) error {
	var person Models.Person

	id := userID(c)

	people, _ := s.getUserPersonsFromDb(id)

	res := createPersonSuccessResponse(&person, people, msgSuccessOnGetUserPersons)

	return c.JSON(http.StatusOK, res)
}

// PostPersonsFromCsv exported
// Mass upload contacts from CSV
func (s *Server) PostPersonsFromCsv(c echo.Context) error {
	people := []Models.Person{}
	var person Models.Person

	id := userID(c)

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
		res := createPersonErrorResponse(&person, people, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}
	jsonParsed, err := gabs.ParseJSON(b)
	if err != nil {
		handleErr(err)
		res := createPersonErrorResponse(&person, people, errMsgBadRequest)
		return c.JSON(http.StatusOK, res)
	}

	children, _ := jsonParsed.Children()
	tx, err := s.db.Beginx()

	for _, child := range children {
		ct := &Models.Person{
			FirstName:   child.S("first_name").Data().(string),
			LastName:    child.S("last_name").Data().(string),
			MiddleName:  Models.HandleNullString(null.StringFrom(child.S("middle_name").Data().(string))),
			Title:       Models.HandleNullString(null.StringFrom(child.S("title").Data().(string))),
			Email:       child.S("email").Data().(string),
			CellPhone:   Models.HandleNullString(null.StringFrom(child.S("cell_phone").Data().(string))),
			HomePhone:   Models.HandleNullString(null.StringFrom(child.S("home_phone").Data().(string))),
			OfficePhone: Models.HandleNullString(null.StringFrom(child.S("office_phone").Data().(string))),
			Fax:         Models.HandleNullString(null.StringFrom(child.S("fax").Data().(string))),
			Address:     Models.HandleNullString(null.StringFrom(child.S("address").Data().(string))),
			City:        Models.HandleNullString(null.StringFrom(child.S("city").Data().(string))),
			Zip:         Models.HandleNullString(null.StringFrom(child.S("zip").Data().(string))),
			StateID:     Models.HandleNullString(null.StringFrom(child.S("state_id").Data().(string))),
			RecruiterID: Models.HandleNullString(null.StringFrom(child.S("recruiter_id").Data().(string))),
			ProviderID:  Models.HandleNullString(null.StringFrom(child.S("provider_id").Data().(string))),
			Notes:       Models.HandleNullString(null.StringFrom(child.S("notes").Data().(string))),
			Lat:	     parseNullableFloat(child.S("lat").Data().(string)),
			Lon:	     parseNullableFloat(child.S("lon").Data().(string)),
		}

		_, err := s.insertPersonIntoDb(id, ct)
		handleErr(err)
	}
	err = tx.Commit()
	handleErr(err)

	rows, e := s.db.Queryx(`
	SELECT * FROM (
		SELECT people_id AS id
		FROM user_contacts uc
		WHERE uc.user_id=$1
		) u
	INNER JOIN people c
	USING (id)`, id)
	if e != nil {
		log.Println(e)
	}

	for rows.Next() {
		var ct Models.Person
		err = rows.StructScan(&ct)
		handleErr(err)
		people = append(people, ct)
	}

	err = rows.Err()
	if err != nil {
		handleErr(err)
	}

	res := createPersonSuccessResponse(&person, people, msgSuccessOnPostPersonsFromCsv)

	return c.JSON(http.StatusOK, res)
}

// PostPerson exported
// create a contact
func (s *Server) PostPerson(c echo.Context) error {
	var person Models.Person
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	b, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		handleErr(err)
		res := createPersonErrorResponse(&person, nil, errMsgBadRequest)
		return c.JSON(http.StatusOK, res)
	}

	p0, err := convertJsonToPerson(b)

	if(err != nil) {
		handleErr(err)
		return c.JSON(http.StatusOK, createPersonErrorResponse(p0, nil, err.Error()))
	}

	pID, err := s.insertPersonIntoDb(id, p0)

	if err != nil {
		handleErr(err)
	}

	p1, err := s.getPersonFromDb(id, pID)

	if err != nil {
		handleErr(err)
	}

	people := []Models.Person{}

	res := createPersonSuccessResponse(p1, people, msgSuccessOnPostPerson)

	return c.JSON(http.StatusOK, res)
}

// DeletePersonByID exported
// Mark a contact as deleted
func (s *Server) DeletePersonByID(c echo.Context) error {
	const deletePersonByIdQuery = `
UPDATE user_contacts
	SET deleted=true
	WHERE people_id=$1 AND user_id=$2
	RETURNING people_id as id`

	var person Models.Person
	person_id := c.Param("id")
	user_id := userID(c)

	err := s.db.QueryRowx(deletePersonByIdQuery, person_id, user_id).StructScan(&person)

	if err != nil {
		handleErr(err)
		res := createPersonErrorResponse(&person, nil, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createPersonSuccessResponse(&person, nil, msgSuccessOnDeletePersonByID)
	return c.JSON(http.StatusOK, res)
}

type personData struct {
	Person Models.Person   `json:"person"`
	People []Models.Person `json:"people"`
}

const(
	msgSuccessOnGetUserPersons = "Successfully retrieved your contacts"
	msgSuccessOnGetPersonByID = "Successfully retrieved contact #1"
	msgSuccessOnPostPerson = "Successfully created contact #1"
	msgSuccessOnPostPersonsFromCsv = "Successfully added people from CSV"
	msgSuccessOnDeletePersonByID = "Successfully deleted that person"
	msgErrorMissingRequiredFieldOnPostPerson = "Missing required field"
)

func createPersonErrorResponse(person *Models.Person, people []Models.Person, errMsg string) *Models.Res {
	anon := personData{
		*person,
		people,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = errMsgExists
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createPersonSuccessResponse(person *Models.Person, people []Models.Person, successMsg string) *Models.Res {
	anon := personData{
		*person,
		people,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}


// this function will find a user's contacts with his/her people_id
func (s *Server) getPersonFromDb(userId string, personId string) (*Models.Person, error) {
	var person Models.Person

	err := s.db.QueryRowx(`
	SELECT * FROM (
			SELECT people_id AS id
			FROM user_contacts uc
			WHERE uc.user_id=$1
		) u
	INNER JOIN people c
	USING (id)
	WHERE u.id=$2`, userId, personId).StructScan(&person)
	return &person, err
}

/**
 this function will insert a person object into people DB table

 @return	string	newly inserted person's ID
 		error	Error occured
*/
func (s *Server) insertPersonIntoDb(userId string, person *Models.Person) (string, error) {
	const insertPersonQuery = `
	INSERT INTO people(
		first_name, last_name, email,
		title, cell_phone, home_phone, office_phone, fax,
		address, zip, city, lat, lon)
	VALUES ($1, $2, $3,
		$4, $5, $6, $7, $8,
		$9, $10, $11, $12, $13)
	RETURNING id`

	const insertUserPersonRelationQuery = `
	INSERT INTO user_contacts(
		user_id, people_id)
	VALUES ($1, $2);`

	var fID string

	return fID, s.tx(func(tx *sqlx.Tx) error {
		err := tx.Get(&fID, insertPersonQuery, person.FirstName, person.LastName, person.Email,
			person.Title, person.CellPhone, person.HomePhone, person.OfficePhone, person.Fax,
			person.Address, person.Zip, person.City, person.Lat, person.Lon)
		if err != nil {
			return err
		}

		_, err = tx.Exec(insertUserPersonRelationQuery, userId, fID)
		return err
	})
}

// this function will find a user's contacts with his/her people_id
func (s *Server) getUserPersonsFromDb(userId string) ([]Models.Person, error) {
	var people []Models.Person
	const getUserPersonsQuery = `
	SELECT * FROM (
		SELECT people_id as id, json_agg(ci) as contacts
		FROM user_contacts uc
		LEFT OUTER JOIN contact_instance ci
		ON ci.user_contacts_id=uc.id
		WHERE uc.user_id=$1
		AND uc.deleted=false
		GROUP BY uc.people_id
	) u
	LEFT OUTER JOIN people c
	USING (id)
	LEFT OUTER JOIN (
		SELECT id "provider_id"
		, id "provider.id"
		, npi "provider.npi"
		, w9 "provider.w9"
		, direct_deposit_form "provider.direct_deposit_form"
		, hourly_rate "provider.hourly_rate"
		, desired_shifts_month "provider.desired_shifts_month"
		, max_shifts_month "provider.max_shifts_month"
		, min_shifts_month "provider.min_shifts_month"
		, full_time "provider.full_time"
		, part_time "provider.part_time"
		, prn "provider.prn"
		, retired "provider.retired"
		, notes "provider.notes"
		, insurance_certificate "provider.insurance_certificate"
		, tb_expiration "provider.tb_expiration"
		, tb_file "provider.tb_file"
		, flu_expiration "provider.flu_expiration"
		, flu_file "provider.flu_file" FROM providers
	) as provider
	USING (provider_id)
	LEFT OUTER JOIN (
		SELECT id "state.id", id "state_id", name "state.name", abbreviation as "state.abbr"
		FROM states
	) st
	USING (state_id)
	LEFT OUTER JOIN (
		SELECT id "recruiter_id", id "recruiter.id"
		FROM recruiters
	) rr
	USING (recruiter_id)
	`
	rows, e := s.db.Queryx(getUserPersonsQuery, userId)
	handleErr(e)

	for rows.Next() {
		var ct Models.Person
		err := rows.StructScan(&ct)
		handleErr(err)
		people = append(people, ct)
	}

	err := rows.Err()
	handleErr(err)

	return people, err
}


// this function will be used to convert Request JSON string into Person object
func convertJsonToPerson(json []byte) (*Models.Person, error) {
	var e error
	jsonParsed := MustParseJSON(string(json))

	// check non-nullable required fields
	first_name, ok := jsonParsed.Path("first_name").Data().(string)
	if ok == false {
		e = errors.New(strings.Join([]string{msgErrorMissingRequiredFieldOnPostPerson, "first_name"}, ":"))
	}

	last_name, ok := jsonParsed.Path("last_name").Data().(string)
	if ok == false {
		e = errors.New(strings.Join([]string{msgErrorMissingRequiredFieldOnPostPerson, "last_name"}, ":"))
	}

	email, ok := jsonParsed.Path("email").Data().(string)
	if ok == false {
		e = errors.New(strings.Join([]string{msgErrorMissingRequiredFieldOnPostPerson, "email"}, ":"))
	}

	// fill other nullable fields
	middle_name, ok := jsonParsed.Path("middle_name").Data().(string)
	title, ok := jsonParsed.Path("title").Data().(string)
	cell_phone, ok := jsonParsed.Path("cell_phone").Data().(string)
	home_phone, ok := jsonParsed.Path("home_phone").Data().(string)
	office_phone, ok := jsonParsed.Path("office_phone").Data().(string)
	fax, ok := jsonParsed.Path("fax").Data().(string)
	address, ok := jsonParsed.Path("address").Data().(string)
	zip, ok := jsonParsed.Path("zip").Data().(string)
	city, ok := jsonParsed.Path("city").Data().(string)
	latStr, ok := jsonParsed.Path("lat").Data().(string)
	lonStr, ok := jsonParsed.Path("lon").Data().(string)

	p := &Models.Person{
		FirstName:	first_name,
		MiddleName:	Models.HandleNullString(null.StringFrom(middle_name)),
		LastName:	last_name,
		Email:		email,
		Title:		Models.HandleNullString(null.StringFrom(title)),
		CellPhone:	Models.HandleNullString(null.StringFrom(cell_phone)),
		HomePhone:	Models.HandleNullString(null.StringFrom(home_phone)),
		OfficePhone:	Models.HandleNullString(null.StringFrom(office_phone)),
		Fax:		Models.HandleNullString(null.StringFrom(fax)),
		Address:	Models.HandleNullString(null.StringFrom(address)),
		Zip:		Models.HandleNullString(null.StringFrom(zip)),
		City:		Models.HandleNullString(null.StringFrom(city)),
		Lat:		parseNullableFloat(latStr),
		Lon:		parseNullableFloat(lonStr),
	}
	return p, e
}
