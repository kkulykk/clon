package utils

import (
	"clon/instructions"
	"clon/services"
	"fmt"
)

type MoveCommand struct {
	Args struct {
		FromPath string
		ToPath   string
	} `positional-args:"yes" required:"2"`
}

var moveCommand MoveCommand

func (options *MoveCommand) Execute([]string) error {
	fmt.Printf("Moving %v -> %v...\n\n", options.Args.FromPath, options.Args.ToPath)

	fromPath := options.Args.FromPath
	toPath := options.Args.ToPath
	sess := services.ConnectAws()

	instructions.Move(sess, fromPath, toPath)

	return nil
}

func init() {
	Parser.AddCommand("move",
		"Move files from source to dest",
		"Moves the contents of the source directory to the destination directory",
		&moveCommand)
}
