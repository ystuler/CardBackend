package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	IP   string
	Port string
}

type DatabaseConfig struct {
	DSN string
}

func NewConfig() *Config {
	viper.SetDefault("server.ip", "localhost")
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("database.dsn", "user=admin password=root dbname=memoryCards sslmode=disable")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config
}
