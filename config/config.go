package config

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Configuration struct {
	Debug      bool      `json:"debug"`
	APIServer  APIServer `json:"apiserver"`
	KubeConfig string    `json:"kubeConfig"`
}

type APIServer struct {
	Address string `json:"address"`
	GRPC    string `json:"grpc"`
	HTTP    string `json:"http"`
}

func ParseConfigFromFile(configPath string) (*Configuration, error) {
	cfg := &Configuration{}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(cfg, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
	if err != nil {
		return nil, err
	}
	if err := validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validate() error {
	return nil
}
