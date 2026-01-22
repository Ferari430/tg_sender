package config

import (
	"os"
	"runtime"
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
	var (
		envPath      string
		DownloadPath string
	)

	System := runtime.GOOS
	switch System {
	case "linux":
		envPath = `/home/user/programmin/tg_sender/.env`
		DownloadPath = `/home/user/programmin/tg_sender/data`

	case "windows":
		//todo
		envPath = `B:\programmin-20260114T065921Z-1-001\programmin\tg_sender\.env`
		DownloadPath = `B:\programmin-20260114T065921Z-1-001\programmin\tg_sender\data`
		//todo: set for windows
	}

	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	//todo: delete hardcode
	t := os.Getenv("TIME")
	a, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return nil, err
	}

	dur := time.Duration(a) * time.Second

	return &Config{BotConfig: BotConfig{Token: os.Getenv("TOKEN")}, DownloaderConfig: &DownloaderConfig{RootDir: DownloadPath},

		TickerConfig: TickerConfig{
			TickTime: dur,
		},
	}, nil
}
