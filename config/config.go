package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Badger BadgerConfig `yaml:"badger"`
	Redis  RedisConfig  `yaml:"redis"`
	BM25   BM25Config   `yaml:"bm25"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type BadgerConfig struct {
	Path string `yaml:"path"`
}

type RedisConfig struct {
	Host               string
	Port               int
	Password           string
	DB                 int
	DialConnectTimeout time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxIdle            int
	MaxActive          int
	IdleTimeout        time.Duration
	MaxConnLifetime    time.Duration
}

type BM25Config struct {
	K1 float64 `yaml:"k1"`
	B  float64 `yaml:"b"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %s", err)
		return nil, err
	}

	return &config, nil
}
