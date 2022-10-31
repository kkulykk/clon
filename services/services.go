package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string

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
