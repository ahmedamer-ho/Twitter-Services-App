
package config

import (
	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	URL           string `mapstructure:"MONGO_URI"`
}

type Config struct {
	Port     int           `mapstructure:"PORT"`
	MongoDB MongoDBConfig `mapstructure:",squash"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("../..") // project root
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
