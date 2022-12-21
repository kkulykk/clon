package utils

import (
	"clon/instructions"
	"clon/services"
)

type CheckCommand struct {
	Args struct {
		RemotePath string
		LocalPath  string
	} `positional-args:"no" required:"2"`
}

var checkCommand CheckCommand

func (x *CheckCommand) Execute(args []string) error {

	sess := services.ConnectAws()
	instructions.Check(sess, "clon-test", x.Args.LocalPath)

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
