package email

import (
	"fmt"
	"net/smtp"
)

type Sender interface {
	SendOTP(to, otp, purpose string) error
}

type SMTPSender struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
	fromName  string
}

func NewSMTPSender(host string, port int, username, password, fromEmail, fromName string) *SMTPSender {
	return &SMTPSender{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *SMTPSender) SendOTP(to, otp, purpose string) error {
	subject := s.getSubject(purpose)
	body := s.getBody(otp, purpose)

	message := fmt.Sprintf("From: %s <%s>\r\n", s.fromName, s.fromEmail)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-version: 1.0;\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\";\r\n"
	message += "\r\n"
	message += body

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *SMTPSender) getSubject(purpose string) string {
	switch purpose {
	case "register":
		return "Verify Your Email Address"
	case "reset_password":
		return "Reset Your Password"
	default:
		return "Your One-Time Password"
	}
}

func (s *SMTPSender) getBody(otp, purpose string) string {
	var action string
	switch purpose {
	case "register":
		action = "verify your email address"
	case "reset_password":
		action = "reset your password"
	default:
		action = "complete your request"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>OTP Verification</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; border-radius: 8px; padding: 30px; margin: 20px 0;">
        <h2 style="color: #007bff; margin-top: 0;">Your One-Time Password</h2>
        <p>You requested to %s. Please use the following OTP code:</p>
        <div style="background-color: #fff; border: 2px solid #007bff; border-radius: 6px; padding: 20px; text-align: center; margin: 25px 0;">
            <span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #007bff;">%s</span>
        </div>
        <p style="color: #dc3545; font-weight: bold;">This code will expire in 5 minutes.</p>
        <p style="color: #6c757d; font-size: 14px; margin-top: 30px;">If you didn't request this code, please ignore this email.</p>
    </div>
</body>
</html>
`, action, otp)
}
