package utils

import "github.com/jessevdk/go-flags"

type Options struct {
	// Example of verbosity with level
	//Help []bool `short:"h" long:"help" description:"Show help message" optional:"yes"`
}

var options Options
var Parser = flags.NewParser(&options, flags.Default)
