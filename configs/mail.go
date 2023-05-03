package configs

import (
	"os"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/gomail.v2"
)

func SendMail(from string, to string, subject string, htmlBody string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(os.Getenv("SMTP_SERVER"), smtpPort, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PW"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
