package main

import (
	"context"
	"fmt"
	"in-memory/config"
	"in-memory/internal/compute"
	"in-memory/internal/compute/parser"
	"in-memory/internal/logger"
	"in-memory/internal/server"
	"in-memory/internal/storage"
	"in-memory/internal/storage/engine"
	"in-memory/internal/wal"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cnf := config.LoadConfig()
	logger := logger.NewLogger(cnf)
	defer logger.Sync()

	logger.Info("init logs")

	prs := parser.NewLineParser()
	logger.Info("init parser")

	hashTable := engine.NewHashTable()
	logger.Info("init hashTable")

	store := storage.NewStorage(hashTable)
	logger.Info("init storage")

	w, err := wal.NewWal(cnf, logger)
	if err != nil {
		log.Fatal("Failed initialisastion WAL", zap.Error(err))
	}

	if w != nil {
		bridge := make(chan *parser.Query)

		go func() {
			err := w.Recovery(bridge)
			if err != nil {
				log.Fatal("Failed recovery from WAL", zap.Error(err))
			}
		}()

		ctx := context.Background()

		for query := range bridge {
			switch query.Cmd {
			case parser.CmdSet:
				store.Set(ctx, query.Key, query.Value)
			case parser.CmdDel:
				store.Del(ctx, query.Key)
			}
		}
	}

	comp := compute.NewCompute(prs, store, logger, w)

	logger.Info("init db")

	server := server.NewServerTSP(&cnf.TCPServer, comp, logger)

	go server.StartServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch
	logger.Info("Received shutdown signal. Starting graceful shutdown...")

	shutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = server.ShutDown(shutCtx)
	fmt.Println("\nServer stoped")
	if err != nil {
		logger.Error("Server shutdown with error", zap.Error(err))
	} else {
		logger.Info("Server shutdown gracefully")
	}

}
