package utils

import (
	"clon/instructions"
	"clon/services"
	"fmt"
)

type DeleteRemoteCommand struct {
	Args struct {
		RemoteName string
	} `positional-args:"yes" required:"1"`
}

var deleteRemoteCommand DeleteRemoteCommand

func (options *DeleteRemoteCommand) Execute([]string) error {
	fmt.Printf("Deleting an existing remote -> %v\n\n", options.Args.RemoteName)

	sess := services.ConnectAws()

	instructions.DeleteRemote(sess, options.Args.RemoteName)

	return nil
}

func init() {
	Parser.AddCommand("delete-remote",
		"Delete an existing remote",
		"Delete an existing remote",
		&deleteRemoteCommand)
}
