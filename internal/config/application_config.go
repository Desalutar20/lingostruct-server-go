package config

type ApplicationConfig struct {
	Port                          uint   `env:"APPLICATION_PORT,notEmpty" validate:"required,port"`
	ClientUrl                     string `env:"APPLICATION_CLIENT_URL,notEmpty" validate:"required,url"`
	AccountVerificationPath       string `env:"APPLICATION_ACCOUNT_VERIFICATION_PATH,notEmpty" validate:"required,startswith=/"`
	ResetPasswordPath             string `env:"APPLICATION_RESET_PASSWORD_PATH,notEmpty" validate:"required,startswith=/"`
	SessionCookieName             string `env:"APPLICATION_SESSION_COOKIE_NAME,notEmpty" validate:"required"`
	CookieSecure                  bool   `env:"APPLICATION_COOKIE_SECURE,notEmpty"`
	AccountVerificationTTLMinutes uint   `env:"APPLICATION_ACCOUNT_VERIFICATION_TTL_MINUTES,notEmpty" validate:"required,min=60,max=1440"`
	SessionTTLMinutes             uint   `env:"APPLICATION_SESSION_TTL_MINUTES,notEmpty" validate:"required,min=1440,max=43200"`
	ResetPasswordTTLMinutes       uint   `env:"APPLICATION_RESET_PASSWORD_TTL_MINUTES,notEmpty" validate:"required,min=5,max=10"`
	LogLevel                      string `env:"APPLICATION_LOG_LEVEL,notEmpty" validate:"oneofci=info error debug warn"`
	PrettyLog                     bool   `env:"APPLICATION_PRETTY_LOG,notEmpty"`
}
