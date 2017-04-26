package Controllers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
	Models "github.com/chrislewispac/rmd-server/models"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

//GetFacilityByIDCtrl ...
func GetFacilityByIDCtrl(c echo.Context) error {
	res := Models.NewResponse()
	res.Msg = "Successfully retrieved facility #1"

	return c.JSON(http.StatusOK, res)
}

//DeleteFacilityByIDCtrl ...
func DeleteFacilityByIDCtrl(c echo.Context) error {
	facilities := []Models.Facility{}
	var facility Models.Facility
	facility_id := c.Param("id")
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	user_id := claims["id"].(string)

	err := Models.DB.QueryRowx("UPDATE user_facilities SET deleted=true WHERE facility_id=$1 AND user_id=$2 RETURNING facility_id as id", facility_id, user_id).StructScan(&facility)
	if err != nil {
		handleErr(err)
	}

	res := createFacilitySuccessResponse(facility, facilities, "Successfully delete that facility")
	return c.JSON(http.StatusOK, res)
}

//GetFacilityByIDCtrl ...
func UpdateFacilityByIDCtrl(c echo.Context) error {
	res := Models.NewResponse()
	res.Msg = "Successfully updated facility #1"

	return c.JSON(http.StatusOK, res)
}

//GetUserFacilitiesCtrl ...
func GetUserFacilitiesCtrl(c echo.Context) error {
	var facilities []Models.Facility
	var facility Models.Facility
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	rows, e := Models.DB.Queryx(`
		SELECT * FROM (
				SELECT uf.facility_id as id, uf.id as user_facilities_id, json_agg(ci) as contacts
				FROM user_facilities uf
				LEFT OUTER JOIN contact_instance ci
				ON ci.user_facilities_id=uf.id
				WHERE uf.user_id=$1
				AND uf.deleted=false
				GROUP BY uf.facility_id, uf.id
		) uf
		LEFT JOIN facilities c
		USING (id)`, id)
	if e != nil {
		log.Println(e)
	}

	for rows.Next() {
		var f Models.Facility
		err := rows.StructScan(&f)
		handleErr(err)
		facilities = append(facilities, f)
	}

	err := rows.Err()
	if err != nil {
		handleErr(err)
	}

	res := createFacilitySuccessResponse(facility, facilities, "Retrieved Facilities")

	return c.JSON(http.StatusOK, res)
}

//PostFacilityCtrl ...
func PostFacilityCtrl(c echo.Context) error {
	x := c.Get("user").(*jwt.Token)
	claims := x.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(b))

	var facility Models.Facility
	facilities := []Models.Facility{}

	name, ok := jsonParsed.Path("name").Data().(string)
	if name == "" {
		errMsg := "Cannot create a facility with a blank name."
		res := createFacilityErrorResponse(facility, facilities, errMsg)

		return c.JSON(http.StatusOK, res)
	}
	if ok == false {
		//this would be an error "required field" TODO
	}

	f := &Models.Facility{
		Name: name,
	}

	insertFacility := `
		INSERT INTO facilities(
		name)
		VALUES ($1)
		RETURNING id`

	insertUserFacilityRelation := `
		INSERT INTO user_facilities(
		user_id, facility_id)
		VALUES ($1, $2);
	`

	tx, err := Models.DB.Beginx()
	var fID string
	err = tx.Get(&fID, insertFacility, f.Name)
	handleErr(err)

	_, err = tx.Exec(insertUserFacilityRelation, id, fID)
	handleErr(err)

	err = tx.Commit()
	handleErr(err)

	if err != nil {
		log.Println("Error was not nil, return something meaningful") //TODO
	}

	err = Models.DB.QueryRowx(`
		SELECT * FROM facilities
		WHERE id=$1`, fID).StructScan(&facility)
	if err != nil {
		handleErr(err)
	}

	anon := struct {
		Facility   Models.Facility   `json:"facility"`
		Facilities []Models.Facility `json:"facilities"`
	}{
		facility,
		facilities,
	}

	data, _ := json.Marshal(anon)

	res := Models.NewResponse()
	res.Msg = "Successfully posted contract"
	res.Data = data

	return c.JSON(http.StatusOK, res)
}

func PostFacilitiesFromCsvCtrl(c echo.Context) error {

	res := Models.NewResponse()
	res.Msg = "Successfully uploaded your facilities from csv file"

	return c.JSON(http.StatusOK, res)
}

func createFacilityErrorResponse(facility Models.Facility, facilities []Models.Facility, errMsg string) *Models.Res {
	anon := struct {
		Facility   Models.Facility   `json:"facility"`
		Facilities []Models.Facility `json:"facilities"`
	}{
		facility,
		facilities,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = "There was an Error"
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createFacilitySuccessResponse(facility Models.Facility, facilities []Models.Facility, successMsg string) *Models.Res {
	anon := struct {
		Facility   Models.Facility   `json:"facility"`
		Facilities []Models.Facility `json:"facilities"`
	}{
		facility,
		facilities,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}
