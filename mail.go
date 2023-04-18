package sti2023

import (
	"net/smtp"
)

type configMail struct {
	SmtpHost string `json:"smtpHost"`
	SmtpPort string `json:"smtpPort"`
	SmtpUser string `json:"smptUser"`
	SmtpPass string `json:"smptPass"`
}
var c configMail
var mailFile string = "mail.json"

func Mail(mail string, code string) bool {
	c.SmtpHost = "michalkukla.xyz"
	c.SmtpPort = "587"
	c.SmtpUser = "michal@michalkukla.xyz"
	c.SmtpPass = ""

	if !ReadJsonFile("./", mailFile, &c) {
		return false
	}
	sendTo := mail
	body := []byte(fmt.Sprintf("From: %s\n"+
		"Ověřovací kód do semestrálního projektu\n\n"+
		"%s", smtpUser, code))

	auth := smtp.PlainAuth("", c.SmtpUser, c.SmtpPass, c.SmtpHost)

	err := smtp.SendMail(c.SmtpHost+":"+c.SmtpPort, auth, c.SmtpUser, sendTo, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("email sent")
	return true
}

