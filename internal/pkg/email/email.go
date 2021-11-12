package email

import (
	gomail "gopkg.in/mail.v2"
)

func SendEmail(fromEmail, fromPassword, toEmail, subject, body string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", fromEmail)
	message.SetHeader("To", toEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, fromPassword)

	return dailer.DialAndSend(message)
}
