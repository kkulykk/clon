package instructions

import "clon/services"

func CreateRemote(remoteName string) {
	sess := services.ConnectAws()

	services.CreateBucket(sess, remoteName)
}

func DeleteRemote(remoteName string) {
	sess := services.ConnectAws()

	services.RemoveBucket(sess, remoteName)
}

func Copy(fromPath string, toPath string) {
	sess := services.ConnectAws()

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
