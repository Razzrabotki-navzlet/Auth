package helpers

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendMail Функция работает, приложение отправляет письма но они не доходят до адресата т.к. видимо гугл/вк груп почему-то их блокирует
func SendMail(from string, to string, subject, Link string) error {
	fromMail := mail.NewEmail("Auth service", from)
	toEmail := mail.NewEmail("User", to)
	plainTextContent := fmt.Sprintf("Click here to reset your password: %s", Link)
	htmlContent := fmt.Sprintf("<p>Click <a href='%s'>here</a> to reset your password.</p>", Link)
	fmt.Println(fromMail, "fromMail")
	fmt.Println(subject, "subject")
	fmt.Println(toEmail, "toEmail")
	fmt.Println(plainTextContent, "plainTextContent")
	fmt.Println(htmlContent, "htmlContent")
	message := mail.NewSingleEmail(fromMail, subject, toEmail, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient("api-key")
	response, err := client.Send(message)
	fmt.Println(err)
	if err != nil {
		return err
	}

	fmt.Printf("Status Code: %d\n", response.StatusCode)
	fmt.Printf("Body: %s\n", response.Body)
	fmt.Printf("Headers: %v\n", response.Headers)

	return nil

}
