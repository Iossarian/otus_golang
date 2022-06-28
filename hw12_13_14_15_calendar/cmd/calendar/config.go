package main

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Logger        LoggerConf
	Host          string `mapstructure:"HOST"`
	Port          string `mapstructure:"PORT"`
	StorageSource string `mapstructure:"STORAGE_SOURCE"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

func NewConfig() Config {
	var config Config
	readConfig()
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	config.Logger.Level = viper.GetString("LOG_LEVEL")

	return config
}

func (c Config) GetAddr() string {
	hostPort := []string{c.Host, c.Port}

	return strings.Join(hostPort, ":")
}

func readConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
