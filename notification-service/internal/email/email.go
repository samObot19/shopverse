package email

import (
    "fmt"
    "net/smtp"
)

type EmailService struct {
    SMTPHost     string
    SMTPPort     string
    Username     string
    Password     string
    FromEmail    string
}

func NewEmailService(smtpHost, smtpPort, username, password, fromEmail string) *EmailService {
    return &EmailService{
        SMTPHost:  smtpHost,
        SMTPPort:  smtpPort,
        Username:  username,
        Password:  password,
        FromEmail: fromEmail,
    }
}

func (e *EmailService) SendEmail(to string, subject string, body string) error {
    auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPHost)

    message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

    addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
    return smtp.SendMail(addr, auth, e.FromEmail, []string{to}, message)
}