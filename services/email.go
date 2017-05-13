package services

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	sp "github.com/SparkPost/gosparkpost"
)

//EmailSettings TODO
type EmailSettings struct {
	BaseUrl string
	SPKey   string
	Client  *sp.Client
}

//InitEmailService TODO
func InitEmailService(baseURL string, apiKey string) *EmailSettings {
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}
	var client sp.Client
	err := client.Init(cfg)
	if err != nil {
		log.Error("SparkPost client init failed:")
	}
	e := &EmailSettings{BaseUrl: baseURL, SPKey: apiKey, Client: &client}
	return e
}

//SendResetPasswordEmail TODO
func (EmailService *EmailSettings) SendResetPasswordEmail(email string, body string) {

	link := fmt.Sprintf("%s/#ResetPassword/%s", EmailService.BaseUrl, body)
	html := fmt.Sprintf("<a href=%s>Click Here To Reset Your Password</a>", link)

	// Create a Transmission using an inline Recipient List
	// and inline email Content.
	tx := &sp.Transmission{
		Recipients: []string{email},
		Content: sp.Content{
			HTML:    html,
			From:    "admin@statrecruit.com",
			Subject: "Password Reset link",
		},
	}
	id, _, err := EmailService.Client.Send(tx)
	if err != nil {
		log.Error(err)
	}

	// The second value returned from Send
	// has more info about the HTTP response, in case
	// you'd like to see more than the Transmission id.
	log.Printf("Transmission sent with id [%s]\n", id)
}

//SendNoEmailForwardingSetupEmail TODO
func (EmailService *EmailSettings) SendNoEmailForwardingSetupEmail(email string, body string) {

	html := fmt.Sprintf("%s", body)

	tx := &sp.Transmission{
		Recipients: []string{email},
		Content: sp.Content{
			HTML:    html,
			From:    "admin@statrecruit.com",
			Subject: "Email Forwarding Failed [ Not Setup ]",
		},
	}
	id, _, err := EmailService.Client.Send(tx)
	if err != nil {
		log.Error(err)
	}

	log.Printf("Transmission sent with id [%s]\n", id)
}
