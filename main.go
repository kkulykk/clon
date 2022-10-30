package main

import (
	"clon/utils"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	if _, err := utils.Parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
