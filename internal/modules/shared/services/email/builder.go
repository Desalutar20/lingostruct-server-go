package email

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

var (
	accountVerificationSubject = mustNonEmptyString("Account verification")
	passwordChangedSubject     = mustNonEmptyString("Your password has been changed")
	resetPasswordSubject       = mustNonEmptyString("Reset your password")
)

func mustNonEmptyString(value string) valueobjects.NonEmptyString {
	result, err := valueobjects.NewNonEmptyString(value, "message", "Email message", 0, 0, nil)
	if err != nil {
		panic(err)
	}

	return result
}

func buildURL(
	clientURL string,
	path string,
	email valueobjects.Email,
	token valueobjects.NonEmptyString,
) (string, error) {
	u, err := url.Parse(clientURL)
	if err != nil {
		return "", domain.InternalError("Failed to parse client URL")
	}

	u.Path = path

	query := u.Query()
	query.Set("token", token.Value())
	query.Set("email", email.Value())

	u.RawQuery = query.Encode()

	return u.String(), nil
}

func BuildAccountVerificationEmail(
	config *config.ApplicationConfig,
	to valueobjects.Email,
	token valueobjects.NonEmptyString,
) (EmailMessage, error) {
	verificationURL, err := buildURL(
		config.ClientUrl,
		config.AccountVerificationPath,
		to,
		token,
	)
	if err != nil {
		return EmailMessage{}, err
	}

	text := mustNonEmptyString(fmt.Sprintf(
		`Hello,

Please verify your account by clicking the link below:

%s

Thank you`,
		verificationURL,
	))

	html := mustNonEmptyString(fmt.Sprintf(
		`<html>
<body style="font-family: Arial, sans-serif; color: #333;">
    <h2 style="color: #2E86C1;">Hello!</h2>
    <p>Please verify your account by clicking the button below:</p>
    <a href="%s"
       target="_blank"
       style="display:inline-block;padding:10px 20px;background-color:#2E86C1;color:white;text-decoration:none;border-radius:5px;font-weight:bold;">
       Verify Account
    </a>
    <p style="margin-top:20px;">If the button doesn't work, copy and paste this link into your browser:</p>
    <p>%s</p>
    <p>Thank you</p>
</body>
</html>`,
		html.EscapeString(verificationURL),
		html.EscapeString(verificationURL),
	))

	return EmailMessage{
		Subject:   accountVerificationSubject,
		To:        to,
		PlainText: text,
		HtmlText:  &html,
	}, nil
}

func BuildResetPasswordEmail(
	config *config.ApplicationConfig,
	to valueobjects.Email,
	token valueobjects.NonEmptyString,
) (EmailMessage, error) {
	resetURL, err := buildURL(
		config.ClientUrl,
		config.ResetPasswordPath,
		to,
		token,
	)
	if err != nil {
		return EmailMessage{}, err
	}

	text := mustNonEmptyString(fmt.Sprintf(
		`Hello,

You requested a password reset.

Open the link below to set a new password:

%s

If you did not request this, you can safely ignore this email.

Thanks`,
		resetURL,
	))

	htmlText := mustNonEmptyString(fmt.Sprintf(
		`<html>
<body style="font-family: Arial, sans-serif; color: #333;">
    <h2>Password reset</h2>
    <p>You requested a password reset.</p>
    <a href="%s"
       target="_blank"
       style="display:inline-block;padding:12px 24px;background-color:#2E86C1;color:#ffffff;text-decoration:none;border-radius:6px;font-weight:bold;">
        Reset password
    </a>
    <p style="margin-top:20px;">
        If the button doesn't work, copy and paste this link into your browser:
    </p>
    <p>%s</p>
    <p>If you did not request a password reset, just ignore this email.</p>
    <p>Thanks</p>
</body>
</html>`,
		html.EscapeString(resetURL),
		html.EscapeString(resetURL),
	))

	return EmailMessage{
		Subject:   resetPasswordSubject,
		To:        to,
		PlainText: text,
		HtmlText:  &htmlText,
	}, nil
}

func BuildPasswordChangedEmail(
	to valueobjects.Email,
) EmailMessage {
	text := mustNonEmptyString(
		`Hello,

Your password has been successfully changed.

If you did not perform this action, please contact support immediately.

Thanks`,
	)

	htmlText := mustNonEmptyString(
		`<html>
<body style="font-family: Arial, sans-serif; color: #333;">
    <h2>Password changed</h2>
    <p>Your password has been successfully changed.</p>
    <p>If you did not perform this action, contact support immediately.</p>
    <p>Thanks</p>
</body>
</html>`,
	)

	return EmailMessage{
		Subject:   passwordChangedSubject,
		To:        to,
		PlainText: text,
		HtmlText:  &htmlText,
	}
}

func SerializeMessage(message EmailMessage) (valueobjects.NonEmptyString, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return valueobjects.NonEmptyString{}, domain.InternalError(
			"Failed to serialize email message",
		)
	}

	return mustNonEmptyString(string(data)), nil
}
