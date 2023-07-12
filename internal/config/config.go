package config

import (
	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HttpPort     string `env:"HTTP_PORT" envDefault:"3333"`
	WOLBroadcast string `env:"WOL_BROADCAST,notEmpty,required"`
	WOLMac       string `env:"WOL_MAC,notEmpty,required"`
	TVEndpoint   string `env:"TV_ENDPOINT,notEmpty,required"`
	TVAuth       string `env:"TV_AUTH,unset,notEmpty,required"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
