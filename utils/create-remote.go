package utils

import (
	"clon/services"
	"fmt"
)

type CreateRemoteCommand struct {
	Args struct {
		RemoteName string
	} `positional-args:"yes" required:"1"`
}

var createRemoteCommand CreateRemoteCommand

func (x *CreateRemoteCommand) Execute([]string) error {
	fmt.Printf("Creating new remote -> :%v\n\n", x.Args.RemoteName)

	sess := services.ConnectAws()

	services.CreateBucket(sess, x.Args.RemoteName)

	return nil
}

func init() {
	Parser.AddCommand("create-remote",
		"Create a new remote with name",
		"Create a new remote with name.",
		&createRemoteCommand)
}
