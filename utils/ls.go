package utils

import (
	"clon/instructions"
	"clon/services"
	"fmt"
)

type LsCommand struct {
	Args struct {
		Path string
	} `positional-args:"yes" required:"1"`
}

var lsCommand LsCommand

func (options *LsCommand) Execute([]string) error {
	fmt.Printf("Showing all the files under the path: %v\n\n", options.Args.Path)

	sess := services.ConnectAws()

	instructions.ListElements(sess, options.Args.Path)

	return nil
}

func init() {
	Parser.AddCommand("ls",
		"List the objects in the path with size, path and update date",
		"Lists the objects in the source path to standard output in a human readable format with size path and update date.",
		&lsCommand)
}
