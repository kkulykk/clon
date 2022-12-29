package instructions

import (
	"clon/services"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

type CheckFilesResult struct {
	filesToUpload []string
	filesToDelete []string
}

// CreateRemote : Create a remote entity in the cloud
func CreateRemote(sess *session.Session, remoteName string) {
	err := services.CreateBucket(sess, remoteName)
	if err != nil {
		return
	}
}

// DeleteRemote : Delete a remote entity in the cloud
func DeleteRemote(sess *session.Session, remoteName string) {
	err := services.RemoveBucket(sess, remoteName)
	if err != nil {
		return
	}
}

// Copy : Helper function to copy files between a remote entity in the cloud and local machine
func Copy(sess *session.Session, fromPath string, toPath string) {
	if services.IsRemotePath(fromPath) {
		bucketName := services.GetBucketNameFromRemotePath(fromPath)
		remoteFilePath := services.GetRemoteFilePath(fromPath)

		services.DownloadFile(sess, bucketName, remoteFilePath, toPath)
	} else {
		bucketName := services.GetBucketNameFromRemotePath(toPath)
		remoteFilePath := services.GetRemoteFilePath(toPath)

		services.UploadFileWithChecksum(sess, bucketName, fromPath, remoteFilePath)
	}
}

// Move : Helper function to move files between a remote entity in the cloud and local machine
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

// ListElements : Helper function to retrieve file data from a remote entity in the cloud
func ListElements(sess *session.Session, path string) {
	services.GetBucketItems(sess, path)
}

// Delete : Helper function to remove file(s) from a remote entity in the cloud
func Delete(sess *session.Session, path string) {
	if services.GetBucketNameFromRemotePath(path) == "" {
		services.ExitErrorf("You should provide remote name to delete something")
	}

	bucket := services.GetBucketNameFromRemotePath(path)

	if services.GetRemoteFilePath(path) == "/" {
		fmt.Println("Detected an attempt to remove all contents from remote.")
		fmt.Println("After performing the operation, you will not be able to restore the data.")

		if services.Confirm("this operation") {
			err := services.DeleteAll(sess, bucket)
			if err != nil {
				return
			}
		} else {
			os.Exit(1)
		}
	} else if services.IsDirectory(path) && path != bucket {
		fmt.Printf("Detected an attempt to remove all contents from a directory %v\n", path)
		fmt.Println("After performing the operation, you will not be able to restore the data.")

		if services.Confirm("this operation") {
			err := services.DeleteDirectory(sess, path)
			if err != nil {
				return
			}
		} else {
			os.Exit(1)
		}
	} else {
		err := services.DeleteFile(sess, path)
		if err != nil {
			return
		}
	}
}

// Size : Helper function to retrieve file(s) size from a remote entity in the cloud
func Size(sess *session.Session, filePath string) {
	if !services.IsRemotePath(filePath) {
		fmt.Println("Please, enter remote file path")

		return
	}

	bucketName := services.GetBucketNameFromRemotePath(filePath)
	remoteFilePath := services.GetRemoteFilePath(filePath)

	services.GetBucketFileSize(sess, bucketName, remoteFilePath)
}

// GetRemotes : Helper function to retrieve all remote entities from the cloud
func GetRemotes(sess *session.Session) {
	services.GetBucketsList(sess)
}

// Check : Helper function to check if local and remote directories are up-to-date
func Check(sess *session.Session, localPath string, remotePath string) {
	if !services.RemotePathExists(sess, remotePath) {
		services.ExitErrorf("Remote path %q does not exist.\nCheck the correctness of input or arguments order", remotePath)
	}

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		services.ExitErrorf("Local path %q does not exist.\nCheck the correctness of input or arguments order", localPath)
	}

	bucket := services.GetBucketNameFromRemotePath(remotePath)
	files := services.CheckFiles(sess, bucket, localPath, remotePath)

	if len(files.FilesToUpload) == 0 && len(files.FilesToDelete) == 0 {
		fmt.Println("Remote and local paths are up-to-date")

		return
	} else {
		fmt.Println("Remote and local paths are NOT synchronized\nCheck a difference below")
	}

	if len(files.FilesToUpload) > 0 {
		fmt.Println("Files to upload:")

		for _, fileToUpdate := range files.FilesToUpload {
			color.Green("	upload: %q", fileToUpdate)
		}

		fmt.Println()
	}

	if len(files.FilesToDelete) > 0 {
		fmt.Println("Files to delete:")

		for _, fileToDelete := range files.FilesToDelete {
			color.Red("	delete: %q", fileToDelete)
		}

		fmt.Println()
	}
}

// Sync : Helper function to make remote directory up-to-date with local
func Sync(sess *session.Session, localPath string, remotePath string) {
	bucket := services.GetBucketNameFromRemotePath(remotePath)
	remotePathPrefix := services.GetRemotePathPrefix(remotePath)
	localPathPrefix := services.GetLocalPathPrefix(localPath)
	files := services.CheckFiles(sess, bucket, localPath, remotePath)

	if len(files.FilesToUpload) > 0 {
		fmt.Println("Files to upload:")
		for _, fileToUpdate := range files.FilesToUpload {
			color.Green("	upload: %q", fileToUpdate)
		}
		if services.Confirm("uploading this file(s)") {

			for _, fileToUpdate := range files.FilesToUpload {
				splitRemoteFileDirectory := strings.Split(fileToUpdate, "/")
				// We should add / at the end of remoteFileDirectory to make it match with localPathPrefix and be able to replace it later
				// With remotePathPrefix
				remoteFileDirectory := strings.Join(splitRemoteFileDirectory[:len(splitRemoteFileDirectory)-1], "/") + "/"

				// We should replace localPathPrefix with remotePathPrefix to set correct bucket path
				services.UploadFileWithChecksum(sess, bucket, fileToUpdate, strings.Replace(remoteFileDirectory, localPathPrefix, remotePathPrefix, 1))
			}
		} else {
			os.Exit(1)
		}
	}

	if len(files.FilesToDelete) > 0 {
		fmt.Println("Files to delete:")
		for _, fileToDelete := range files.FilesToDelete {
			color.Red("	delete: %q", fileToDelete)
		}

		if services.Confirm("deleting this file(s)") {
			for _, fileToDelete := range files.FilesToDelete {
				remoteFilePathToDelete := bucket + ":" + strings.Join(strings.Split(fileToDelete, "/"), ":")
				err := services.DeleteFile(sess, remoteFilePathToDelete)
				if err != nil {
					fmt.Printf("Error while deleting %q file\n", fileToDelete)
				}
			}
		} else {
			os.Exit(1)
		}
	}

	if len(files.FilesToUpload) == 0 && len(files.FilesToDelete) == 0 {
		fmt.Println("No files need to be synchronized")
	}
}
