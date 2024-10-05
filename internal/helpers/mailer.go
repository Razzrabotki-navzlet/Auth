package helpers

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

func SendMail() error {
	// Sender data.
	from := "sashazakirov66@gmail.com"
	password := ""
	to := []string{
		"sasha_zakirov_2014@mail.ru",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "465"

	// Message.
	subject := "Subject: Test Email from Go!\n" // Заголовок
	body := "This is a test email message."
	message := []byte(subject + "\n\n" + body)

	// Set up the TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create a new SMTP client
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Authenticate
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = c.Auth(auth); err != nil {
		fmt.Println(err)
		return err
	}

	// To send the email, we need to specify the sender and the recipient
	if err = c.Mail(from); err != nil {
		fmt.Println(err)
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Send the email body
	w, err := c.Data()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Quit and close the connection
	c.Quit()
	fmt.Println("Email Sent Successfully!")
	return nil

}
