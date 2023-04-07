package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	CONFIG_FILE = os.Getenv("SNUUZE_CONFIG_FILE")

	hostingConfig viper.Viper
)

func init() {
	if CONFIG_FILE != "" {
		CONFIG_FILE = "config.yaml"
	}
}

// LoadConfig loads the configuration for the application
func LoadHostingConfig() error {
	viper.SetConfigFile(CONFIG_FILE)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	err = viper.Unmarshal(&hostingConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}
	return nil
}
