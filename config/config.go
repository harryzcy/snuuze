package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/harryzcy/snuuze/types"
)

var (
	config types.Config
)

// LoadWorkflows loads the workflows configurations
func LoadConfig() error {
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

func GetConfig() types.Config {
	return config
}

func GetRules() []types.Rule {
	return config.Rules
}
