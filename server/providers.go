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

const msgSuccessOnGetUserProviders = "Successfully retrieved providers"
const msgSuccessOnGetProviderByID = "Successfully retrived a provider"

//GetProviderByID ...
func (s *Server) GetProviderByID(c echo.Context) error {
	var provider Models.Provider
	providers := []Models.Provider{}

	providerId := c.Param("id")
	_ = userID(c) // TODO explicit error?

	if err := s.getProvider(providerId, &provider); err != nil {
		res := createProviderErrorResponse(&provider, providers, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createProviderSuccessResponse(&provider, providers, msgSuccessOnGetProviderByID)
	return c.JSON(http.StatusOK, res)
}

func (s *Server) getProvider(id string, f *Models.Provider) error {
	const q = `SELECT * FROM providers WHERE id=$1`
	return s.db.QueryRowx(q, id).StructScan(f)
}

//DeleteProviderByID ...
func (s *Server) DeleteProviderByID(c echo.Context) error {
	var provider Models.Provider
	providers := []Models.Provider{}

	providerId := c.Param("id")
	user_id := userID(c)

	err := s.db.QueryRowx("UPDATE user_providers SET deleted=true WHERE provider_id=$1 AND user_id=$2 RETURNING provider_id as id", providerId, user_id).StructScan(&provider)
	if err != nil {
		handleErr(err)
		res := createProviderErrorResponse(&provider, providers, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createProviderSuccessResponse(&provider, providers, "Successfully deleted that provider")
	return c.JSON(http.StatusOK, res)
}

//GetProviderByID ...
func (s *Server) UpdateProviderByID(c echo.Context) error {
	//TODO
	var provider Models.Provider
	providers := []Models.Provider{}
	res := createProviderErrorResponse(&provider, providers, errMsgUnimplemented)
	return c.JSON(http.StatusOK, res)
}

//GetUserProviders ...
func (s *Server) GetUserProviders(c echo.Context) error {
	providers := []Models.Provider{}
	var provider Models.Provider

	id := userID(c)

	rows, err := s.db.Queryx(`
	SELECT * FROM (
			SELECT uf.provider_id as id, uf.id as user_providers_id, json_agg(ci) as contacts
			FROM user_providers uf
			LEFT OUTER JOIN contact_instance ci
			ON ci.user_providers_id=uf.id
			WHERE uf.user_id=$1
			AND uf.deleted=false
			GROUP BY uf.provider_id, uf.id
	) uf
	LEFT JOIN providers c
	USING (id)`, id)
	if err != nil {
		handleErr(err)
		res := createProviderErrorResponse(&provider, providers, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	for rows.Next() {
		var f Models.Provider
		err := rows.StructScan(&f)
		handleErr(err)
		providers = append(providers, f)
	}
	err = rows.Err()
	if err != nil {
		handleErr(err)
		res := createProviderErrorResponse(&provider, providers, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createProviderSuccessResponse(&provider, providers, msgSuccessOnGetUserProviders)

	return c.JSON(http.StatusOK, res)
}

//PostProvider ...
func (s *Server) PostProvider(c echo.Context) error {
	var provider Models.Provider
	providers := []Models.Provider{}

	// try to get logged in user_id from JWT
	id := userID(c)

	// check if the request body contains valid JSON request
	jsonParsed, err := gabs.ParseJSONBuffer(c.Request().Body)
	if err != nil {
		handleErr(err)
		res := createProviderErrorResponse(&provider, providers, "Invalid provider input.")
		return c.JSON(http.StatusOK, res)
	}

	// extract data from parsed JSON
	hourlyRate, ok := jsonParsed.Path("hourly_rate").Data().(int)
	if !ok || hourlyRate == 0 {
		errMsg := "Cannot create a provider with a blank hourly_rate."
		res := createProviderErrorResponse(&provider, providers, errMsg)
		return c.JSON(http.StatusOK, res)
	}

	// try insert Provider into DB
	fID, err := s.insertProvider(id, hourlyRate)
	if err != nil {
		handleErr(err)
		res := createProviderErrorResponse(&provider, providers, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	// fetch inserted Provider for response
	err = s.db.QueryRowx(`
	SELECT * FROM providers
	WHERE id=$1`, fID).StructScan(&provider)
	if err != nil {
		handleErr(err)
		res := createProviderErrorResponse(&provider, providers, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createProviderSuccessResponse(&provider, providers, "Successfully added your provider")
	return c.JSON(http.StatusOK, res)
}

func (s *Server) insertProvider(userId string, hourlyRate int) (string, error) {
	const insertProvider = `
	INSERT INTO providers(
	hourly_rate)
	VALUES ($1)
	RETURNING id`

	const insertUserProviderRelation = `
	INSERT INTO user_providers(
	user_id, provider_id)
	VALUES ($1, $2);`

	var fID string
	return fID, s.tx(func(tx *sqlx.Tx) error {
		err := tx.Get(&fID, insertProvider, hourlyRate)
		if err != nil {
			return err
		}

		_, err = tx.Exec(insertUserProviderRelation, userId, fID)
		return err
	})
}

func (s *Server) PostProvidersFromCsv(c echo.Context) error {
	//TODO
	var provider Models.Provider
	providers := []Models.Provider{}
	res := createProviderErrorResponse(&provider, providers, errMsgUnimplemented)
	return c.JSON(http.StatusOK, res)
}


type providerData struct {
	Provider   Models.Provider   `json:"provider"`
	Providers []Models.Provider `json:"providers"`
}

func createProviderErrorResponse(provider *Models.Provider, providers []Models.Provider, errMsg string) *Models.Res {
	anon := providerData{
		*provider,
		providers,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = errMsgExists
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createProviderSuccessResponse(provider *Models.Provider, providers []Models.Provider, successMsg string) *Models.Res {
	anon := providerData{
		*provider,
		providers,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}

func (s *Server) getProviderFromDb(providerId string) (*Models.Provider, error) {
	var provider Models.Provider

	err := s.db.QueryRowx(`
	SELECT * FROM providers p
	WHERE p.id=$1`, providerId).StructScan(&provider)
	return &provider, err
}


/**
 this function will insert a provider object into providers DB table

 @return	string	newly inserted person's ID
 		error	Error occured
*/
func (s *Server) insertProviderIntoDb(provider *Models.Provider) (string, error) {
	return "0", nil
	//TODO:
	//
	//fID string
	//return fID, s.tx(func(tx *sqlx.Tx) error {
	//	//err := tx.Get(&fID, insertPersonQuery, provider.FirstName, person.LastName, person.Email,
	//	//	person.Title, provider.CellPhone, provider.HomePhone, person.OfficePhone, person.Fax,
	//	//	person.Address, person.Zip, person.City, person.Lat, person.Lon)
	//	//if err != nil {
	//	//	return err
	//	//}
	//
	//	_, err = tx.Exec(insertUserPersonRelationQuery, userId, fID)
	//	return err
	//})
}

