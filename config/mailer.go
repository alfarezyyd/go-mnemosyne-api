package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/wneessen/go-mail"
	"html/template"
)

type MailerService struct {
	ViperConfig *viper.Viper
}

func NewMailerService(viperConfig *viper.Viper) *MailerService {
	return &MailerService{ViperConfig: viperConfig}
}

type EmailPayload struct {
	Title     string
	Recipient string
	Body      string
	Sender    string
}

// SendEmail sends an email with the provided data

func (mailerService *MailerService) SendEmail(recipientEmail string, subjectEmail string, templateFile string, data EmailPayload) error {
	// Parse the HTML template
	fmt.Println("Sending email to ", recipientEmail)
	fmt.Println(mailerService.ViperConfig.GetString("EMAIL_USERNAME"))

	templateFilePath, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Println("ERROR parsing email template file")
		return err
	}
	// Render the template with the data

	var body bytes.Buffer

	if err := templateFilePath.Execute(&body, data); err != nil {
		fmt.Println("ERROR AT 39")
		return err
	}
	// Create a new email message
	mailerMessage := mail.NewMsg()
	fmt.Println(mailerService.ViperConfig.GetString("EMAIL_USERNAME"))
	err = mailerMessage.From(mailerService.ViperConfig.GetString("EMAIL_USERNAME"))
	fmt.Println(err)
	if err != nil {
		return err
	} // Sender email
	fmt.Println("AWdadwadwada1212123123w12312312")

	err = mailerMessage.To(recipientEmail)
	fmt.Println(err)

	if err != nil {
	} // Recipient email
	fmt.Println("adwa")

	mailerMessage.Subject(subjectEmail)                           // Subject
	mailerMessage.SetBodyString(mail.TypeTextHTML, body.String()) // Email body (HTML)

	fmt.Println(mailerService.ViperConfig.GetString("EMAIL_USERNAME"))
	fmt.Println(mailerService.ViperConfig.GetString("EMAIL_PASSWORD"))
	// Set up SMTP client
	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(mailerService.ViperConfig.GetString("EMAIL_USERNAME")),
		mail.WithPassword(mailerService.ViperConfig.GetString("EMAIL_PASSWORD")), // Use App Password if using Gmail
		mail.WithTLSPolicy(mail.DefaultTLSPolicy),
	)

	if err != nil {
		fmt.Println("ERROR AT 58")

		return err
	}

	// Send the email
	if err := client.DialAndSend(mailerMessage); err != nil {
		fmt.Println("ERROR AT 65")

		return err
	}

	return nil
}
