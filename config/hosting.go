package config

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unicode"

	"github.com/harryzcy/snuuze/types"
	"github.com/spf13/viper"
)

var (
	CONFIG_FILE = defaultValue(os.Getenv("SNUUZE_CONFIG_FILE"), "config.yaml")

	hostingConfig types.HostingConfig
)

// LoadConfig loads the configuration for the application
func LoadHostingConfig() error {
	c := viper.New()
	c.SetConfigFile(CONFIG_FILE)
	c.SetEnvPrefix("SNUUZE")
	c.AddConfigPath(".")
	c.AddConfigPath(".github")
	err := c.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	err = c.Unmarshal(&hostingConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}

	loadEnvs(c)

	return nil
}

func loadEnvs(c *viper.Viper) {
	v := reflect.ValueOf(&hostingConfig).Elem()
	loadEnv(c, "", v)
}

func loadEnv(c *viper.Viper, parent string, v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)

		name := t.Field(i).Tag.Get("yaml")
		if parent != "" {
			name = parent + "." + name
		}

		if field.Type().Kind() == reflect.Struct {
			loadEnv(c, name, field)
			continue
		}

		envName := "SNUUZE_" + toEnvName(name)
		value := os.Getenv(envName)
		if value == "" {
			continue
		}

		if field.IsValid() && field.CanSet() {
			fmt.Println(76, field.Kind())
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Bool:
				isTrue, err := strconv.ParseBool(value)
				if err != nil {
					panic(err)
				}
				field.SetBool(isTrue)
			case reflect.Int, reflect.Int64:
				if value != "" {
					intValue, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						panic(err)
					}
					if !field.OverflowInt(intValue) {
						field.SetInt(intValue)
					}
				}
			default:
				panic("unsupported type")
			}
		}
	}
}

func toEnvName(name string) string {
	if len(name) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteRune(unicode.ToUpper(rune(name[0])))

	previousLower := name[0] >= 'a' && name[0] <= 'z'
	for _, r := range name[1:] {
		if r == '.' {
			buf.WriteRune('_')
			previousLower = false
			continue
		}

		if r >= 'A' && r <= 'Z' {
			if previousLower {
				buf.WriteRune('_')
			}
		}

		buf.WriteRune(unicode.ToUpper(r))
		previousLower = r >= 'a' && r <= 'z'
	}
	return buf.String()
}

func GetHostingConfig() types.HostingConfig {
	return hostingConfig
}

// TempDir returns the temporary directory to store the data
func TempDir() string {
	return defaultValue(hostingConfig.Data.TempDir, os.TempDir())
}

func GitHubToken() string {
	return defaultValue(hostingConfig.GitHub.Token, os.Getenv("GITHUB_TOKEN"))
}

func GoPath() string {
	return os.Getenv("GOPATH")
}

func GoProxy() string {
	return defaultValue(os.Getenv("GOPROXY"), "https://proxy.golang.org,direct")
}

func defaultValue[T comparable](value, defaultValue T) T {
	if value == *new(T) {
		return defaultValue
	}

	return value
}
