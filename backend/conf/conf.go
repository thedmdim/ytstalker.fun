package conf

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Addr           string   `json:"addr"`
	YouTubeApiUrl  string   `json:"youtube_api_url"`
	YouTubeApiKeys []string `json:"youtube_api_keys"`
	DSN            string   `json:"dsn"`
}

func ParseConfig(path string) *Config {
	// parse config

	config := &Config{}

	config.Addr = os.Getenv("ADDR")
	
	config.DSN = os.Getenv("DSN")
	if config.DSN == "" {
		config.DSN = "server.db"
	}

	config.YouTubeApiKeys = strings.Split(os.Getenv("YT_API_KEYS"), ",")
	if len(config.YouTubeApiKeys) == 0 {
		log.Fatal("You forgot to provide YouTube API keys!")
	}

	return config
}