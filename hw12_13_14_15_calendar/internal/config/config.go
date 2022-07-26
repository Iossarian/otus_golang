package config

import (
	"github.com/spf13/viper"
	"net"
	"strings"
	"time"
)

type Config struct {
	Host                     string        `mapstructure:"HOST"`
	HTTPPort                 string        `mapstructure:"HTTP_PORT"`
	GrpcPort                 string        `mapstructure:"GRPC_PORT"`
	StorageSource            string        `mapstructure:"STORAGE_SOURCE"`
	LoggingFile              string        `mapstructure:"LOGGING_FILE"`
	DBUser                   string        `mapstructure:"DB_USER"`
	DBPassword               string        `mapstructure:"DB_PASSWORD"`
	DBTable                  string        `mapstructure:"DB_NAME"`
	LogLevel                 string        `mapstructure:"LOG_LEVEL"`
	EventScanTimeout         time.Duration `mapstructure:"EVENT_SCAN_TIMEOUT"`
	AMQPAddress              string        `mapstructure:"AMQP_ADDRESS"`
	AMQPQueueName            string        `mapstructure:"AMQP_QUEUE_NAME"`
	AMQPHandlersCount        int           `mapstructure:"AMQP_HANDLERS_COUNT"`
	AMQPMaxReconnectionTries int           `mapstructure:"AMQP_MAX_RECONNECT_TRIES"`
	AMQPReconnectTimeout     time.Duration `mapstructure:"AMQP_RECONNECT_TIMEOUT"`
}

func New(configPath string) Config {
	var config Config
	readConfig(configPath)
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func (c Config) GetHTTPAddr() string {
	return net.JoinHostPort(c.Host, c.HTTPPort)
}

func (c Config) GetGRPCAddr() string {
	return net.JoinHostPort(c.Host, c.GrpcPort)
}

func readConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
