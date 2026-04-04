package compute

import(
	"in-memory/internal/compute/parser"
	"in-memory/internal/compute/storage"
)

type Compute struct {
	parser *parser.Parser
	storage *storage.Storage
	loger *zap.Logger
}

func NewCompute(p parser.Parser, s storage.Storage, l *zap.Logger) *Compute{
	return &Compute{
		parser: p,
		storage: s,
		loger: l,
	}
}

func (c *Compute) HandleRead(query string) error {

}


