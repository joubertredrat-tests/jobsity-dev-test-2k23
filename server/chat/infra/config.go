package infra

import (
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

const (
	ENV_FILE = ".env"
)

type Config struct {
	ApiHost       string `env:"API_HOST,required"`
	ApiPort       string `env:"API_PORT,required"`
	MongoHost     string `env:"MONGO_HOST,required"`
	MongoPort     string `env:"MONGO_PORT,required"`
	MongoDatabase string `env:"MONGO_DATABASE,required"`
	MongoUser     string `env:"MONGO_USER,required"`
	MongoPassword string `env:"MONGO_PASSWORD,required"`
}

func NewConfig() (Config, error) {
	if err := loadEnv(); err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func loadEnv() error {
	if _, err := os.Stat(ENV_FILE); os.IsNotExist(err) {
		return nil
	}

	return godotenv.Load(ENV_FILE)
}
