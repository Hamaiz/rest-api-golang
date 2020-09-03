package email

import (
	"net/smtp"
	"os"
)

func SignUpEmail(e string, n string, url string) error {
	// email credentials
	email := os.Getenv("GM_EMAIL")
	pass := os.Getenv("GM_PASS")

	// sending email to
	to := []string{e}

	// smtp configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// email configuration
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Files Signup\n"
	msg := SignUpEmailMsg(n, url)

	// message configuration
	message := []byte(subject + mime + msg)

	// authenticating user
	auth := smtp.PlainAuth("", email, pass, smtpHost)

	// sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, to, message)

	if err != nil {
		return err
	}

	return nil
}
