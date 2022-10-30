package utils

import (
	"clon/instructions"
	"fmt"
)

type CreateRemoteCommand struct {
	Args struct {
		RemoteName string
	} `positional-args:"yes" required:"1"`
}

var createRemoteCommand CreateRemoteCommand

func (options *CreateRemoteCommand) Execute([]string) error {
	fmt.Printf("Creating new remote -> %v\n\n", options.Args.RemoteName)

	instructions.CreateRemote(options.Args.RemoteName)

	return nil
}

func init() {
	Parser.AddCommand("create-remote",
		"Create a new remote with name",
		"Create a new remote with name.",
		&createRemoteCommand)
}
