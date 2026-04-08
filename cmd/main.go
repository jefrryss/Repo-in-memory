package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"in-memory/internal/compute"
	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
	"in-memory/internal/storage/engine"
	"in-memory/config"
	"in-memory/internal/logger"
)



func main() {
	cnf := config.LoadConfig()
	logger := logger.NewLogger(cnf)
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