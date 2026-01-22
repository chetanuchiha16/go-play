package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
}

func Load() Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil{
		log.Println("no env file found",)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Unable to decode into the struct")
	}
	return cfg
}