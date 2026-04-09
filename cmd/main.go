package main

import (
	"in-memory/internal/compute"
	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
	"in-memory/internal/storage/engine"
	"in-memory/config"
	"in-memory/internal/logger"
	"in-memory/internal/server"
)



func main() {
	cnf := config.LoadConfig()
	logger := logger.NewLogger(cnf)
	
	logger.Info("init logs")
	
	prs := parser.NewLineParser()
	hashTable := engine.NewHashTable() 
	store := storage.NewStorage(hashTable)
	
	comp := compute.NewCompute(prs, store, logger)
	
	logger.Info("init db")
	server := server.NewServerTSP(&cnf.TCPServer, comp, logger)
	server.StartServer()
}