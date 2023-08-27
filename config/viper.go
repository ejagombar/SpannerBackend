package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	CLIENT_ID     string `mapstructure:"CLIENT_ID"`
	CLIENT_SECRET string `mapstructure:"CLIENT_SECRET"`
	PORT          string `mapstructure:"PORT"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			CLIENT_ID:     os.Getenv("CLIENT_ID"),
			CLIENT_SECRET: os.Getenv("CLIENT_SECRET"),
			PORT:          os.Getenv("PORT"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.CLIENT_ID == "" {
		err = errors.New("CLIENT_ID is required")
		return
	}

	if config.CLIENT_SECRET == "" {
		err = errors.New("CLIENT_SECRET is required")
		return
	}

	return
}
