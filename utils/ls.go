package utils

import (
	"clon/services"
	"fmt"
)

type LsCommand struct {
	Args struct {
		Path string
	} `positional-args:"yes" required:"1"`
}

var lsCommand LsCommand

func (x *LsCommand) Execute([]string) error {
	fmt.Printf("Showing all the files under the path: %v\n\n", x.Args.Path)

	sess := services.ConnectAws()

	services.GetBucketItems(sess, x.Args.Path)

	return nil
}

func init() {
	Parser.AddCommand("ls",
		"List the objects in the path with size, path and update date",
		"Lists the objects in the source path to standard output in a human readable format with size path and update date.",
		&lsCommand)
}
