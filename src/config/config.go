package config

import (
	"sync"

	"github.com/spf13/viper"
	"io.hyperd/inspectmx/logger"
)

var (
	once           sync.Once
	configInstance IMXConfig
)

// OrangeIdeasConfig exported
type IMXConfig struct {
	Paths       PathsConfig
	Server      ServerConfig
	Environment EnvironmentConfig
	Redis       RedisConfig
	Email       EmailConfig
}

// PathsConfig configuration exported
type PathsConfig struct {
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

// ServerConfig configuration exported
type ServerConfig struct {
	HTTPPort int `mapstructure:"http_port"`
	TLSPort  int `mapstructure:"tls_port"`
}

// EnvironmentConfig exported
type EnvironmentConfig struct {
	Development bool   `mapstructure:"development"`
	LogLevel    string `mapstructure:"log_level"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type EmailConfig struct {
	AllowedProviders []string `mapstructure:"allowed_providers"`
}

func Instance() IMXConfig {
	once.Do(func() {
		// var (
		// 	env string
		// 	ok  bool
		// )
		// if env, ok = os.LookupEnv("ENV"); !ok {
		// 	logger.Info("configuration environment not set, falling back to production", nil)
		// 	env = "prod"
		// } else {
		// 	env = strings.ToLower(env)
		// }

		logger.Info("loading configuration", nil)
		// setup viper and populate the instance config obj
		viper.SetConfigName(".config")
		viper.SetConfigType("yml")
		// viper.AddConfigPath(".config/" + env + "/")

		if err := viper.ReadInConfig(); err != nil { // Find and read the config file
			exit(err)
		}

		err := viper.Unmarshal(&configInstance)
		if err != nil {
			exit(err)
		}
	})

	return configInstance
}

func exit(err error) {
	logger.Fatal("can't take it anymore", logger.WithFields{"error": err})
}
