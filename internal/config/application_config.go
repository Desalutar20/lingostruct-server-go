package config

type ApplicationConfig struct {
	ClientUrl                     string `env:"APPLICATION_CLIENT_URL,notEmpty" validate:"required,url"`
	AccountVerificationPath       string `env:"APPLICATION_ACCOUNT_VERIFICATION_PATH,notEmpty" validate:"required,startswith=/"`
	ResetPasswordPath             string `env:"APPLICATION_RESET_PASSWORD_PATH,notEmpty" validate:"required,startswith=/"`
	SessionCookieName             string `env:"APPLICATION_SESSION_COOKIE_NAME,notEmpty" validate:"required"`
	ApplicationCookieSecure       bool   `env:"APPLICATION_COOKIE_SECURE,notEmpty"`
	AccountVerificationTTLMinutes uint   `env:"APPLICATION_ACCOUNT_VERIFICATION_TTL_MINUTES,notEmpty" validate:"required,min=60,max=1440"`
	SessionTTLMinutes             uint   `env:"APPLICATION_SESSION_TTL_MINUTES,notEmpty" validate:"required,min=1440,max=43200"`
	ResetPasswordTTLMinutes       uint   `env:"APPLICATION_RESET_PASSWORD_TTL_MINUTES,notEmpty" validate:"required,min=5,max=10"`
}

// "Application": {
// 	"ClientUrl": "http://localhost:4000",
// 	"AccountVerificationPath": "/auth/account-verification",
// 	"ResetPasswordPath": "/auth/reset-password",
// 	"SessionCookieName": "session",
// 	"AccountVerificationTTLMinutes": 1440,
// 	"SessionTTLMinutes": 43200,
// 	"ResetPasswordTTLMinutes": 10
