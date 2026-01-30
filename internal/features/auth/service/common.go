package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/constants"
	"github.com/Desalutar20/lingostruct-server-go/pkg/random"
)

func (s *Service) generateSession(ctx context.Context, userId string) (string, error) {
	sessionId, err := random.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	if err := s.redis.SetEx(ctx, s.generateSessionKey(sessionId), userId, time.Minute*time.Duration(s.config.SessionTTLMinutes)).Err(); err != nil {
		return "", err
	}

	return sessionId, nil

}

func (s *Service) generateVerificationKey(token string) string {
	return fmt.Sprintf("%s%s", constants.RedisVerificationPrefix, token)
}

func (s *Service) generateSessionKey(sessionId string) string {
	return fmt.Sprintf("%s%s", constants.RedisSessionPrefix, sessionId)
}

func (s *Service) generateResetPasswordkey(token string) string {
	return fmt.Sprintf("%s%s", constants.RedisResetPasswordPrefix, token)
}

func (s *Service) sendVerificationEmail(token, to string) error {
	subject := "Verify your email address"

	verifyURL := fmt.Sprintf("%s%s?token=%s&email=%s", s.config.ClientUrl, s.config.AccountVerificationPath, token, url.QueryEscape(to))

	text := fmt.Sprintf(
		"Welcome to Lingostruct!\r\n\r\n"+
			"Thanks for signing up. To verify your email address, open the link below:\r\n\r\n"+
			"%s\r\n\r\n"+
			"If you didn’t create an account, you can safely ignore this email.\r\n\r\n"+
			"— Lingostruct Team",
		verifyURL)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 20px;">
    <div style="max-width: 600px; margin: auto; background: #ffffff; padding: 24px; border-radius: 8px;">
      <h2>Welcome to Lingostruct 👋</h2>

      <p>
        Thanks for signing up! To complete your registration,
        please verify your email address.
      </p>

      <p style="text-align: center; margin: 32px 0;">
        <a href="%s" _target="blank"
           style="
             background-color: #4f46e5;
             color: #ffffff;
             padding: 12px 24px;
             text-decoration: none;
             border-radius: 6px;
             font-weight: bold;
           ">
          Verify email
        </a>
      </p>

      <p>If the button doesn’t work, open this link:</p>

      <p style="word-break: break-all;">
        <a href="%s">%s</a>
      </p>

      <hr style="margin: 32px 0;" />

      <p style="font-size: 12px; color: #666;">
        If you didn’t create an account, you can safely ignore this email.
      </p>
    </div>
  </body>
</html>
`, verifyURL, verifyURL, verifyURL)

	return s.emailSender.Send(subject, text, html, []string{to})
}

func (s *Service) sendResetPasswordEmail(token, to string) error {
	subject := "Reset your password"

	resetURL := fmt.Sprintf("%s%s?token=%s&email=%s", s.config.ClientUrl, s.config.ResetPasswordPath, token, url.QueryEscape(to))

	text := fmt.Sprintf(
		"Hello!\r\n\r\n"+
			"We received a request to reset your password. Click the link below to reset it:\r\n\r\n"+
			"%s\r\n\r\n"+
			"If you didn’t request a password reset, you can safely ignore this email.\r\n\r\n"+
			"— Lingostruct Team",
		resetURL)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 20px;">
    <div style="max-width: 600px; margin: auto; background: #ffffff; padding: 24px; border-radius: 8px;">
      <h2>Password Reset Request 🔑</h2>

      <p>
        We received a request to reset your password. Click the button below to proceed.
      </p>

      <p style="text-align: center; margin: 32px 0;">
        <a href="%s" target="_blank"
           style="
             background-color: #4f46e5;
             color: #ffffff;
             padding: 12px 24px;
             text-decoration: none;
             border-radius: 6px;
             font-weight: bold;
           ">
          Reset Password
        </a>
      </p>

      <p>If the button doesn’t work, open this link:</p>

      <p style="word-break: break-all;">
        <a href="%s">%s</a>
      </p>

      <hr style="margin: 32px 0;" />

      <p style="font-size: 12px; color: #666;">
        If you didn’t request a password reset, you can safely ignore this email.
      </p>
    </div>
  </body>
</html>
`, resetURL, resetURL, resetURL)

	return s.emailSender.Send(subject, text, html, []string{to})
}

func (h *Service) generateSessionCookieOptions(sessionId string) *http.Cookie {
	return &http.Cookie{
		Name:     h.config.SessionCookieName,
		Value:    sessionId,
		HttpOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().UTC().Add(time.Minute * time.Duration(h.config.SessionTTLMinutes))}
}
