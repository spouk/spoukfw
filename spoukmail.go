package spoukfw

import (
	"fmt"
	"errors"
	"crypto/tls"
	gml "gopkg.in/gomail.v2"
)

type (
	SpoukMail struct {
		MailMessage MailMessage
	}
	MailMessage struct {
		To         string
		From       string
		Message    string
		Subject    string
		FileAttach string `fullpath to file`
		Host       string
		Port       int
		Username   string
		Password   string
	}

)
func newSpoukMail() *SpoukMail {
	return &SpoukMail{MailMessage:MailMessage{}}
}
func (mail MailMessage) MailSend(message *MailMessage) (error) {
	d := gml.NewPlainDialer(message.Host, message.Port, message.Username, message.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gml.NewMessage()
	m.SetHeader("From", message.From)
	m.SetHeader("To", message.To)
	//	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/html", message.Message)
	if message.FileAttach != "" {
		m.Attach(message.FileAttach)
	}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("[sendemail] ошибка отправки сообщения %v\n", err)
		return errors.New(fmt.Sprintf("[sendemail] ошибка отправки сообщения %v\n", err))
	}
	return nil
}
