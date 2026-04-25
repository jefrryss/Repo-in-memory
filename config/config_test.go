package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_DefaultsOnly(t *testing.T) {
	t.Setenv("CONFIG_PATH", "/path/that/does/not/exist.yaml")

	cfg := LoadConfig()

	require.NotNil(t, cfg, "Конфигурация должна быть инициализирована")

	assert.Equal(t, "in_memory", cfg.Engine.Type, "Тип БД по умолчанию должен быть in_memory")
	
	assert.Equal(t, "127.0.0.1:3223", cfg.TCPServer.Address)
	assert.Equal(t, 100, cfg.TCPServer.MaxConnections)
	assert.Equal(t, "4KB", cfg.TCPServer.MaxMessageSize)
	assert.Equal(t, 5*time.Minute, cfg.TCPServer.IdleTimeout)
	
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "./log/output.log", cfg.Logging.Output)
}

func TestLoadConfig_WithValidFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	yamlContent := []byte(`
engine:
  type: "postgres"
network:
  address: "0.0.0.0:8080"
  max_connections: 500
logging:
  level: "debug"
`)

	err := os.WriteFile(configPath, yamlContent, 0644)
	require.NoError(t, err, "Не удалось создать временный файл конфигурации")

	t.Setenv("CONFIG_PATH", configPath)

	cfg := LoadConfig()

	require.NotNil(t, cfg)


	assert.Equal(t, "postgres", cfg.Engine.Type, "Тип движка должен прочитаться из файла")
	assert.Equal(t, "0.0.0.0:8080", cfg.TCPServer.Address, "Адрес сервера должен прочитаться из файла")
	assert.Equal(t, 500, cfg.TCPServer.MaxConnections)
	assert.Equal(t, "debug", cfg.Logging.Level)

	assert.Equal(t, "4KB", cfg.TCPServer.MaxMessageSize, "Размер сообщения должен fallback'нуться к дефолту")
	assert.Equal(t, 5*time.Minute, cfg.TCPServer.IdleTimeout)
	assert.Equal(t, "./log/output.log", cfg.Logging.Output)
}