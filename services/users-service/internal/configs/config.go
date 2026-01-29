package config

import (
	"github.com/spf13/viper"
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
	Port     int            `mapstructure:"PORT"`
	Keycloak KeycloakConfig `mapstructure:",squash"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Bind all expected environment variables
	viper.BindEnv("PORT")
	viper.BindEnv("KEYCLOAK_URL")
	viper.BindEnv("KEYCLOAK_REALM")
	viper.BindEnv("KEYCLOAK_CLIENT_ID")
	viper.BindEnv("KEYCLOAK_CLIENT_SECRET")
	viper.BindEnv("KEYCLOAK_ADMIN_USERNAME")
	viper.BindEnv("KEYCLOAK_ADMIN_PASSWORD")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	// Ignore error if config file is not found
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Fallback manual assignments for squashed structs if Unmarshal fails to map from Env
	if cfg.Port == 0 {
		cfg.Port = viper.GetInt("PORT")
	}
	if cfg.Keycloak.URL == "" {
		cfg.Keycloak.URL = viper.GetString("KEYCLOAK_URL")
		cfg.Keycloak.Realm = viper.GetString("KEYCLOAK_REALM")
		cfg.Keycloak.ClientID = viper.GetString("KEYCLOAK_CLIENT_ID")
		cfg.Keycloak.ClientSecret = viper.GetString("KEYCLOAK_CLIENT_SECRET")
		cfg.Keycloak.AdminUsername = viper.GetString("KEYCLOAK_ADMIN_USERNAME")
		cfg.Keycloak.AdminPassword = viper.GetString("KEYCLOAK_ADMIN_PASSWORD")
	}

	return &cfg, nil
}
