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
	c := viper.New()
	c.SetConfigName("snuuze")
	c.SetConfigType("yaml")
	c.AddConfigPath(".")
	c.AddConfigPath(".github")
	err := c.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// use default config
			setDefaultConfig()
			return nil
		}
		return fmt.Errorf("fatal error config file: %w", err)
	} else {
		err = c.Unmarshal(&config)
		if err != nil {
			return fmt.Errorf("unable to decode into struct, %w", err)
		}
	}
	return nil
}

func setDefaultConfig() {
	config.Version = "1"
	config.Presets = []string{"base"}
	config.Rules = []types.Rule{
		{
			Name:            "all dependencies",
			PackageManagers: types.AllPackageManagers,
		},
	}
}

func GetConfig() types.Config {
	return config
}

func GetRules() []types.Rule {
	return config.Rules
}
