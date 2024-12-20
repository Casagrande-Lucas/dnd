package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		APP    *APP
		Server *Server
		DB     *DB
		CORS   *CORS
	}

	Server struct {
		Port string
	}

	APP struct {
		ENV string
	}

	DB struct {
		Type       string
		DSN        string
	}

	CORS struct {
		AllowOrigins     []string
		AllowMethods     []string
		AllowHeaders     []string
		ExposeHeaders    []string
		AllowCredentials bool
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("fatal error config file: default %v\n", err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			log.Fatalf("fatal error unmarshal: default %v\n", err)
		}
	})

	return configInstance
}
