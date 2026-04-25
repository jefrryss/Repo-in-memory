package parser

import (
	"fmt"
	"strings"
)

type Command string

const (
	CmdSet Command = "SET"
	CmdGet Command = "GET"
	CmdDel Command = "DEL"
)

type Query struct {
	Cmd   Command
	Key   string
	Value string
}

func (q *Query) Marshal() string {
	return fmt.Sprintf("%s %s %s\n", q.Cmd, q.Key, q.Value)
}

func UnMarshal(str string) *Query {
	str = strings.TrimSpace(str) 
	parts := strings.Fields(str) 
	
	if len(parts) < 2 {
		return nil 
	}

	switch parts[0] {
	case string(CmdSet):
		if len(parts) < 3 { return nil }
		return &Query{Cmd: CmdSet, Key: parts[1], Value: parts[2]}
	case string(CmdGet):
		return &Query{Cmd: CmdGet, Key: parts[1]}
	case string(CmdDel):
		return &Query{Cmd: CmdDel, Key: parts[1]}
	default:
		return nil
	}
}