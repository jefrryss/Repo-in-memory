package compute

import (
	"errors"

	"go.uber.org/zap"
	"context"
	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
	"in-memory/internal/wal"
)

type ctxKey string
const ClientIpKey ctxKey = "ClientIp" 


type Compute struct {
	parser  parser.Parser
	logger *zap.Logger
	storageInMemory storage.Storage
	wal wal.WAL
}

func NewCompute(p parser.Parser, s storage.Storage, l *zap.Logger, w wal.WAL) *Compute {
	return &Compute{
		parser:  p,
		logger:  l,
		storageInMemory: s,
		wal: w,
	}
}

func (c *Compute) HandleQuery(ctx context.Context, queryStr string) (string, error) {
	ip, ok := ctx.Value(ClientIpKey).(string)
	if !ok {
		ip = "unknown"
	}

	query, err := c.parser.Parse(ctx, queryStr)
	if err != nil {
		c.logger.Error("failde query", zap.Error(err), zap.String("ip", ip))
		return "", err
	}

	switch query.Cmd {
	case parser.CmdSet:
		if c.wal != nil {
			err := c.wal.Write(ctx, query) 
			if err != nil {
				c.logger.Error("failed write to WAL", zap.String("key", query.Key), zap.String("ip", ip), zap.Error(err))
				return "", err
			}
		}
		err = c.storageInMemory.Set(ctx, query.Key, query.Value)
		if err != nil {
			c.logger.Error("failed SET", zap.Error(err), zap.String("key", query.Key), zap.String("ip", ip))
			return "", err
		}
		return "success", nil

	case parser.CmdGet:
		value, err := c.storageInMemory.Get(ctx, query.Key)
		if err != nil {
			c.logger.Info("failed GET (key not found)", zap.Error(err), zap.String("key", query.Key), zap.String("ip", ip))
			return "", err
		}
		return value, nil

	case parser.CmdDel:
		if c.wal != nil {
			err := c.wal.Write(ctx, query) 
			if err != nil {
				c.logger.Error("failed write to WAL", zap.String("ip", ip), zap.String("key", query.Key), zap.Error(err))
				return "", err
			}
		}
		err = c.storageInMemory.Del(ctx, query.Key)
		if err != nil {
			c.logger.Error("failed DEL", zap.Error(err), zap.String("key", query.Key), zap.String("ip", ip))
			return "", err
		}
		return "success", nil

	default:
		c.logger.Error("unknown command passed parser", zap.String("ip", ip))
		return "", errors.New("internal error: unknown command")
	}
}