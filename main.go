package main

import (
	"clon/instructions"
	"clon/services"
)

func main() {
	services.LoadEnv()

	sess := services.ConnectAws()
	toPath := "clon-demo:folder3"
	fromPath := "./remote/folder5/subfolder1/subfolder2/subfolder3/subfolder4/subfolder5/subfolder6/subfolder7.4/ls.go"

	//fmt.Println("GetLocalFilePaths", services.GetLocalFilePaths(localPath))
	//fmt.Println("GetRemoteFilePaths", services.GetRemoteFilePaths(sess, remotePath))

	instructions.Move(sess, fromPath, toPath)

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
