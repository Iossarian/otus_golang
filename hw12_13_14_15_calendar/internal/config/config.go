package config

import (
	"github.com/spf13/viper"
	"net"
	"strings"
)

type Config struct {
	Host          string `mapstructure:"HOST"`
	Port          string `mapstructure:"PORT"`
	StorageSource string `mapstructure:"STORAGE_SOURCE"`
	LoggingFile   string `mapstructure:"LOGGING_FILE"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBTable       string `mapstructure:"DB_NAME"`
}

func NewConfig(configPath string) Config {
	var config Config
	readConfig(configPath)
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func (c Config) GetAddr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func readConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
