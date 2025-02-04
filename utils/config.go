package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	SecretKey     string        `mapstructure:"SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string, filename string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
