package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	MattermostTeamName string `env:"MM_TEAM" default:"vk"`
	MattermostToken    string `env:"MM_TOKEN"`
	MattermostChannel  string `env:"MM_CHANNEL"`
	MattermostServer   string `env:"MM_SERVER"`
	Database           DatabaseConfig
}

type DatabaseConfig struct {
	Server   string `env:"DATABASE_SERVER"`
	User     string `env:"DATABASE_USER"`
	Password string `env:"DATABASE_PASSWORD"`
}

func LoadConfig() *Config {
	cfgApp := &Config{}
	parseConfig(cfgApp)
	return cfgApp
}

func parseConfig(cfg *Config) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}
