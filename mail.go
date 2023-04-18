package sti2023

import (
	"fmt"
	"net/smtp"
)

type configMail struct {
	SmtpHost string `json:"smtpHost"`
	SmtpPort string `json:"smtpPort"`
	SmtpUser string `json:"smtpUser"`
	SmtpPass string `json:"smtpPass"`
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
	sendTo := []string{mail}
	body := []byte(fmt.Sprintf("From: %s\n"+
		"Ověřovací kód do semestrálního projektu\n\n"+
		"%s", c.SmtpUser, code))

	auth := smtp.PlainAuth("", c.SmtpUser, c.SmtpPass, c.SmtpHost)

	err := smtp.SendMail(c.SmtpHost+":"+c.SmtpPort, auth, c.SmtpUser, sendTo, body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("email sent")
	return true
}
