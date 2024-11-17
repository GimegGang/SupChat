package Config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Address     string        `yaml:"address" required:"true" default:":8080"`
	Timeout     time.Duration `yaml:"timeout" required:"true" default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" required:"true" default:"60s"`
	StoragePath string        `yaml:"storage_path" required:"true" default:"storage/storage.db"`
}

func MustLoad(path string) *Config {
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("config file not found: %s", path)
		return nil
	}
	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		log.Fatalf("config file error: %s", path)
	}
	return &config
}
