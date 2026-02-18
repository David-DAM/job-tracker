package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port       int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	AppName    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, using system env")
		return nil, err
	}

	cfg := &Config{
		Port:       viper.GetInt("PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		AppName:    viper.GetString("OTEL_SERVICE_NAME"),
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	return cfg, nil
}
