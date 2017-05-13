package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	Models "gitlab.com/chrislewispac/rmd-server/models"
)

//ReceiveEmailCV exported
func (s *Server) ReceiveEmailCV(c echo.Context) error {
	var userHandle string
	var fromAddress string
	var user Models.User

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
	}
	jsonParsed, err := gabs.ParseJSON([]byte(b))
	if err != nil {
		handleErr(err)
	}

	children, _ := jsonParsed.Children()
	for _, child := range children {
		recipient, ok := child.Path("msys.relay_message.rcpt_to").Data().(string)
		if !ok {
			log.Error("error parsing recipient field.")
		}
		fromAddress, _ = child.Path("msys.relay_message.friendly_from").Data().(string)
		userHandle = strings.Split(recipient, "@")[0]
	}

	err = s.db.QueryRowx(`SELECT * FROM users WHERE email_forwarding_handle=$1`, userHandle).StructScan(&user)
	if err != nil {
		s.email.SendNoEmailForwardingSetupEmail(fromAddress, fmt.Sprintf("Unfortunately, there is no user associated with %s. If you would like to activate email forwarding please login to your settings page.", userHandle))
		log.Error("No User Found... Error Email Sent.")
	}

	return c.JSON(http.StatusOK, "ok")
}
