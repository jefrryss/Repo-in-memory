package compute

import(
	"in-memory/internal/compute/parser"
)

type Compute struct {
	parser parser.Parser
}

func NewCompute(p parser.Parser) *Compute{
	return &Compute{
		Parser: p,
	}
}