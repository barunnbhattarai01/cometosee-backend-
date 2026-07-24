package common

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(to, subject, body string) error {
	from := "barunnbhattarai@gmail.com"
	password := os.Getenv("APP_PASSWORD")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf(
		"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n%s",
		subject,
		body,
	))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)
}
