package config

import (
	"github.com/chetanuchiha16/go-play/internal/logger"

	"github.com/spf13/viper"
)

type Config struct {
	ENV          string `mapstructure:"ENV"`
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	JWT_SECRET   string `mapstructure:"JWT_SECRET"`
}

func Load() Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetDefault("ENV", "development")
	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Info().Msg("no env file found")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Log.Fatal().Msg("Unable to decode into the struct")
	}
	return cfg
}
