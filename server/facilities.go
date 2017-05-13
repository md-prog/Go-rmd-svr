package server

import (
	"encoding/json"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	null "gopkg.in/guregu/null.v3"

	Models "gitlab.com/chrislewispac/rmd-server/models"
)

//GetFacilityByID ...
func (s *Server) GetFacilityByID(c echo.Context) error {
	var facility Models.Facility
	facilities := []Models.Facility{}

	facilityId := c.Param("id")

	if err := s.getFacility(facilityId, &facility); err != nil {
		res := createFacilityErrorResponse(facility, facilities, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createFacilitySuccessResponse(facility, facilities, "Successfully retrieved your facility")
	return c.JSON(http.StatusOK, res)
}

func (s *Server) getFacility(id string, f *Models.Facility) error {
	//TODO proper join
	const q = `SELECT * FROM facilities WHERE id=$1`
	return s.db.QueryRowx(q, id).StructScan(f)
}

//DeleteFacilityByID ...
func (s *Server) DeleteFacilityByID(c echo.Context) error {
	var facility Models.Facility
	facilities := []Models.Facility{}

	facility_id := c.Param("id")
	user_id := userID(c)

	err := s.db.QueryRowx("UPDATE user_facilities SET deleted=true WHERE facility_id=$1 AND user_id=$2 RETURNING facility_id as id", facility_id, user_id).StructScan(&facility)
	if err != nil {
		handleErr(err)
		res := createFacilityErrorResponse(facility, facilities, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createFacilitySuccessResponse(facility, facilities, "Successfully deleted that facility")
	return c.JSON(http.StatusOK, res)
}

//GetFacilityByID ...
func (s *Server) UpdateFacilityByID(c echo.Context) error {
	//TODO
	var facility Models.Facility
	facilities := []Models.Facility{}
	res := createFacilityErrorResponse(facility, facilities, errMsgUnimplemented)
	return c.JSON(http.StatusOK, res)
}

//GetUserFacilities ...
func (s *Server) GetUserFacilities(c echo.Context) error {
	facilities := []Models.Facility{}
	var facility Models.Facility

	id := userID(c)

	rows, err := s.db.Queryx(`
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
	if err != nil {
		handleErr(err)
		res := createFacilityErrorResponse(facility, facilities, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	for rows.Next() {
		var f Models.Facility
		err := rows.StructScan(&f)
		handleErr(err)
		facilities = append(facilities, f)
	}
	err = rows.Err()
	if err != nil {
		handleErr(err)
		res := createFacilityErrorResponse(facility, facilities, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createFacilitySuccessResponse(facility, facilities, "Retrieved Facilities")

	return c.JSON(http.StatusOK, res)
}

//PostFacility ...
func (s *Server) PostFacility(c echo.Context) error {
	var facility Models.Facility
	facilities := []Models.Facility{}

	id := userID(c)

	jsonParsed, err := gabs.ParseJSONBuffer(c.Request().Body)
	if err != nil {
		handleErr(err)
		res := createFacilityErrorResponse(facility, facilities, "Invalid facility input.")
		return c.JSON(http.StatusOK, res)
	}

	name, ok := jsonParsed.Path("name").Data().(string)
	if !ok || name == "" {
		errMsg := "Cannot create a facility with a blank name."
		res := createFacilityErrorResponse(facility, facilities, errMsg)
		return c.JSON(http.StatusOK, res)
	}

	fID, err := s.insertFacility(id, name)
	if err != nil {
		handleErr(err)
		res := createFacilityErrorResponse(facility, facilities, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	if err := s.getFacility(fID, &facility); err != nil {
		handleErr(err)
		res := createFacilityErrorResponse(facility, facilities, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createFacilitySuccessResponse(facility, facilities, "Successfully added your facility")
	return c.JSON(http.StatusOK, res)
}

// insertFacility inserts a new facility and returns the facility id.
func (s *Server) insertFacility(userID, facilityName string) (string, error) {
	const insertFacility = `
	INSERT INTO facilities(
	name)
	VALUES ($1)
	RETURNING id`

	const insertUserFacilityRelation = `
	INSERT INTO user_facilities(
	user_id, facility_id)
	VALUES ($1, $2);`

	var fID string
	return fID, s.tx(func(tx *sqlx.Tx) error {
		err := tx.Get(&fID, insertFacility, facilityName)
		if err != nil {
			return err
		}
		_, err = tx.Exec(insertUserFacilityRelation, userID, fID)
		return err
	})
}

func (s *Server) PostFacilitiesFromCsv(c echo.Context) error {
	//TODO
	var facility Models.Facility
	facilities := []Models.Facility{}
	res := createFacilityErrorResponse(facility, facilities, errMsgUnimplemented)
	return c.JSON(http.StatusOK, res)
}

func createFacilityErrorResponse(facility Models.Facility, facilities []Models.Facility, errMsg string) *Models.Res {
	anon := facilityData{
		facility,
		facilities,
	}
	data, err := json.Marshal(anon)
	if err != nil {
		handleErr(err)
		return errResponding()
	}
	res := Models.NewResponse()
	res.Msg = "There was an Error"
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createFacilitySuccessResponse(facility Models.Facility, facilities []Models.Facility, successMsg string) *Models.Res {
	anon := facilityData{
		facility,
		facilities,
	}
	data, err := json.Marshal(anon)
	if err != nil {
		handleErr(err)
		return errResponding()
	}
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}

type facilityData struct {
	Facility   Models.Facility   `json:"facility"`
	Facilities []Models.Facility `json:"facilities"`
}
