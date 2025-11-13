package config

import "os"

type Config struct {
	DatabaseConfig DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
}

func MustLoad() *Config {
	return &Config{
		DatabaseConfig: DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_DB"),
			SslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
	}
}
