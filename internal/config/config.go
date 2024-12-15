package config

import (
	"flag"
	"net"
	"os"
	"strconv"

	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Env  string `env:"ENV" envDefault:"dev"`
	DB   DBConfig
	App  AppConfig
	S3   S3Config
	SMTP SMTPConfig
}

type DBConfig struct {
	URL    string `env:"DATABASE_URL" env-required:"true"`
	ApiKey string `env:"API_KEY" env-required:"true"`
}

type AppConfig struct {
	Host         string `env:"APP_HOST" envDefault:"localhost"`
	Port         int    `env:"APP_PORT" envDefault:"8080"`
	AppUrl       string `env:"APP_URL" envDefault:"http://localhost:8080" env-required:"true"`
	JwtSecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
	JwtExpiresIn string `env:"JWT_EXPIRES_IN" env-default:"24h"`
	FromEmail    string `env:"FROM_EMAIL" env-required:"true"`
	GrpcServer   string `env:"GRPC_SERVER" env-required:"true" envDefault:"localhost:50051"`
}

type S3Config struct {
	Region    string `env:"S3_REGION" env-required:"true"`
	Endpoint  string `env:"S3_ENDPOINT" env-required:"true"`
	Bucket    string `env:"S3_BUCKET" env-required:"true"`
	AccessKey string `env:"S3_ACCESS_KEY_ID" env-required:"true"`
	SecretKey string `env:"S3_SECRET_ACCESS_KEY" env-required:"true"`
	Domain    string `env:"S3_DOMAIN" env-required:"true"`
}

type SMTPConfig struct {
	Host     string `env:"SMTP_HOST" env-required:"true"`
	Port     int    `env:"SMTP_PORT" env-required:"true"`
	Username string `env:"SMTP_USERNAME" env-required:"true"`
	Password string `env:"SMTP_PASSWORD" env-required:"true"`
}

func (config *DBConfig) DSN() string {
	return config.URL
}

func (c *AppConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func MustLoad() *Config {
	var cfg Config
	var err error

	configPath := fetchConfigPath()
	if configPath != "" {
		err = godotenv.Load(configPath)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		logger.Info("Ошибка при загрузке конфигурации: %v", zap.Error(err))
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Конфигурация некорректна:" + err.Error())
	}

	return &cfg

}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "Путь к файлу конфигурации")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
