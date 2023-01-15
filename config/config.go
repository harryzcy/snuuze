package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName("sailor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	return nil
}
