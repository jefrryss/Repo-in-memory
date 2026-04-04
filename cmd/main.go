package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"in-memory/internal/compute"
	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
	"in-memory/internal/storage/engine"
)

func InitLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()

	config.OutputPaths = []string{"db.log"}
	config.ErrorOutputPaths = []string{"db.log"}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func main() {
	logger := InitLogger()
	defer logger.Sync()
	
	logger.Info("start init logs")
	
	prs := parser.NewLineParser()
	mem := engine.NewMemory() 
	store := storage.NewStorage(mem)
	
	comp := compute.NewCompute(prs, store, logger)
	
	logger.Info("start init db")
	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "exit" {
			logger.Info("end of work db")
			return
		}
		
		result, err := comp.HandleQuery(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Result: %s\n", result)
		}
	}
}