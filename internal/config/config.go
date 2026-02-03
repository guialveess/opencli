package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Project string `mapstructure:"project"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(filepath.Join(home, ".config", "opcli"))

	viper.BindEnv("base_url", "OPENPROJECT_BASE_URL")
	viper.BindEnv("project", "OPENPROJECT_PROJECT")
	viper.BindEnv("api_key", "OPENPROJECT_API_KEY")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.BaseURL == "" || cfg.APIKey == "" || cfg.Project == "" {
		return nil, errors.New("base_url, api_key ausente ou project ausente")
	}

	return &cfg, nil
}
