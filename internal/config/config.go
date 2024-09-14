package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"strconv"
)

type Config struct {
	Env string `env:"ENV" envDefault:"dev"`
	DB  DBConfig
	App AppConfig
}

type DBConfig struct {
	URL    string `env:"DATABASE_URL" env-required:"true"`
	ApiKey string `env:"API_KEY" env-required:"true"`
}

type AppConfig struct {
	Host string `env:"APP_HOST" envDefault:"localhost"`
	Port int    `env:"APP_PORT" envDefault:"8080"`
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
		log.Printf("No loading .env file: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("config is empty: " + err.Error())
	}

	return &cfg

}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
