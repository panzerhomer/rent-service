package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database `yaml:"postgres"`
	App      `yaml:"app"`
}

type App struct {
	Env    string `yaml:"env"`
	Secret string `yaml:"secret" env:"SECRET"`
	Port   int    `yaml:"port" env:"APP_PORT"`
}

type Database struct {
	Port       int    `yaml:"port" env:"POSTGRES_PORT"`
	Host       string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	User       string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password   string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"password"`
	PostgresDB string `yaml:"postgres-db" env:"POSTGRES_DB" env-default:"postgres"`
}

func ReadConfig() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		return nil, fmt.Errorf("read yml error: %v", err.Error())
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env error: %v", err.Error())
	}

	return &cfg, nil
}
