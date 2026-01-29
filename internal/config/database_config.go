package config

type DatabseConfig struct {
	Host     string `env:"DATABASE_HOST,notEmpty" validate:"required,hostname|ip"`
	Port     uint   `env:"DATABASE_PORT,notEmpty" validate:"required,port"`
	Name     string `env:"DATABASE_NAME,notEmpty" validate:"required"`
	User     string `env:"DATABASE_USER,notEmpty" validate:"required"`
	Password string `env:"DATABASE_PASSWORD,notEmpty" validate:"required"`
	Ssl      bool   `env:"DATABASE_SSL,notEmpty"`
}
