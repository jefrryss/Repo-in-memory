package compute

import (
	"errors"

	"go.uber.org/zap"

	"in-memory/internal/compute/parser"
	"in-memory/internal/storage"
)

type Compute struct {
	parser  parser.Parser
	storage storage.Storage
	logger  *zap.Logger 
}

func NewCompute(p parser.Parser, s storage.Storage, l *zap.Logger) *Compute {
	return &Compute{
		parser:  p,
		storage: s,
		logger:  l,
	}
}

func (c *Compute) HandleQuery(queryStr string) (string, error) {
	c.logger.Debug("start parsing", zap.String("query", queryStr))
	
	query, err := c.parser.Parse(queryStr)
	if err != nil {
		c.logger.Warn("Invalid query syntax", zap.Error(err), zap.String("query", queryStr))
		return "", err
	}

	switch query.Command {
	case "SET":
		err := c.storage.Set(query.Key, query.Value)
		if err != nil {
			c.logger.Error("failed SET", zap.Error(err), zap.String("key", query.Key))
			return "", err
		}
		return "ok", nil

	case "GET":
		value, err := c.storage.Get(query.Key)
		if err != nil {
			c.logger.Info("failed GET (key not found)", zap.Error(err), zap.String("key", query.Key))
			return "", err
		}
		return value, nil

	case "DEL":
		err := c.storage.Del(query.Key)
		if err != nil {
			c.logger.Error("failed DEL", zap.Error(err), zap.String("key", query.Key))
			return "", err
		}
		return "ok", nil

	default:
		c.logger.Error("unknown command passed parser", zap.String("command", query.Command))
		return "", errors.New("internal error: unknown command")
	}
}