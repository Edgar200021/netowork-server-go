package emailclient

import (
	"fmt"
	"net/smtp"
	"net/url"
)

func (s *EmailClient) SendResetPasswordEmail(token, to string) error {
	subject := "Password Recovery"

	resetPasswordLink := fmt.Sprintf(
		"%s%s?token=%s&email=%s", s.appConfig.ClientUrl, s.appConfig.ResetPasswordPath,
		token, to,
	)

	if _, err := url.Parse(resetPasswordLink); err != nil {
		return err
	}

	body := fmt.Sprintf(
		`
	<html>
		<body>
			<h2>Hello!</h2>
			<p>We received a request to reset your password.</p>
			<p>You can reset your password by clicking the link below:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>If you did not request this, please ignore this email.</p>
			<p>Thank you!</p>
		</body>
	</html>
`, resetPasswordLink,
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
		return fmt.Errorf("failed to send reset password email: %w", err)
	}

	return nil
}
