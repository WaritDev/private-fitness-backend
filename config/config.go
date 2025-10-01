package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string `mapstructure:"PORT"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	Schema           string `mapstructure:"SCHEMA"`
	UseDocker        bool   `mapstructure:"USE_DOCKER"`
}

func ProvideConfig() *Config {
	config := &Config{}
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading config file", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln("Unable to decode into struct", err)
	}

	return config
}