package parser

type Command string


var (
	CmdSet Command = "SET"
	CmdGet Command = "GET"
	CmdDel Command = "DEL"
)

type Query struct {
	Cmd Command
	Key string 
	Value string
}