package email

import (
	"cat_ben/src/config"
	"gopkg.in/gomail.v2"
	"log"
)

func Send(subject string, body string) {
	from := config.Config.EmailFrom
	to := config.Config.EmailTo
	authCode := config.Config.EmailAuth
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	d := gomail.NewDialer("smtp.163.com", 465, from, authCode)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		log.Printf("send email err:%s,%s", subject, err)
	}
}
