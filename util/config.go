package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver       string        `mapstructure:"DB_DRIVER"`
	DBSource       string        `mapstructure:"DB_SOURCE"`
	ServerAddress  string        `mapstructure:"SERVER_ADDRESS"`
	SecretKey      string        `mapstructure:"SECRET_KEY"`
	AccessDuration time.Duration `mapstructure:"ACCESS_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
