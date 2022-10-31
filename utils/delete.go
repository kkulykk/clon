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
	_, err := Parser.AddCommand("delete",
		"Remove the files in path",
		"Remove the files in path.",
		&deleteCommand)
	if err != nil {
		return
	}
}
