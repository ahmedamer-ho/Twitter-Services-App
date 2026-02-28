package config

import (
	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	URL string `mapstructure:"MONGO_URI"`
}

type NATSConfig struct {
	URL     string `mapstructure:"NATS_URL"`
	Subject string `mapstructure:"NATS_TWEET_SUBJECT"`
}

type Config struct {
	Port    int           `mapstructure:"PORT"`
	MongoDB MongoDBConfig `mapstructure:",squash"`
	NATS    NATSConfig    `mapstructure:",squash"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Bind environment variables to ensure Unmarshal works even without a config file
	viper.BindEnv("MONGO_URI")
	viper.BindEnv("PORT")
	viper.BindEnv("NATS_URL")
	viper.BindEnv("NATS_TWEET_SUBJECT")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")

	// Ignore error if config file is not found
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Double check fields if unmarshal missed them (squash can be tricky)
	if cfg.MongoDB.URL == "" {
		cfg.MongoDB.URL = viper.GetString("MONGO_URI")
	}
	if cfg.Port == 0 {
		cfg.Port = viper.GetInt("PORT")
	}
	if cfg.NATS.URL == "" {
		cfg.NATS.URL = viper.GetString("NATS_URL")
	}
	if cfg.NATS.Subject == "" {
		cfg.NATS.Subject = viper.GetString("NATS_TWEET_SUBJECT")
	}

	return &cfg, nil
}
