package utils

import (
	"clon/instructions"
	"clon/services"
	"fmt"
)

type SizeCommand struct {
	Args struct {
		FilePath string
	} `positional-args:"yes" required:"1"`
}

var sizeCommand SizeCommand

func (options *SizeCommand) Execute([]string) error {
	fmt.Printf("Retrieving size of file -> %v\n\n", options.Args.FilePath)

	sess := services.ConnectAws()

	instructions.Size(sess, options.Args.FilePath)

	return nil
}

func init() {
	_, err := Parser.AddCommand("size",
		"Return size of file in bytes",
		"Return size of file in bytes",
		&sizeCommand)
	if err != nil {
		return
	}
}
