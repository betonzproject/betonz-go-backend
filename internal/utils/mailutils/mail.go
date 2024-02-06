package mailutils

import (
	"net/smtp"
	"os"
)

func SendMail(to string, message string, subject string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	senderName := os.Getenv("EMAIL_SENDER_NAME")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	password := os.Getenv("SMTP_PASSWORD")

	if os.Getenv("ENVIRONMENT") == "development" && os.Getenv("SEND_TEST_EMAILS_TO") != "" {
		to = os.Getenv("SEND_TEST_EMAILS_TO")
	}

	msg :=
		"From: " + senderName + " <" + from + ">\n" +
			"To: " + to + "\n" +
			"Subject: " + subject + "\n" +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
			message

	auth := smtp.PlainAuth("", from, password, host)
	return smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
}
