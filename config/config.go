package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TelegramToken                 string
	MaximumAudioDurationInSeconds int
}

var cfg *Config

// Load initializes the configuration
func Load() {
	cfg = &Config{
		TelegramToken: strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
	}
	if len(cfg.TelegramToken) == 0 {
		panic("environment variable 'TELEGRAM_BOT_TOKEN' must not be empty")
	}
	maximumAudioDurationInSeconds, err := strconv.Atoi(strings.TrimSpace(os.Getenv("MAXIMUM_AUDIO_DURATION_IN_SECONDS")))
	if err != nil {
		cfg.MaximumAudioDurationInSeconds = 480
	} else {
		cfg.MaximumAudioDurationInSeconds = maximumAudioDurationInSeconds
	}
}

// Get returns the configuration
func Get() *Config {
	return cfg
}
