package config

type SmtpConfig struct {
	Host     string `env:"SMTP_HOST,notEmpty" validate:"hostname|ip"`
	Port     uint   `env:"SMTP_PORT,notEmpty" validate:"port"`
	User     string `env:"SMTP_USER,notEmpty" validate:"required,email"`
	Password string `env:"SMTP_PASSWORD,notEmpty" validate:"required"`
	From     string `env:"SMTP_FROM,notEmpty" validate:"required"`
}
