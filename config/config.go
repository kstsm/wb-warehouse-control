package config

import (
	"os"
	"time"

	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Postgres Postgres
	JWT      JWT
}

type Server struct {
	Host string
	Port int
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	Ssl      string
}

type JWT struct {
	Secret string
	TTL    time.Duration
	Issuer string
}

func GetConfig() Config {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		slog.Fatal("Failed to read .env file", "error", err)
		os.Exit(1)
	}

	return Config{
		Server: Server{
			Host: viper.GetString("SRV_HOST"),
			Port: viper.GetInt("SRV_PORT"),
		},
		Postgres: Postgres{
			Username: viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			DBName:   viper.GetString("POSTGRES_DB"),
			Ssl:      viper.GetString("POSTGRES_SSL"),
		},
		JWT: JWT{
			Secret: viper.GetString("JWT_SECRET"),
			TTL:    viper.GetDuration("JWT_TTL"),
			Issuer: viper.GetString("JWT_ISSUER"),
		},
	}
}
