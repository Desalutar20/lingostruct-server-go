package config

type RedisConfig struct {
	Host     string `env:"REDIS_HOST,notEmpty" validate:"hostname|ip"`
	Port     uint   `env:"REDIS_PORT,notEmpty" validate:"port"`
	User     string `env:"REDIS_USER,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	Database uint   `env:"REDIS_DATABASE,notEmpty" validate:"min=0,max=15"`
}

// /   "Redis": {
// // 	"Host": "Localhost",
// // 	"Port": 6379,
// // 	"User": "",
// // 	"Password": "",
// // 	"Database": 0
// //   },
