package config

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var Instance = &Config{}

type Config struct {
	Debug     bool      `json:"debug"`
	APIServer APIServer `json:"apiserver"`
}

type APIServer struct {
	GRPC string `json:"grpc"`
	HTTP string `json:"http"`
}

func ParseConfigFromFile(configPath string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(Instance, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
	if err != nil {
		return err
	}
	if err := validate(); err != nil {
		return err
	}
	return nil
}

func validate() error {
	return nil
}
