package utils

import (
	"clon/instructions"
	"clon/services"
)

type SyncCommand struct {
	Soft  bool `short:"s" long:"soft" description:"Soft syncing of files"`
	Force bool `short:"f" long:"force" description:"Force syncing of files"`
	Args  struct {
		LocalPath  string
		RemotePath string
		// TODO! What does positional-args:"no" mean?
	} `positional-args:"no" required:"2"`
}

var syncCommand SyncCommand

func (options *SyncCommand) Execute(args []string) error {
	//remotePath := "clon-demo:"
	//localPath := "./remote"
	sess := services.ConnectAws()
	instructions.Sync(sess, options.Args.LocalPath, options.Args.RemotePath)

	return nil
}

func init() {
	_, err := Parser.AddCommand("sync",
		"Make source and dest identical",
		"Make source and dest identical, modifying destination only. Use -s to soft sync of files \n"+
			" (add only non-existent file by names) or -f to force sync (syncing based on name, size and \n"+
			"modification date of files).",
		&syncCommand)
	if err != nil {
		return
	}
}
