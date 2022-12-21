package utils

import (
	"fmt"
)

type SyncCommand struct {
	Soft  bool `short:"s" long:"soft" description:"Soft syncing of files"`
	Force bool `short:"f" long:"force" description:"Force syncing of files"`
	Args  struct {
		SyncPath string
	} `positional-args:"yes" required:"1"`
}

var syncCommand SyncCommand

func (x *SyncCommand) Execute(args []string) error {

	//sess := services.ConnectAws()
	//instructions.Sync(sess, "clon-test", x.Args.SyncPath)

	fmt.Printf("Synced (force=%v): %#v\n", x.Soft, args)
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
