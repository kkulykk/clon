package utils

import (
	"clon/instructions"
	"clon/services"
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

	sess := services.ConnectAws()

	instructions.CreateRemote(sess, options.Args.RemoteName)

	return nil
}

func init() {
	_, err := Parser.AddCommand("create-remote",
		"Create a new remote with name",
		"Create a new remote with name.",
		&createRemoteCommand)
	if err != nil {
		return
	}
}
