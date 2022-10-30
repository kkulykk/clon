package instructions

import (
	"clon/services"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

func CreateRemote(sess *session.Session, remoteName string) {
	services.CreateBucket(sess, remoteName)
}

func DeleteRemote(sess *session.Session, remoteName string) {
	services.RemoveBucket(sess, remoteName)
}

func Copy(sess *session.Session, fromPath string, toPath string) {
	if services.IsRemotePath(fromPath) {
		bucketName := services.GetBucketNameFromRemotePath(fromPath)
		remoteFilePath := services.GetRemoteFilePath(fromPath)

		services.DownloadFile(sess, bucketName, remoteFilePath, toPath)
	} else {
		bucketName := services.GetBucketNameFromRemotePath(toPath)
		remoteFilePath := services.GetRemoteFilePath(toPath)

		services.UploadFile(sess, bucketName, fromPath, remoteFilePath)
	}
}

func Move(sess *session.Session, fromPath string, toPath string) {
	Copy(sess, fromPath, toPath)

	if services.IsRemotePath(fromPath) {
		bucketName := services.GetBucketNameFromRemotePath(fromPath)
		remoteFilePath := services.GetRemoteFilePath(fromPath)

		services.DeleteBucketFile(sess, bucketName, remoteFilePath)
	} else {
		e := os.Remove(fromPath)

		if e != nil {
			log.Fatal(e)
		}
	}
}
