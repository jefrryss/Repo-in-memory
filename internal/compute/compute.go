package compute

import (
	"errors"

	"go.uber.org/zap"

	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
)

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

func (c *Compute) HandleQuery(queryStr string) (string, error) {	
	query, err := c.parser.Parse(queryStr)
	if err != nil {
		c.logger.Error("failde query", zap.Error(err))
		return "", err
	}

	switch query.Cmd {
	case parser.CmdSet:
		err := c.storageInMemory.Set(query.Key, query.Value)
		if err != nil {
			c.logger.Error("failed SET", zap.Error(err), zap.String("key", query.Key))
			return "", err
		}
		return "succes", nil

	case parser.CmdGet:
		value, err := c.storageInMemory.Get(query.Key)
		if err != nil {
			c.logger.Info("failed GET (key not found)", zap.Error(err), zap.String("key", query.Key))
			return "", err
		}
		return value, nil

	case parser.CmdDel:
		err := c.storageInMemory.Del(query.Key)
		if err != nil {
			c.logger.Error("failed DEL", zap.Error(err), zap.String("key", query.Key))
			return "", err
		}
		return "succes", nil

	default:
		c.logger.Error("unknown command passed parser")
		return "", errors.New("internal error: unknown command")
	}
}