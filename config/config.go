package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Engine    EngineConfig  `yaml:"engine"`
	TCPServer ServerConfig  `yaml:"network"`
	Logging   LoggingConfig `yaml:"logging"`
	WAL       WalConfig     `yaml:"wal"`
}

type EngineConfig struct {
	Type string `yaml:"type" env-default:"in_memory"`
}

type ServerConfig struct {
	Address        string        `yaml:"address" env-default:"127.0.0.1:3223"`
	MaxConnections int           `yaml:"max_connections" env-default:"100"`
	MaxMessageSize string        `yaml:"max_message_size" env-default:"4KB"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"5m"`
}

type LoggingConfig struct {
	Level  string `yaml:"level" env-default:"info"`
	Output string `yaml:"output" env-default:"./log/output.log"`
}

type WalConfig struct {
	TurnOn               bool          `yaml:"turn_on" env-default:"false"`
	FlushingBatchSize    int           `yaml:"flushing_batch_size" env-default:"100"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout" env-default:"10ms"`
	MaxSegmentSize       string        `yaml:"max_segment_size" env-default:"10MB"`
	Directory            string        `yaml:"data_directory" env-default:"./data/wal"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	var cnf Config

	if _, err := os.Stat(configPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Config file not found, using default values")
			if err := cleanenv.ReadEnv(&cnf); err != nil {
				log.Fatalf("Error read default parametrs %s", err.Error())
			}
			return &cnf
		}
		log.Fatalf("Error checking config file: %v", err)
	}

	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	return &cnf
}

func PasreSize(sizeConf string) (int, error) {
	var size string
	var bytes string
	sizeConf = strings.TrimSpace(sizeConf)
	for _, val := range sizeConf {
		if unicode.IsDigit(val) {
			size += string(val)
		} else {
			bytes += string(val)
		}
	}

	if size == "" {
		return 0, errors.New("size empty")
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return 0, err
	}
	bytes = strings.TrimSpace(strings.ToUpper(bytes))
	switch bytes {
	case "", "Б", "B", "BYTE", "BYTES":
		return sizeInt, nil
	case "КБ", "KB", "K":
		return sizeInt * 1024, nil
	case "МБ", "MB", "M":
		return sizeInt * 1024 * 1024, nil
	case "ГБ", "GB", "G":
		return sizeInt * 1024 * 1024 * 1024, nil
	default:
		return 0, errors.New("Unknown format: " + bytes)
	}
}
