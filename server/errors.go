package server

import (
	"strings"

	"gitlab.com/chrislewispac/rmd-server/models"
	null "gopkg.in/guregu/null.v3"

	log "github.com/Sirupsen/logrus"
)

const (
	//TODO use these from contracts and auth
	errMsgBadRequest     = "Uh oh, we could understand the format of your request. Please let us know in our chat about this error!"
	errMsgSomethingWrong = "Something Went Wrong!"
	errMsgUnimplemented  = "Not Yet Implemented!"
	errMsgExists         = "There was an Error"
)

func errResponding() *Models.Res {
	res := Models.NewResponse()
	res.Msg = "There was an Error"
	res.Error = null.StringFrom("Something went wrong building the response")
	return res
}

func handleErr(err error) {
	if err != nil {
		log.Warn(err)
	}
	return
}

func errMissingRequiredFields(fields ...string) string {
	return "Missing required fields: " + strings.Join(fields, ", ")
}
