package utils

import "github.com/jessevdk/go-flags"

type Options struct{}

var options Options
var Parser = flags.NewParser(&options, flags.Default)
