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
	//remotePath := "clon-demo:"
	//localPath := "./remote"
	sess := services.ConnectAws()
	instructions.Check(sess, options.Args.LocalPath, options.Args.RemotePath)

	return nil
}

func init() {
	_, err := Parser.AddCommand("check",
		//TODO: update descriptions
		"Make source and dest identical",
		"Make source and dest identical, modifying destination only. Use -s to soft sync of files \n"+
			" (add only non-existent file by names) or -f to force sync (syncing based on name, size and \n"+
			"modification date of files).",
		&checkCommand)
	if err != nil {
		return
	}
}
