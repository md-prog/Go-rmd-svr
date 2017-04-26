package Controllers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	Models "github.com/chrislewispac/rmd-server/models"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	null "gopkg.in/guregu/null.v3"
	"net/http"
	"time"
)

//GetContractByIDCtrl ...
func GetContractByIDCtrl(c echo.Context) error {
	res := Models.NewResponse()
	res.Msg = "Successfully retrieved contract #1"

	return c.JSON(http.StatusOK, res)
}

//GetContractByIDCtrl ...
func UpdateContractByIDCtrl(c echo.Context) error {
	res := Models.NewResponse()
	res.Msg = "Successfully updated contract #1"

	return c.JSON(http.StatusOK, res)
}

//GetUserContractsCtrl ...
func GetUserContractsCtrl(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	rows, e := Models.DB.Queryx(`
		SELECT * FROM (
			SELECT uc.contract_id as id, uc.id as user_contracts_id, json_agg(cf.user_facilities_id) as facilities, json_agg(ci) as contacts
			FROM user_contracts uc
			LEFT OUTER JOIN contract_facilities cf
			ON cf.contract_id=uc.contract_id
			LEFT OUTER JOIN contact_instance ci
			ON ci.user_contracts_id=uc.id
			WHERE uc.user_id=$1
			AND uc.deleted=false
			GROUP BY uc.contract_id, uc.id
		) ucf
		LEFT JOIN contracts c
		USING (id)`, id)
	if e != nil {
		handleErr(e)
	}

	var contracts []Models.Contract
	for rows.Next() {
		var ct Models.Contract
		err := rows.StructScan(&ct)
		handleErr(err)
		contracts = append(contracts, ct)
	}

	err := rows.Err()
	if err != nil {
		handleErr(err)
	}

	res := Models.NewResponse()
	res.Msg = "Successfully retrieved your contracts"
	if len(contracts) > 0 {
		anon := struct {
			Contracts []Models.Contract `json:"contracts"`
			Contract  []int             `json:"contract"`
		}{
			contracts,
			nil,
		}
		data, _ := json.Marshal(anon)
		res.Data = data
	} else {
		anon := struct {
			Contracts []int `json:"contracts"`
			Contract  []int `json:"contract"`
		}{
			make([]int, 0),
			nil,
		}
		data, _ := json.Marshal(anon)
		res.Data = data
	}

	return c.JSON(http.StatusOK, res)
}

//PostContractCtrl ...
func PostContractCtrl(c echo.Context) error {
	x := c.Get("user").(*jwt.Token)
	claims := x.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	jsonParsed := GetJsonBody(c)

	var rd null.Time
	var sd null.Time

	start_date, ok := jsonParsed.Path("start_date").Data().(string)
	if ok != true {
		sd = null.NewTime(time.Now(), false)
	} else {
		s, _ := time.Parse(time.RFC3339, start_date)
		sd = null.NewTime(s, true)
	}
	renewal_date, ok := jsonParsed.Path("renewal_date").Data().(string)
	if ok != true {
		rd = null.NewTime(time.Now(), false)
	} else {
		r, _ := time.Parse(time.RFC3339, renewal_date)
		rd = null.NewTime(r, true)
	}
	current_contract_status, _ := jsonParsed.Path("current_contract_status").Data().(string)

	ct := &Models.Contract{
		StartDate:   sd,
		RenewalDate: rd,
	}

	log.Println(current_contract_status)

	insertContract := `
		INSERT INTO contracts(
		start_date)
		VALUES ($1)
		RETURNING id`

	insertUserContractRelation := `
		INSERT INTO user_contracts(
		user_id, contract_id)
		VALUES ($1, $2);
	`

	// insertContractFacilitiesRelation := `
	// 	INSERT INTO contracts_facilities
	// 		( contract_id
	// 		, facility_id)
	// 	VALUES
	// 		( $1
	// 		, $2
	// 		)`

	tx, err := Models.DB.Beginx()
	var ctID string
	err = tx.Get(&ctID, insertContract, ct.StartDate)
	handleErr(err)

	_, err = tx.Exec(insertUserContractRelation, id, ctID)
	handleErr(err)

	err = tx.Commit()
	handleErr(err)

	if err != nil {
		log.Println("Error was not nil, return something meaningful") //TODO
	}

	var contract Models.Contract
	err = Models.DB.QueryRowx(`
		SELECT * FROM contracts
		WHERE id=$1`, ctID).StructScan(&contract)
	if err != nil {
		handleErr(err)
	}

	contracts := []Models.Contract{}
	anon := struct {
		Contract  Models.Contract   `json:"contract"`
		Contracts []Models.Contract `json:"contracts"`
	}{
		contract,
		contracts,
	}

	data, _ := json.Marshal(anon)

	res := Models.NewResponse()
	res.Msg = "Successfully posted contract"
	res.Data = data

	return c.JSON(http.StatusOK, res)
}

func PostContractsFromCsvCtrl(c echo.Context) error {

	res := Models.NewResponse()
	res.Msg = "Successfully uploaded your contracts from csv file"

	return c.JSON(http.StatusOK, res)
}
