package utils

import (
	"clon/instructions"
	"clon/services"
)

type CheckCommand struct {
	Args struct {
		LocalPath  string
		RemotePath string
	} `positional-args:"no" required:"2"`
}

var checkCommand CheckCommand

func (options *CheckCommand) Execute(args []string) error {
	sess := services.ConnectAws()
	instructions.Check(sess, options.Args.LocalPath, options.Args.RemotePath)

	return nil
}

func init() {
	_, err := Parser.AddCommand("check",
		"Checks the files in the source and destination match",
		"Checks the files in the source and destination match. It compares sizes and hashes (MD5 or SHA1) and logs a report of files that don't match. It doesn't alter the source or destination.",
		&checkCommand)
	if err != nil {
		return
	}
}
