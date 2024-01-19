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
	REFRESH_TOKEN string `mapstructure:"REFRESH_TOKEN"`
	ACCESS_TOKEN  string `mapstructure:"ACCESS_TOKEN"`
	TOKEN_TIMEOUT string `mapstructure:"TOKEN_TIMEOUT"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			CLIENT_ID:     os.Getenv("CLIENT_ID"),
			CLIENT_SECRET: os.Getenv("CLIENT_SECRET"),
			PORT:          os.Getenv("PORT"),
			REFRESH_TOKEN: os.Getenv("REFRESH_TOKEN"),
			ACCESS_TOKEN:  os.Getenv("ACCESS_TOKEN"),
			TOKEN_TIMEOUT: os.Getenv("TOKEN_TIMEOUT"),
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
		return config, err
	}

	if config.CLIENT_SECRET == "" {
		err = errors.New("CLIENT_SECRET is required")
		return config, err
	}

	return config, nil
}

func UpdateToken(AccessToken, RefreshToken, TokenTimeout string) {
	viper.Set("ACCESS_TOKEN", AccessToken)
	viper.Set("REFRESH_TOKEN", RefreshToken)
	viper.Set("TOKEN_TIMEOUT", TokenTimeout)
	viper.WriteConfig()
}
