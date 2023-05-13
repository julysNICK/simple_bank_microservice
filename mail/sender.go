package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtAuthAddress = "smtp.gmail.com"
	smtpAddr       = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFile []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddr     string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddr string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddr:     fromEmailAddr,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFile []string,
) error {
	e := email.NewEmail()

	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddr)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, file := range attachFile {
		_, err := e.AttachFile(file)

		if err != nil {
			return fmt.Errorf("error when attach file %s: %s", file, err.Error())
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddr, sender.fromEmailPassword, smtAuthAddress)

	return e.Send(smtpAddr, smtpAuth)

}
