package helpers

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

func SendMail(to []string, subject string, body string) error {
	from := "sashazakirov66@gmail.com"
	password := "sfdn azzm qdjo zxga"

	smtpHost := "smtp.gmail.com"
	smtpPort := "465"

	subject = subject + "\n"
	message := []byte(subject + "\n\n" + body)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = c.Auth(auth); err != nil {
		fmt.Println(err)
		return err
	}

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

	c.Quit()
	fmt.Println("Email Sent Successfully!")
	return nil

}
