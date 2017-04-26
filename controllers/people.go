package Controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"

	Models "github.com/chrislewispac/rmd-server/models"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

//GetPersonByIDCtrl exported
func GetPersonByIDCtrl(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	contact_id := c.Param("id")

	var person Models.Person
	err := Models.DB.QueryRowx(`
		SELECT * FROM (
					SELECT people_id AS id
					FROM user_contacts uc
					WHERE uc.user_id=$1
				) u
		INNER JOIN people c
		USING (id)
		WHERE u.id=$2`, id, contact_id).StructScan(&person)
	if err != nil {
		handleErr(err)
	}

	anon := struct {
		Person Models.Person `json:"person"`
	}{
		person,
	}

	data, _ := json.Marshal(anon)

	res := Models.NewResponse()
	res.Msg = "Successfully retrieved contact #1"
	res.Data = data

	return c.JSON(http.StatusOK, res)
}

//GetPersonByIDCtrl exported
func UpdatePersonByIDCtrl(c echo.Context) error {

	res := Models.NewResponse()
	res.Msg = "Successfully updated contact #1"

	return c.JSON(http.StatusOK, res)
}

func PostContactCtrl(c echo.Context) error {

	res := Models.NewResponse()
	res.Msg = "Successfully added your contact"

	return c.JSON(http.StatusOK, res)
}

//GetUserPersonsCtrl ...
func GetUserPersonsCtrl(c echo.Context) error {
	people := []Models.Person{}
	var person Models.Person
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	rows, e := Models.DB.Queryx(`
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
		`, id)
	if e != nil {
		log.Println(e)
	}

	for rows.Next() {
		var ct Models.Person
		err := rows.StructScan(&ct)
		handleErr(err)
		people = append(people, ct)
	}

	err := rows.Err()
	if err != nil {
		handleErr(err)
	}

	res := createPersonSuccessResponse(person, people, "Successfully retrieved your contacts")

	return c.JSON(http.StatusOK, res)
}

func PostPersonsFromCsvCtrl(c echo.Context) error {
	people := []Models.Person{}
	var person Models.Person
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	jsonParsed := GetJsonBody(c)

	children, _ := jsonParsed.Children()
	tx, err := Models.DB.Beginx()
	contact := `
	INSERT INTO people
		( first_name
		, last_name
		, middle_name
		, title
		, email
	 	, cell_phone
	 	, home_phone
	 	, office_phone
	  , fax
		, address
		, city
		, zip
		, state_id
		,	recruiter_id
		,	provider_id
		, notes
		, lat
		, lon )
	VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	RETURNING id`
	contact_rel := `
		INSERT INTO user_contacts
			( people_id
			, user_id)
		VALUES ($1, $2)`

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
		}

		if lat, ok := child.S("lon").Data().(float64); ok {
			ct.Lat = null.FloatFrom(lat)
		} else {
			ct.Lat = null.Float{}
		}

		if lon, ok := child.S("lon").Data().(float64); ok {
			ct.Lon = null.FloatFrom(lon)
		} else {
			ct.Lon = null.Float{}
		}

		var ctID string
		err = tx.Get(&ctID, contact, ct.FirstName, ct.LastName, ct.MiddleName,
			ct.Title, ct.Email, ct.CellPhone, ct.HomePhone, ct.OfficePhone, ct.Fax,
			ct.Address, ct.City, ct.Zip, ct.StateID, ct.RecruiterID, ct.ProviderID, ct.Notes, ct.Lat, ct.Lon)
		handleErr(err)
		_, err = tx.Exec(contact_rel, ctID, id)
		handleErr(err)
	}
	err = tx.Commit()
	handleErr(err)

	rows, e := Models.DB.Queryx(`
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

	res := createPersonSuccessResponse(person, people, "Successfully added people from CSV")

	return c.JSON(http.StatusOK, res)
}

func PostPersonCtrl(c echo.Context) error {
	people := []Models.Person{}
	var person Models.Person
	x := c.Get("user").(*jwt.Token)
	claims := x.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	jsonParsed := GetJsonBody(c)

	first_name, ok := jsonParsed.Path("first_name").Data().(string)
	if ok == false {
		//this would be an error "required field" TODO
	}

	middle_name, ok := jsonParsed.Path("middle_name").Data().(string)
	if ok == false {
		//this is not a "required field" TODO
	}

	last_name, ok := jsonParsed.Path("last_name").Data().(string)
	if ok == false {
		//this would be an error "required field" TODO
	}

	email, ok := jsonParsed.Path("email").Data().(string)
	if ok == false {
		//this would be an error "required field" TODO
	}

	p := &Models.Person{
		FirstName:  first_name,
		MiddleName: null.StringFrom(middle_name),
		LastName:   last_name,
		Email:      email,
	}

	insertPerson := `
		INSERT INTO people(
		first_name, middle_name, last_name, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	insertUserPersonRelation := `
		INSERT INTO user_contacts(
		user_id, people_id)
		VALUES ($1, $2);
	`

	tx, err := Models.DB.Beginx()
	var pID string
	err = tx.Get(&pID, insertPerson, p.FirstName, p.MiddleName, p.LastName, p.Email)
	handleErr(err)

	_, err = tx.Exec(insertUserPersonRelation, id, pID)
	handleErr(err)

	err = tx.Commit()
	handleErr(err)

	if err != nil {
		log.Println("Error was not nil, return something meaningful") //TODO
	}

	err = Models.DB.QueryRowx(`
		SELECT * FROM people
		WHERE id=$1`, pID).StructScan(&person)
	if err != nil {
		handleErr(err)
	}

	anon := struct {
		Person Models.Person   `json:"person"`
		People []Models.Person `json:"people"`
	}{
		person,
		people,
	}

	data, _ := json.Marshal(anon)

	res := Models.NewResponse()
	res.Msg = "Successfully posted contract"
	res.Data = data

	return c.JSON(http.StatusOK, res)
}

func createPersonErrorResponse(person Models.Person, people []Models.Person, errMsg string) *Models.Res {
	anon := struct {
		Person Models.Person   `json:"person"`
		People []Models.Person `json:"people"`
	}{
		person,
		people,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = "There was an Error"
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createPersonSuccessResponse(person Models.Person, people []Models.Person, successMsg string) *Models.Res {
	anon := struct {
		Person Models.Person   `json:"person"`
		People []Models.Person `json:"people"`
	}{
		person,
		people,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}
