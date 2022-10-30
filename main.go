package main

import (
	"clon/services"
	"clon/utils"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	services.LoadEnv()

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

	//services.GetBucketsList(Sess)
	//getBucketItems(sess, "clon-demo")
	//createBucket(sess, "clon-demo")
	//uploadFile(sess, "clon-demo", "./main.go", "./")
	//uploadFile(sess, "clon-demo", "./image.jpeg", "./img/")

	//downloadFile(sess, "clon-demo", "./img/image.jpeg", "./")
	//deleteBucketFile(sess, "clon-demo", "./img/image.jpeg")

	//awsAccessKeyID := services.GetEnvWithKey("AWS_ACCESS_KEY_ID")
	//fmt.Println("My access key ID is ", awsAccessKeyID)
}
