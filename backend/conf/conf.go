package conf

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Domain string `json:"domain"`
	Addr string `json:"addr"`
	YouTubeApiUrl string `json:"youtube_api_url"`
	YouTubeApiKeys []string `json:"youtube_api_keys"`
	DSN string `json:"dsn"`
}

func ParseConfig(path string) *Config {
	// parse config
	configFile, err := os.Open(path)
    if err != nil {
        log.Fatalf("cannot open %s: %s", path, err.Error())
    }
    defer configFile.Close()

	log.Printf("successfully opened %s", path)

	config := &Config{}
    json.NewDecoder(configFile).Decode(config)

	// paste default value
	if config.YouTubeApiUrl == "" {
		config.YouTubeApiUrl = "https://www.googleapis.com/youtube/v3"
	}

	if len(config.YouTubeApiKeys) == 0 {
		log.Fatal("You forgot to provide YouTube API keys!")
	}

	if config.DSN == "" {
		config.DSN = "server.db"
	}

	return config
}