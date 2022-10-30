package utils

import (
	"clon/instructions"
	"clon/services"
	"fmt"
)

type RemotesCommand struct {
}

var remotesCommand RemotesCommand

func (options *RemotesCommand) Execute([]string) error {
	fmt.Println("Retrieving list of available remotes ...\n")
	fmt.Println("Remotes:")

	sess := services.ConnectAws()

	instructions.GetRemotes(sess)

	return nil
}

func init() {
	Parser.AddCommand("remotes",
		"Get list of available remotes",
		"Get list of available remotes",
		&remotesCommand)
}
