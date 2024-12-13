package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	GRPCServer GRPCServer `yaml:"grpc_server"`
	Postgresql Postgresql `yaml:"postgresql"`
	Migrations Migrations `yaml:"migrations"`
}

type GRPCServer struct {
	Port int `yaml:"port" env-default:"50000"`
}

type Postgresql struct {
	Url     string   `yaml:"-" env:"POSTGRES_URL" env-required:"true"`
	Options []string `yaml:"options"`
}

func (db *Postgresql) DSN(options []string) string {
	if options == nil {
		options = db.Options
	}

	opts := strings.Join(options, "&")
	if len(opts) == 0 {
		return db.Url
	}

	return fmt.Sprintf("%s?%s", db.Url, opts)
}

type Migrations struct {
	Path          string `yaml:"path" env-default:"./migrations"`
	TableName     string `yaml:"table_name" env-default:"migrations"`
	PostgresqlUrl string `yaml:"-" env:"MIGRATIONS_POSTGRES_URL"`
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
