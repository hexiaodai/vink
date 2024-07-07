package config

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var Instance *Configuration

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

func ParseConfigFromFile(configPath string) error {
	cfg := &Configuration{}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(cfg, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
	Instance = cfg
	if err != nil {
		return err
	}
	return validate()
}

func validate() error {
	return nil
}
