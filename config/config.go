package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	TokenTTL   time.Duration `yaml:"token-ttl" env-default:""`
	GRPCConfig GRPCConfig    `yaml:"grpc"`
	Database   Database      `yaml:"database"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout"`
}

type Database struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Port     int    `yaml:"port" env-default:"5432"`
	Server   string `yaml:"server" env-default:"localhost"`
	Database string `yaml:"database" env-default:"postgres"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config file is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not exist " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "config.yaml", "config file path")
	flag.Parse()

	if res == "" {
		res = "./config/config.yaml"
	}
	return res
}
