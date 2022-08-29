package config

import (
	"fmt"
	"net"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http"`

	GRPC struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"grpc"`

	Logger struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logger"`

	DB struct {
		ConnectionAddr string `yaml:"connectionAddr"`
	} `yaml:"db"`

	AMQP struct {
		Addr             string        `yaml:"addr"`
		Name             string        `yaml:"name"`
		HandlersCount    int           `yaml:"handlersCount"`
		EventScanTimeout time.Duration `yaml:"eventScanTimeout"`
	} `yaml:"amqp"`

	StorageSource string `yaml:"storageSource"`
}

func New(configPath string) *Config {
	f, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("fail to open config fail: %w", err))
	}

	config := &Config{}
	if err := yaml.Unmarshal(f, config); err != nil {
		fmt.Println(fmt.Errorf("fail to decode config file: %w", err))
	}

	return config
}

func (c Config) GetHTTPAddr() string {
	return net.JoinHostPort(c.HTTP.Host, c.HTTP.Port)
}

func (c Config) GetGRPCAddr() string {
	return net.JoinHostPort(c.GRPC.Host, c.GRPC.Port)
}
