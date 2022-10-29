package utils

import (
	"fmt"
)

type CopyCommand struct {
	Force bool `short:"f" long:"force" description:"Force copy of files"`
}

var copyCommand CopyCommand

func (options *CopyCommand) Execute(args []string) error {
	fmt.Printf("Copying: %#v\n", args)
	return nil
}

func init() {
	Parser.AddCommand("copy",
		"Copy contents",
		"Copy file(s) or directories from remote or local or vice versa. Use -f to force copy.",
		&copyCommand)
}
