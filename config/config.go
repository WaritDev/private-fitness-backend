package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port       string `mapstructure:"PORT"`
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBParams   string `mapstructure:"DB_PARAMS"`
	UseDocker  bool   `mapstructure:"USE_DOCKER"`
	TZ         string `mapstructure:"TZ"`

	// Slip2Go API Configuration
	Slip2GoSecretKey string `mapstructure:"SLIP2GO_SECRET_KEY"`
	MockSlip2Go      bool   `mapstructure:"MOCK_SLIP2GO"`
}

func ProvideConfig() *Config {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("DB_DRIVER", "mysql")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "3307")
	v.SetDefault("DB_PARAMS", "parseTime=true&loc=Asia%2FBangkok&charset=utf8mb4")
	v.SetDefault("TZ", "Asia/Bangkok")

	if err := v.ReadInConfig(); err != nil {
		log.Println("⚠️  No .env file found, using env vars/defaults:", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalln("Unable to decode into struct", err)
	}

	return cfg
}
