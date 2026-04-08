package logger

import (
	"go.uber.org/zap"
	"in-memory/config"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"fmt"
	"os"
)

func NewLogger(cnf *config.Config) *zap.Logger{

	zapLevel, err := zapcore.ParseLevel(cnf.Logging.Level)
	if err != nil {
		zapLevel = zapcore.InfoLevel
	}
	zapConfig := zap.NewDevelopmentConfig()
	
	createLogsIfnotExists(cnf.Logging.Output)
	paths := []string{cnf.Logging.Output}
	
	zapConfig.Level = zap.NewAtomicLevelAt(zapLevel)
	zapConfig.OutputPaths = paths
	zapConfig.ErrorOutputPaths = paths
	
	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("Error building logger: %v", err))
	}

	return logger
} 

func createLogsIfnotExists(path string) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir , 0750); err != nil {
		panic(fmt.Sprintf("Ошибка при создании директории для логов '%s': %v", dir, err))
	}

}