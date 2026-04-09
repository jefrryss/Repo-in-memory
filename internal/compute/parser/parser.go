package parser

import (
	"strings"
	"errors"
)

var (
	ErrNotEnoughArgs = errors.New("not enought arguments")
	ErrInvalidCommand = errors.New("invalid command")
	ErrEmptyQuery = errors.New("empty query")
)


type Parser interface {
	Parse(val string) (*Query, error)
}

type LineParser struct {}


func NewLineParser() Parser {
	return &LineParser{}
}


func (l *LineParser) Parse(val string) (*Query, error) {
	command := strings.Fields(val)

	if len(command) == 0 {
		return nil, ErrEmptyQuery
	}

	cmd := command[0]

	switch cmd {
	case "SET":
		if len(command) != 3 {
			return nil, ErrNotEnoughArgs
		}

		query := &Query{
			Cmd: CmdSet,
			Key: command[1],
			Value: command[2],
		}
		return query, nil
	case "GET":
		if len(command) != 2 {
			return nil, ErrNotEnoughArgs
		}

		query := &Query{
			Cmd: CmdGet,
			Key: command[1],
		}
		return query, nil
	case "DEL":
		if len(command) != 2 {
			return nil, ErrNotEnoughArgs
		}

		query := &Query{
			Cmd: CmdDel,
			Key: command[1],
		}
		return query, nil
	default:
		return nil, ErrInvalidCommand
	}
}