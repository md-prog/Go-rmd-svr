package Models

import (
	"github.com/jmoiron/sqlx/types"
	uuid "github.com/satori/go.uuid"
	null "gopkg.in/guregu/null.v3"
)

// Response is a response type
type Response struct {
	Text   string `json:"text"`
	Errors string `json:"errors"`
}

type Res struct {
	Msg   string         `json:"msg"`
	Token string         `json:"token"`
	Uuid  string         `json:"uuid"`
	Error null.String    `json:"error"`
	Data  types.JSONText `json:"data"`
}

func NewResponse() *Res {
	r := &Res{
		Uuid: uuid.NewV4().String(),
	}
	return r
}
