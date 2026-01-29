package config

type RateLimitConfig struct {
	SignUp         uint `env:"RATE_LIMIT_SIGN_UP,notEmpty" validate:"required,min=3,max=10"`
	SignIn         uint `env:"RATE_LIMIT_SIGN_IN,notEmpty" validate:"required,min=3,max=10"`
	VerifyAccount  uint `env:"RATE_LIMIT_VERIFY_ACCOUNT,notEmpty" validate:"required,min=3,max=10"`
	GetMe          uint `env:"RATE_LIMIT_GET_ME,notEmpty" validate:"required,min=5,max=20"`
	ForgotPassword uint `env:"RATE_LIMIT_FORGOT_PASSWORD,notEmpty" validate:"required,min=3,max=5"`
	ResetPassword  uint `env:"RATE_LIMIT_RESET_PASSWORD,notEmpty" validate:"required,min=3,max=5"`
	Logout         uint `env:"RATE_LIMIT_LOGOUT,notEmpty" validate:"required,min=3,max=10"`
}
