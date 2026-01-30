package config

import (
	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	URL string `mapstructure:"MONGO_URI"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"KAFKA_BROKERS"`
	Topic   string   `mapstructure:"KAFKA_TWEET_TOPIC"`
}

type Config struct {
	Port    int           `mapstructure:"PORT"`
	MongoDB MongoDBConfig `mapstructure:",squash"`
	Kafka   KafkaConfig   `mapstructure:",squash"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Bind environment variables to ensure Unmarshal works even without a config file
	viper.BindEnv("MONGO_URI")
	viper.BindEnv("PORT")
	viper.BindEnv("KAFKA_BROKERS")
	viper.BindEnv("KAFKA_TWEET_TOPIC")

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
	if len(cfg.Kafka.Brokers) == 0 {
		cfg.Kafka.Brokers = viper.GetStringSlice("KAFKA_BROKERS")
	}
	if cfg.Kafka.Topic == "" {
		cfg.Kafka.Topic = viper.GetString("KAFKA_TWEET_TOPIC")
	}

	return &cfg, nil
}
