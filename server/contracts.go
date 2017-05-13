package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	null "gopkg.in/guregu/null.v3"

	Models "gitlab.com/chrislewispac/rmd-server/models"
)

//GetContractByID ...
func (s *Server) GetContractByID(c echo.Context) error {
	contracts := []Models.Contract{}
	var contract Models.Contract

	uID := userID(c)
	contractID := c.Param("id")

	if err := s.getContract(uID, contractID, &contract); err != nil {
		handleErr(err)
		res := createContractErrorResponse(contract, contracts, "Something Went Wrong!")
		return c.JSON(http.StatusOK, res)
	}

	res := createContractSuccessResponse(contract, contracts, "Successfully retrieved your contract")

	return c.JSON(http.StatusOK, res)
}

//DeleteContractByID ...
func (s *Server) DeleteContractByID(c echo.Context) error {
	contracts := []Models.Contract{}
	var contract Models.Contract

	contractID := c.Param("id")
	uID := userID(c)

	err := s.db.QueryRowx("UPDATE user_contracts SET deleted=true WHERE contract_id=$1 AND user_id=$2 RETURNING contract_id as id", contractID, uID).StructScan(&contract)
	if err != nil {
		handleErr(err)
		res := createContractErrorResponse(contract, contracts, "Something Went Wrong!")
		return c.JSON(http.StatusOK, res)
	}

	res := createContractSuccessResponse(contract, contracts, "Successfully deleted that contract")
	return c.JSON(http.StatusOK, res)
}

//UpdateContractByID ...
func (s *Server) UpdateContractByID(c echo.Context) error {
	contracts := []Models.Contract{}
	var contract Models.Contract
	//TODO implement
	res := createContractErrorResponse(contract, contracts, "Not yet implemented!")
	return c.JSON(http.StatusOK, res)
}

//GetUserContracts ...
func (s *Server) GetUserContracts(c echo.Context) error {
	var contract Models.Contract
	var contracts []Models.Contract

	id := userID(c)

	rows, e := s.db.Queryx(`
	SELECT * FROM (
		SELECT uc.contract_id as id, uc.id as user_contracts_id, json_agg(uf.facility_id) as facilities, json_agg(ci) as contacts
		FROM user_contracts uc
		LEFT JOIN contract_facilities cf
		ON cf.contract_id=uc.contract_id
		LEFT JOIN contact_instance ci
		ON ci.user_contracts_id=uc.id
		LEFT JOIN user_facilities uf
		ON cf.user_facilities_id=uf.id
		WHERE uc.user_id=$1
		AND uc.deleted=false
		GROUP BY uc.contract_id, uc.id
	) ucf
	LEFT JOIN contracts c
	USING (id)`, id)
	if e != nil {
		handleErr(e)
		res := createContractErrorResponse(contract, contracts, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	for rows.Next() {
		var ct Models.Contract
		err := rows.StructScan(&ct)
		if err != nil {
			handleErr(err)
			continue
		}
		contracts = append(contracts, ct)
	}
	err := rows.Err()
	if err != nil {
		handleErr(err)
		res := createContractErrorResponse(contract, contracts, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	res := createContractSuccessResponse(contract, contracts, "Successfully retrieved your contracts")
	return c.JSON(http.StatusOK, res)
}

//PostContract ...
func (s *Server) PostContract(c echo.Context) error {
	contracts := []Models.Contract{}
	var contract Models.Contract
	x := c.Get("user").(*jwt.Token)
	claims := x.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
	}
	jsonParsed := GetJsonBody(string(b))

	var rd null.Time
	var sd null.Time

	startDate, ok := jsonParsed.Path("start_date").Data().(string)
	if ok != true {
		sd = null.NewTime(time.Now(), false)
	} else {
		s, _ := time.Parse(time.RFC3339, startDate)
		sd = null.NewTime(s, true)
	}
	renewalDate, ok := jsonParsed.Path("renewal_date").Data().(string)
	if ok != true {
		rd = null.NewTime(time.Now(), false)
	} else {
		r, _ := time.Parse(time.RFC3339, renewalDate)
		rd = null.NewTime(r, true)
	}
	// current_contract_status, _ := jsonParsed.Path("current_contract_status").Data().(string)
	contractName, ok := jsonParsed.Path("name").Data().(string)
	if ok == false {
		res := createContractErrorResponse(contract, contracts, "There was an error parsing the name field")
		return c.JSON(http.StatusOK, res)
	}

	ct := &Models.Contract{
		StartDate:   sd,
		RenewalDate: rd,
		Name:        contractName,
	}

	if ct.Name == "" {
		res := createContractErrorResponse(contract, contracts, "Contract name is required!")
		return c.JSON(http.StatusOK, res)
	}

	ctID, _, err := s.insertContract(id, ct)
	if err != nil {
		handleErr(err)
		res := createContractErrorResponse(contract, contracts, "Something Went Wrong!")
		return c.JSON(http.StatusOK, res)
	}

	if err := s.getContract(id, ctID, &contract); err != nil {
		handleErr(err)
		res := createContractErrorResponse(contract, contracts, "Something Went Wrong!")
		return c.JSON(http.StatusOK, res)
	}

	res := createContractSuccessResponse(contract, contracts, "Successfully added your contract")

	return c.JSON(http.StatusOK, res)
}

func (s *Server) getContract(userID, contractID string, c *Models.Contract) error {
	return s.db.QueryRowx(`
	SELECT * FROM (
		SELECT uc.contract_id as id, uc.id as user_contracts_id, json_agg(uf.facility_id) as facilities, json_agg(ci) as contacts
		FROM user_contracts uc
		LEFT JOIN contract_facilities cf
		ON cf.contract_id=uc.contract_id
		LEFT JOIN contact_instance ci
		ON ci.user_contracts_id=uc.id
		LEFT JOIN user_facilities uf
		ON cf.user_facilities_id=uf.id
		WHERE uc.user_id=$1 AND uc.contract_id=$2
		AND uc.deleted=false
		GROUP BY uc.contract_id, uc.id
	) ucf
	LEFT JOIN contracts c
	USING (id)`, userID, contractID).StructScan(c)
}

// insertContract inserts a contract and user_contract, returning their IDs.
func (s *Server) insertContract(userID string, ct *Models.Contract) (string, string, error) {
	log.Println("inserting: ", ct)

	const insertContract = `
	INSERT INTO contracts(
	name, start_date, renewal_date, file)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	const insertUserContractRelation = `
	INSERT INTO user_contracts(
	user_id, contract_id)
	VALUES ($1, $2) RETURNING id;`

	const selectUserFacility = `SELECT id FROM user_facilities
	WHERE user_id=$1 AND facility_id=$2;`

	const insertContractFacilitiesRelation = `INSERT INTO contract_facilities(user_facilities_id, contract_id) VALUES ($1, $2);`

	var ctID string
	var ucID string
	return ctID, ucID, s.tx(func(tx *sqlx.Tx) error {
		err := tx.Get(&ctID, insertContract, ct.Name, ct.StartDate, ct.RenewalDate, ct.File)
		if err != nil {
			return err
		}
		err = tx.Get(&ucID, insertUserContractRelation, userID, ctID)
		if err != nil {
			return err
		}

		if ct.Facilities.Valid {
			var facilities []string
			if err := ct.Facilities.Unmarshal(&facilities); err != nil {
				return err
			}
			for i := range facilities {
				var userFacID string
				err := tx.Get(&userFacID, selectUserFacility, userID, facilities[i])
				if err == sql.ErrNoRows {
					return fmt.Errorf("no facility %q for user %q", facilities[i], userID)
				}
				if err != nil {
					return err
				}

				_, err = tx.Exec(insertContractFacilitiesRelation, userFacID, ctID)
				if err != nil {
					return err
				}
			}
		}

		//TODO insert contacts

		return nil
	})
}

//PostContractsFromCsv TODO
func (s *Server) PostContractsFromCsv(c echo.Context) error {
	contracts := []Models.Contract{}
	var contract Models.Contract
	//TODO implement
	res := createContractErrorResponse(contract, contracts, "Not yet implemented!")
	return c.JSON(http.StatusOK, res)
}

func createContractErrorResponse(contract Models.Contract, contracts []Models.Contract, errMsg string) *Models.Res {
	data, err := json.Marshal(contractData{
		&contract,
		contracts,
	})
	if err != nil {
		handleErr(err)
		return errResponding()
	}
	res := Models.NewResponse()
	res.Msg = errMsgExists
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createContractSuccessResponse(contract Models.Contract, contracts []Models.Contract, successMsg string) *Models.Res {
	data, err := json.Marshal(contractData{
		&contract,
		contracts,
	})
	if err != nil {
		handleErr(err)
		return errResponding()
	}
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}

type contractData struct {
	Contract  *Models.Contract  `json:"contract"`
	Contracts []Models.Contract `json:"contracts"`
}
