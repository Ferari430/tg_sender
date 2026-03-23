package config

import (
	"flag"
	"log"
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
	KafkaConfig      KafkaConfig
}

type KafkaConfig struct {
	BrokersAddr     string
	ConsumerGroupID string
	Topic           string
}

type TickerConfig struct {
	TickTime time.Duration
}

type BotConfig struct {
	Token      string
	WithClient bool
	UseProxy   bool   // Добавьте
	ProxyURL   string // Добавьте, например "socks5://127.0.0.1:1080"
}

type DownloaderConfig struct {
	RootDir string
}

func InitConfig() (*Config, error) {
	log.Println("initialaize config...")
	var (
		envPath      string
		DownloadPath string
		withC        bool
	)

	System := runtime.GOOS
	switch System {
	case "linux":
		envPath = `/home/user/programmin/obsidian_Project/prog/tg_sender/.env`
		DownloadPath = `/home/user/data`

	case "windows":
		envPath = `B:\programmin-20260114T065921Z-1-001\programmin\obsidian_Project\prog\tg_sender\.env`
		DownloadPath = `B:\programmin-20260114T065921Z-1-001\programmin\obsidian_Project\prog\tg_sender\data\new`
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

	flag.BoolVar(&withC, "withClient", true, "tgBot with custom client")
	flag.Parse()

	wc, ok := os.LookupEnv("WITH_CLIENT")
	if ok {
		boolVal, _ := strconv.ParseBool(wc)
		withC = boolVal
	}

	return &Config{BotConfig: BotConfig{Token: os.Getenv("TOKEN"),
		WithClient: withC,
		UseProxy:   true,
		ProxyURL:   "socks5://127.0.0.1:1080",
	}, DownloaderConfig: &DownloaderConfig{RootDir: DownloadPath},

		TickerConfig: TickerConfig{
			TickTime: dur,
		},

		KafkaConfig: KafkaConfig{
			BrokersAddr:     os.Getenv("KAFKA_PORT"),
			ConsumerGroupID: os.Getenv("CONSUMER_GROUP_ID"),
			Topic:           os.Getenv("TOPIC"),
		},
	}, nil
}
