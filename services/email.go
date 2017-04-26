package services

import (
	log "github.com/Sirupsen/logrus"
	"net/smtp"
)

type EmailSettings struct {
	Address  string
	Password string
}

var Email EmailSettings

func SendEmail(body string) {
	from := Email.Address
	pass := Email.Password
	to := "chrislewispac@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, Email")
}
