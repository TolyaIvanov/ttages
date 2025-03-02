package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	GRPCPort    string      `yaml:"grpc_port" env-default:"50051"`
	StoragePath string      `yaml:"storage_path" env-default:"./storage"`
	Postgres    Database    `yaml:"database" env-required:"true"`
	Redis       RedisConfig `yaml:"redis" env-required:"true"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"postgres"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	DBName   string `yaml:"dbname" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env-default:"redis:6379"`
	Password string `yaml:"password" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist %s", err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded config: %+v\n", cfg)

	return &cfg
}
