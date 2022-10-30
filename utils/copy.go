package utils

import (
	"clon/services"
	"fmt"
)

type CopyCommand struct {
	Force bool `short:"f" long:"force" description:"Force copy of files"`
	Args  struct {
		FromPath string
		ToPath   string
	} `positional-args:"yes" required:"2"`
}

var copyCommand CopyCommand

func (options *CopyCommand) Execute(args []string) error {
	fmt.Printf("Copying %v -> %v\n\n", options.Args.FromPath, options.Args.ToPath)

	sess := services.ConnectAws()

	fromPath := options.Args.FromPath
	toPath := options.Args.ToPath

	if services.IsRemotePath(fromPath) {
		bucketName := services.GetBucketNameFromRemotePath(fromPath)
		remoteFilePath := services.GetRemoteFilePath(fromPath)
		services.DownloadFile(sess, bucketName, remoteFilePath, toPath)
	} else {
		bucketName := services.GetBucketNameFromRemotePath(toPath)
		remoteFilePath := services.GetRemoteFilePath(toPath)
		services.UploadFile(sess, bucketName, fromPath, remoteFilePath)
	}

	return nil
}

func init() {
	Parser.AddCommand("copy",
		"Copy contents",
		"Copy file(s) or directories from remote or local or vice versa. Use -f to force copy.",
		&copyCommand)
}
