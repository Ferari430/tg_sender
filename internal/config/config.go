package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotConfig *BotConfig
}

type BotConfig struct {
	Token string
}

func InitConfig() (*Config, error) {
	err := godotenv.Load(`B:\programmin-20260114T065921Z-1-001\programmin\tg_sender\.env`)
	if err != nil {
		return nil, err
	}

	return &Config{BotConfig: &BotConfig{Token: os.Getenv("TOKEN")}}, nil
}
