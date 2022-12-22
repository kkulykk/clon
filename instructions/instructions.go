package instructions

import (
	"clon/services"
	"crypto/md5"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
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

		services.UploadFile(sess, bucketName, fromPath, remoteFilePath)
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

		if services.Confirm() {
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

		if services.Confirm() {
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
// func Check(sess *session.Session, bucket string, path string) {
func Check(sess *session.Session) {
	bucket := "clon-demo"
	remotePath := "clon-demo"
	localPath := "./remote"
	remotePathPrefix := services.GetRemoteFilePathPrefix(remotePath)
	localPathPrefix := localPath

	if !strings.HasSuffix(remotePathPrefix, "/") {
		remotePathPrefix = remotePathPrefix + "/"
	}

	if !strings.HasSuffix(localPathPrefix, "/") {
		localPathPrefix = localPathPrefix + "/"
	}
	// Remove ./ from local file prefix if its path starts with it
	if strings.HasPrefix(localPathPrefix, "./") {
		localPathPrefix = strings.Replace(localPathPrefix, "./", "", 1)
	}

	var filesToUpdate []string
	remoteFiles, _ := services.GetAwsS3ItemMap(sess, bucket, remotePath)
	remoteFilesPaths := make([]string, len(remoteFiles))

	i := 0
	for remotePath := range remoteFiles {
		remoteFilesPaths[i] = remotePath
		i++
	}

	filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			contents, err := os.ReadFile(path)

			if err == nil {
				// Should remote /remote from the beginning of the local path to conform with remote path
				localFilePathOnRemote := strings.Replace(path, localPathPrefix, "", 1)
				localFileMd5Sum := md5.Sum(contents)
				localFileChecksum := fmt.Sprintf("%x", localFileMd5Sum)

				if localFileChecksum != remoteFiles[localFilePathOnRemote] {
					filesToUpdate = append(filesToUpdate, path)
				}

				filteredRemoteFilesPaths := make([]string, 0)

				for _, remoteFilePath := range remoteFilesPaths {
					if localFilePathOnRemote != remoteFilePath {
						filteredRemoteFilesPaths = append(filteredRemoteFilesPaths, remoteFilePath)
					}
				}

				remoteFilesPaths = filteredRemoteFilesPaths
			} else {
				services.ExitErrorf("Error reading file %q", path)
			}
		}

		return nil
	})

	//fmt.Println("filesToUpdate", filesToUpdate)
	//fmt.Println("filesToDelete", remoteFilesPaths)

	if len(filesToUpdate) > 0 {
		fmt.Println("Files to update:\n")

		for _, fileToUpdate := range filesToUpdate {
			color.Green("	update: %q", fileToUpdate)
		}

		fmt.Println()
	}

	if len(remoteFilesPaths) > 0 {
		fmt.Println("Files to delete:")

		for _, fileToDelete := range remoteFilesPaths {
			color.Red("	delete: %q", fileToDelete)
		}

		fmt.Println()
	}

	return

	//remoteFiles, err := services.GetAwsS3ItemMap(sess, bucket)
	//files := services.CheckFiles(remoteFiles, path)
	//
	//if len(files.FilesToDelete) == 0 && len(files.FilesToUpload) == 0 {
	//	fmt.Println("Local and remote storages are up to date.")
	//}
	//
	//if len(files.FilesToDelete) > 0 && len(files.FilesToUpload) == 0 {
	//	fmt.Println("No files need to be updated. The following files need to be deleted on remote: ")
	//	fmt.Println(files.FilesToDelete)
	//}
	//
	//if len(files.FilesToUpload) > 0 && len(files.FilesToDelete) == 0 {
	//	fmt.Println("No files need to be deleted on remote. The following files need to be updated on remote: ")
	//	fmt.Println(files.FilesToUpload)
	//}
	//
	//if len(files.FilesToUpload) > 0 && len(files.FilesToDelete) > 0 {
	//	fmt.Println("The following files need to be updated on remote: ")
	//	fmt.Println(files.FilesToUpload)
	//	fmt.Println("The following files need to be deleted on remote: ")
	//	fmt.Println(files.FilesToUpload)
	//}
	//
	//if err != nil {
	//	return
	//}
}
