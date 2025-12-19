
package config

import (
	"github.com/spf13/viper"
	"path/filepath"
)

type KeycloakConfig struct {
	URL           string `mapstructure:"KEYCLOAK_URL"`
	Realm         string `mapstructure:"KEYCLOAK_REALM"`
	ClientID      string `mapstructure:"KEYCLOAK_CLIENT_ID"`
	ClientSecret  string `mapstructure:"KEYCLOAK_CLIENT_SECRET"`
	AdminUsername string `mapstructure:"KEYCLOAK_ADMIN_USERNAME"`
	AdminPassword string `mapstructure:"KEYCLOAK_ADMIN_PASSWORD"`
}

type Config struct {
	Port     int           `mapstructure:"PORT"`
	Keycloak KeycloakConfig `mapstructure:",squash"`
}

func LoadConfig() (*Config ,error) {
// Get the absolute path to the .env file
	envPath := filepath.Join("..", ".env") // Goes up one level from internal/configs

	viper.SetConfigFile(envPath)
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