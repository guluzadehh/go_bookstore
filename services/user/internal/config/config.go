package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	GRPCServer GRPCServer `yaml:"grpc_server"`
}

type GRPCServer struct {
	Port int `yaml:"port" env-default:"50000"`
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file `%s` does not exist", cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("can't read config file `%s` and env variables: %s", cfgPath, err)
	}

	return &cfg
}
