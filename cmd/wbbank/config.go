package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServerConf
}

type HTTPServerConf struct {
	HostPort string `toml:"httpServer.hostPort"`
}

// NewConfig make a config from configFilePath.
func NewConfig() Config {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return config
}
