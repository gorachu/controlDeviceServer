package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string `yaml:"env" env-default:"local"`
	StoragePath      string `yaml:"storage_path" env-required:"true"`
	HTTPServer       `yaml:"http_server"`
	PublicAddress    string `yaml:"public_address"`
	GoogleAppScripts string `yaml:"google_app_scripts"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8443"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH_CONTROLDEVICESERVER")
	if configPath == "" {
		log.Fatal("config path is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
