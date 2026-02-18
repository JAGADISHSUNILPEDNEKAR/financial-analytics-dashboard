package config

import (
	"time"
)

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	Auth   AuthConfig
}

type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type AuthConfig struct {
	ServiceURL string
}

func Load() (*Config, error) {
	// In a real app, we would load from env vars
	return &Config{
		Server: ServerConfig{
			Address:      ":8080",
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		},
		Redis: RedisConfig{
			Address: "localhost:6379",
		},
		Auth: AuthConfig{
			ServiceURL: "http://localhost:8082",
		},
	}, nil
}
