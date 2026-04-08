package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Engine    EngineConfig  `yaml:"engine"`
	TCPServer ServerConfig  `yaml:"network"`
	Logging   LoggingConfig `yaml:"logging"` 
}

type EngineConfig struct {
	Type string `yaml:"type" env-default:"in_memory"`
}

type ServerConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"` 
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"` 
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	if _, err := os.Stat(configPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Config file not found: %s", configPath)
		}
		log.Fatalf("Error checking config file: %v", err)
	}

	var cnf Config
	
	
	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("Error reading config: %v", err) 
	}

	return &cnf
}