package main

import (
	"clon/instructions"
	"clon/services"
)

func main() {
	services.LoadEnv()

	sess := services.ConnectAws()
	instructions.Check(sess)

	//if _, err := utils.Parser.Parse(); err != nil {
	//	switch flagsErr := err.(type) {
	//	case flags.ErrorType:
	//		if flagsErr == flags.ErrHelp {
	//			os.Exit(0)
	//		}
	//		os.Exit(1)
	//	default:
	//		os.Exit(1)
	//	}
	//}
}
