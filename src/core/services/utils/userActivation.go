package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/labstack/gommon/log"
	"html/template"
	"math/big"
	"net/smtp"
	"os"
	"time"
)

func GenerateUserActivationCode() (string, error) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	activationCode := make([]byte, length)
	for index := range activationCode {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Errorf("Erro ao gerar código de ativação da conta do usuário: ", err.Error())
			return "", err
		}

		activationCode[index] = charset[randomInt.Int64()]
	}

	return string(activationCode), nil
}

func SendActivationEmail(userData user.User) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUserEmail := os.Getenv("SMTP_USER_EMAIL")
	smtpUserPassword := os.Getenv("SMTP_USER_PASSWORD")
	auth := smtp.PlainAuth("", smtpUserEmail, smtpUserPassword, smtpHost)

	type EmailData struct {
		UserName       string
		ActivationCode string
		CurrentYear    int
	}

	userName := fmt.Sprintf("%s %s", userData.FirstName(), userData.LastName())
	emailData := EmailData{
		UserName:       userName,
		ActivationCode: userData.ActivationCode(),
		CurrentYear:    time.Now().Year(),
	}

	emailTemplate, err := template.ParseFiles("core/services/resources/email_template.html")
	if err != nil {
		log.Error("Erro ao carregar o template do email de ativação da conta: ", err)
		return err
	}

	var body bytes.Buffer
	body.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n")
	body.WriteString(fmt.Sprintf("Subject: Bem-vindo(a) ao Você na Câmara!\n"))
	body.WriteString(fmt.Sprintf("To: %s\n", userData.Email()))
	body.WriteString(fmt.Sprintf("From: %s\n\n", smtpUserEmail))

	err = emailTemplate.Execute(&body, emailData)
	if err != nil {
		log.Error("Erro ao executar o template do email de ativação da conta: ", err)
		return err
	}

	smtpAddress := fmt.Sprint(smtpHost, ":", smtpPort)
	err = smtp.SendMail(smtpAddress, auth, smtpUserEmail, []string{userData.Email()}, body.Bytes())
	if err != nil {
		log.Errorf("Erro ao enviar email de ativação da conta do usuário %s: %s", userData.Email(), err.Error())
		return err
	}

	log.Info("Email de ativação da conta enviado para o usuário ", userData.Email())
	return nil
}
