package utils

import (
	"clon/instructions"
	"clon/services"
	"fmt"
)

type DeleteCommand struct {
	Args struct {
		RemotePath string
	} `positional-args:"yes"`
}

var deleteCommand DeleteCommand

func (options *DeleteCommand) Execute([]string) error {
	fmt.Printf("Deleting all the files under the path: %v\n\n", options.Args.RemotePath)

	sess := services.ConnectAws()

	instructions.Delete(sess, options.Args.RemotePath)

	return nil
}

func init() {
	Parser.AddCommand("delete",
		"List the objects in the path with size, path and update date",
		"Lists the objects in the source path to standard output in a human readable format with size path and update date.",
		&deleteCommand)
}
