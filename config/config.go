package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Version string
	Presets []string
	Rules   []Rule
}

type Rule struct {
	PackageManager string
	PackageName    string
	PackageType    string // direct, indirect, production, development, all
	Labels         []string
}

var (
	config Config
)

func Load() error {
	viper.SetConfigName("snuuze")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".github")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}
	return nil
}

func GetConfig() Config {
	return config
}

func GetRules() []Rule {
	return config.Rules
}
