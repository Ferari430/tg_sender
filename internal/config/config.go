package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BotConfig        BotConfig
	DownloaderConfig *DownloaderConfig
	TickerConfig     TickerConfig
}

type TickerConfig struct {
	TickTime time.Duration
}

type BotConfig struct {
	Token string
}

type DownloaderConfig struct {
	RootDir string
}

func InitConfig() (*Config, error) {
	err := godotenv.Load(`B:\programmin-20260114T065921Z-1-001\programmin\tg_sender\.env`)
	if err != nil {
		return nil, err
	}

	//todo: delete hardcode
	path := `B:\programmin-20260114T065921Z-1-001\programmin\tg_sender\data`

	t := os.Getenv("TIME")
	a, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return nil, err
	}
	dur := time.Duration(a) * time.Second

	return &Config{BotConfig: BotConfig{Token: os.Getenv("TOKEN")}, DownloaderConfig: &DownloaderConfig{RootDir: path},

		TickerConfig: TickerConfig{
			TickTime: dur,
		},
	}, nil
}
