package config

import (
	"os"
)

type Config struct {
	DB  *DBConfig
	Web *WebConfig
}

type DBConfig struct {
	Type        string
	Host        string
	Port        string
	Username    string
	Password    string
	Name        string
	SSLMode     string
	DBBatchSize int
}

type WebConfig struct {
	HttpPorts string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Type:     os.Getenv("DB_TYPE"), //mysql
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
		},
		Web: &WebConfig{
			HttpPorts: os.Getenv("HTTP_PORTS"),
		},
	}
}
