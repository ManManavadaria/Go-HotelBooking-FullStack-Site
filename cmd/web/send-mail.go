package main

import (
	"fmt"
	"log"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
	"gopkg.in/gomail.v2"
)

// var mail models.MailData

func ReadMail() {
	go func() {
		for {
			mail := <-app.MailChan
			fmt.Println(mail)
			if err := SendEmail(mail); err != nil {
				log.Printf("Error sending email: %v", err)
			}
		}
	}()

}

func SendEmail(mailData models.MailData) error {
	// Create a new message
	m := gomail.NewMessage()

	// Set sender and recipient
	m.SetHeader("From", mailData.From)
	m.SetHeader("To", mailData.To)

	// Set subject and body
	m.SetHeader("Subject", mailData.Subject)
	m.SetBody("text/html", mailData.Data)

	// Create a new SMTP client
	d := gomail.NewDialer("smtp.gmail.com", 587, mailData.From, "Ma@235689")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
