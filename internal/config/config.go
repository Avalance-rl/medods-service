package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer   `yaml:"http_server"`
	Database `yaml:"database"`
	Redis `yaml:"redis"`
	JWTSecretKey string `yaml:"jwt_secret_key"` 
}

type HTTPServer struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type Database struct {
	Address  string `yaml:"address"`	
}

type Redis struct {
	Address string `yaml:"address"`
	Password string `yaml:"password"`
	DB int `yaml:"DB"`
}

var CFG Config

func Init() {
	CFG = mustLoad("configs/main.yml")
}



func mustLoad(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config files does not exist: %s", path)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return cfg
}