package services

import (
	"bytes"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/labstack/gommon/log"
	"html/template"
	"net/smtp"
	"os"
	"time"
)

type Email struct{}

func NewEmailService() *Email {
	return &Email{}
}

func (instance Email) SendUserAccountActivationEmail(userData user.User) error {
	subject := "Bem-vindo(a) ao Você na Câmara!"

	templatePath := "core/services/resources/email_template.html"

	emailData := map[string]interface{}{
		"user_name":       fmt.Sprintf("%s %s", userData.FirstName(), userData.LastName()),
		"activation_code": userData.ActivationCode(),
		"current_year":    time.Now().Year(),
	}

	err := sendEmail(subject, userData.Email(), templatePath, emailData)
	if err != nil {
		log.Errorf("Erro ao enviar email de ativação da conta do usuário %s: %s", userData.Email(), err)
		return err
	}

	return nil
}

func sendEmail(subject string, to string, templatePath string, emailData interface{}) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUserEmail := os.Getenv("SMTP_USER_EMAIL")
	smtpUserPassword := os.Getenv("SMTP_USER_PASSWORD")
	auth := smtp.PlainAuth("", smtpUserEmail, smtpUserPassword, smtpHost)

	emailTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Errorf("Erro ao carregar o template do email '%s' com o caminho %s: %s", subject, templatePath, err)
		return err
	}

	var body bytes.Buffer
	body.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n")
	body.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	body.WriteString(fmt.Sprintf("To: %s\n", to))
	body.WriteString(fmt.Sprintf("From: %s\n\n", smtpUserEmail))

	err = emailTemplate.Execute(&body, emailData)
	if err != nil {
		log.Errorf("Erro ao executar o template do email '%s': %s", subject, err)
		return err
	}

	smtpAddress := fmt.Sprint(smtpHost, ":", smtpPort)
	err = smtp.SendMail(smtpAddress, auth, smtpUserEmail, []string{to}, body.Bytes())
	if err != nil {
		log.Errorf("Erro ao enviar email '%s' para o usuário %s: %s", subject, to, err.Error())
		return err
	}

	log.Infof("Email '%s' enviado para o usuário %s", subject, to)
	return nil
}
