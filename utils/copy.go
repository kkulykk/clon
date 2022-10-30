package utils

import (
	"clon/instructions"
	"fmt"
)

type CopyCommand struct {
	Force bool `short:"f" long:"force" description:"Force copy of files"`
	Args  struct {
		FromPath string
		ToPath   string
	} `positional-args:"yes" required:"2"`
}

var copyCommand CopyCommand

func (options *CopyCommand) Execute(args []string) error {
	fmt.Printf("Copying %v -> %v\n\n", options.Args.FromPath, options.Args.ToPath)

	fromPath := options.Args.FromPath
	toPath := options.Args.ToPath

	instructions.Copy(fromPath, toPath)

	return nil
}

func init() {
	Parser.AddCommand("copy",
		"Copy contents",
		"Copy file(s) or directories from remote or local or vice versa. Use -f to force copy.",
		&copyCommand)
}
