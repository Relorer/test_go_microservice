package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Reindexer struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"reindexer"`
	App struct {
		Port            int   `mapstructure:"port"`
		TTL             int64 `mapstructure:"ttl"`
		CleanupInterval int64 `mapstructure:"cleanupinterval"`
	} `mapstructure:"app"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
