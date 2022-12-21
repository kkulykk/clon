package services

import (
	"crypto/md5"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string

type CheckFilesResult struct {
	FilesToUpload []string
	FilesToDelete []string
}

// GetEnvWithKey : Get environmental variable value by name
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

// LoadEnv : Set up environmental variables from .env file
func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

// ConnectAws : Return connection session object
func ConnectAws() *session.Session {
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	if err != nil {
		panic(err)
	}

	return sess
}

// ExitErrorf : Helper function for better error handling
func ExitErrorf(msg string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, msg+"\n", args...)
	if err != nil {
		return
	}
	os.Exit(1)
}

// GetFileNameByPath : Return a file name from a given path
func GetFileNameByPath(path string) string {
	splitPath := strings.Split(path, "/")
	return splitPath[len(splitPath)-1]
}

// IsRemotePath : Return true if a given path is remote
func IsRemotePath(path string) bool {
	return strings.Contains(path, ":")
}

// GetBucketNameFromRemotePath : Return a bucket name from a given path
func GetBucketNameFromRemotePath(path string) string {
	return strings.Split(path, ":")[0]
}

// GetRemoteFilePathPrefix : Return a prefix for a path
func GetRemoteFilePathPrefix(path string) string {
	// Remove bucket name and replace : with /
	return strings.Join(strings.Split(path, ":")[1:], "/")
}

// GetRemoteFilePath : Return a path without a bucket name
func GetRemoteFilePath(path string) string {
	// Remove bucket name and replace : with /
	return "/" + GetRemoteFilePathPrefix(path)
}

// Confirm : Return true if a user confirms the action
func Confirm() bool {
	var input string

	fmt.Printf("Do you want to continue with this operation? [y|n]: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		panic(err)
	}
	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	}
	return false
}

// IsDirectory : Return true if a given path is a directory
func IsDirectory(path string) bool {
	return strings.HasSuffix(path, ":") || strings.HasSuffix(path, "/")
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// CheckFiles : Iterate through all the files and compare checksums
func CheckFiles(remoteItems map[string]string, paths string) *CheckFilesResult {
	var filesToUpdate []string
	var filesToDelete []string
	var actualFiles []string

	remoteFiles := make([]string, 0, len(remoteItems))

	for k := range remoteItems {
		remoteFiles = append(remoteFiles, k)
	}

	var numfiles int
	result := CheckFilesResult{}

	fmt.Println(paths)

	filepath.Walk(paths, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			numfiles++
			relativeFile := strings.Replace(path, paths, "", 1)
			fmt.Println(remoteItems)
			checksumRemote, _ := remoteItems[relativeFile]

			validateChecksum(path, checksumRemote, filesToUpdate, actualFiles)
		}

		return nil
	})

	filesToDelete = difference(remoteFiles, append(actualFiles, filesToUpdate...))

	fmt.Println(filesToUpdate)
	fmt.Println(filesToDelete)
	fmt.Println(actualFiles)

	result.FilesToUpload = filesToUpdate
	result.FilesToDelete = filesToDelete

	return &result
}

func validateChecksum(filename string, checksumRemote string, filesToUpdate []string, actualFiles []string) {

	if checksumRemote == "" {
		filesToUpdate = append(filesToUpdate, filename)
		return
	}

	contents, err := os.ReadFile(filename)
	if err == nil {
		sum := md5.Sum(contents)
		sumString := fmt.Sprintf("%x", sum)
		if sumString != checksumRemote {
			filesToUpdate = append(filesToUpdate, filename)
			return
		} else {
			actualFiles = append(actualFiles, filename)
			return
		}
	}
}
