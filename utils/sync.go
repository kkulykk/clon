package utils

import (
	"fmt"
)

type SyncCommand struct {
	Soft  bool `short:"s" long:"soft" description:"Soft syncing of files"`
	Force bool `short:"f" long:"force" description:"Force syncing of files"`
}

var syncCommand SyncCommand

func (x *SyncCommand) Execute(args []string) error {
	fmt.Printf("Syncing (force=%v): %#v\n", x.Force, args)
	return nil
}

func init() {
	Parser.AddCommand("sync",
		"Make source and dest identical",
		"Make source and dest identical, modifying destination only. Use -s to soft sync of files \n"+
			" (add only non-existent file by names) or -f to force sync (syncing based on name, size and \n"+
			"modification date of files).",
		&syncCommand)
}
