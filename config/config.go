package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load() error {
	viper.SetConfigName("latte")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".github")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	return nil
}
