package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

type Config struct {
	HTTP          HTTP    `toml:"http"`
	MongoDB       MongoDB `toml:"mongoDB"`
	AuthServerURL string  `toml:"authServerURL"`
}

type HTTP struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

type MongoDB struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

func Parse(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := toml.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg, nil
}

func Generate() (string, error) {
	config := Config{}
	configStr, err := toml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(configStr), nil
}
