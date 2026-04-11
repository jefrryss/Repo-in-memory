package compute

import (
	"errors"

	"go.uber.org/zap"
	"context"
	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
)

type ctxKey string
const ClientIpKey string = "ClientIp" 


type Compute struct {
	parser  parser.Parser
	logger *zap.Logger
	storageInMemory storage.Storage
}

func NewCompute(p parser.Parser, s storage.Storage, l *zap.Logger) *Compute {
	return &Compute{
		parser:  p,
		logger:  l,
		storageInMemory: s,
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
		err := c.storageInMemory.Set(ctx, query.Key, query.Value)
		if err != nil {
			c.logger.Error("failed SET", zap.Error(err), zap.String("key", query.Key), zap.String("ip", ip))
			return "", err
		}
		return "succes", nil

	case parser.CmdGet:
		value, err := c.storageInMemory.Get(ctx, query.Key)
		if err != nil {
			c.logger.Info("failed GET (key not found)", zap.Error(err), zap.String("key", query.Key), zap.String("ip", ip))
			return "", err
		}
		return value, nil

	case parser.CmdDel:
		err := c.storageInMemory.Del(ctx, query.Key)
		if err != nil {
			c.logger.Error("failed DEL", zap.Error(err), zap.String("key", query.Key), zap.String("ip", ip))
			return "", err
		}
		return "succes", nil

	default:
		c.logger.Error("unknown command passed parser", zap.String("ip", ip))
		return "", errors.New("internal error: unknown command")
	}
}