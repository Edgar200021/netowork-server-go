package emailclient

import (
	"fmt"
	"net/smtp"
	"net/url"
)

func (s *EmailClient) SendVerificationEmail(token, to string) error {
	subject := "Account Verification"

	verificationLink := fmt.Sprintf(
		"%s%s?token=%s&email=%s", s.appConfig.ClientUrl, s.appConfig.AccountVerificationPath,
		token, to,
	)

	if _, err := url.Parse(verificationLink); err != nil {
		return err
	}

	body := fmt.Sprintf(
		`
		<html>
			<body>
				<h2>Hello!</h2>
				<p>Please verify your account by clicking the link below:</p>
				<p><a href="%s">Verify Account</a></p>
				<p>Thank you!</p>
			</body>
		</html>
	`, verificationLink,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n%s\r\n",
		s.from, to, subject, body,
	))

	err := smtp.SendMail(s.addr, s.auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}
