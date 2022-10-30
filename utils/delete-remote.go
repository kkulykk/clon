package utils

import (
	"clon/services"
	"fmt"
)

type DeleteRemoteCommand struct {
	Args struct {
		RemoteName string
	} `positional-args:"yes" required:"1"`
}

var deleteRemoteCommand DeleteRemoteCommand

func (x *DeleteRemoteCommand) Execute([]string) error {
	fmt.Printf("Deleting an existing remote -> :%v\n\n", x.Args.RemoteName)

	sess := services.ConnectAws()

	services.RemoveBucket(sess, x.Args.RemoteName)

	return nil
}

func init() {
	Parser.AddCommand("delete-remote",
		"Delete an existing remote",
		"Delete an existing remote",
		&deleteRemoteCommand)
}
