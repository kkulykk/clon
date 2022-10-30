package utils

import (
	"fmt"
)

type AboutCommand struct {
}

var aboutCommand AboutCommand

func (options *AboutCommand) Execute([]string) error {
	fmt.Printf("   ____   _                     \n  / ___| | |   ___    _ __      \n | |     | |  / _ \\  | '_ \\     \n | |___  | | | (_) | | | | |  _ \n  \\____| |_|  \\___/  |_| |_| (_)\n                                \n")
	fmt.Printf("Version 1.0\n")
	fmt.Printf("(C) Bohdan Mykhailiv, Roman Kulyk  2022\n\n")
	fmt.Printf("Clon is a command line program to easily manage files on cloud storage.\n")
	fmt.Printf("The program was created as a part of the final project on Operational Systems course at UCU.\n\n")
	fmt.Printf("To get started, read the help message: clon -h.\n\n")
	return nil
}

func init() {
	Parser.AddCommand("about",
		"Information about the program",
		"Get detailed information about the Clon program.",
		&aboutCommand)
}
