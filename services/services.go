package services

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"io"
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
func Confirm(action string) bool {
	var input string

	fmt.Printf("Do you want to continue with %s? [y|n]: ", action)
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

// ValidateChecksum : Check if checksum of local file matches with remote file checksum
func ValidateChecksum(localFileContent []byte, remoteFileChecksum string) bool {
	localFileMd5Sum := md5.Sum(localFileContent)
	localFileChecksum := fmt.Sprintf("%x", localFileMd5Sum)
	// TODO: fix checksum considering encryption

	if localFileChecksum == remoteFileChecksum {
		return true
	}

	return false
}

// GetRemotePathPrefix : Return prefix of file with remote path
func GetRemotePathPrefix(remotePath string) string {
	remotePathPrefix := GetRemoteFilePathPrefix(remotePath)

	if !strings.HasSuffix(remotePathPrefix, "/") && remotePathPrefix != "" {
		remotePathPrefix = remotePathPrefix + "/"
	}

	return remotePathPrefix
}

// GetLocalPathPrefix : Return prefix of file with local path
func GetLocalPathPrefix(localPath string) string {
	localPathPrefix := localPath

	if !strings.HasSuffix(localPathPrefix, "/") {
		localPathPrefix = localPathPrefix + "/"
	}

	// Remove ./ from local file prefix if its path starts with it
	if strings.HasPrefix(localPathPrefix, "./") {
		localPathPrefix = strings.Replace(localPathPrefix, "./", "", 1)
	}

	return localPathPrefix
}

// ShouldIgnoreFile Check if path starts with paths in .clonignore file
func ShouldIgnoreFile(suffixesToIgnore []string, pathToCheck string) bool {
	for _, suffixToIgnore := range suffixesToIgnore {
		// Check if suffixToIgnore is directory path
		if strings.HasSuffix(suffixToIgnore, "/") {
			if strings.HasPrefix(pathToCheck, suffixToIgnore) {
				return true
			}
		} else {
			if pathToCheck == suffixToIgnore {
				return true
			}
		}
	}

	return false
}

// CheckFiles : Iterate through all the files and compare checksums
func CheckFiles(sess *session.Session, bucket string, localPath string, remotePath string) CheckFilesResult {
	checkFilesResult := CheckFilesResult{}
	remotePathPrefix := GetRemotePathPrefix(remotePath)
	localPathPrefix := GetLocalPathPrefix(localPath)
	clonignoreFilePath := localPath + "/" + ".clonignore"
	remoteFiles, _ := GetAwsS3ItemMap(sess, bucket, remotePath)
	remoteFilesPaths := make([]string, len(remoteFiles))
	var clonignoreFilesPathsToSkip []string
	var filesToUpdate []string

	i := 0
	for remotePath := range remoteFiles {
		remoteFilesPaths[i] = remotePath
		i++
	}

	// Check if .clonignore file exists
	if _, err := os.Stat(clonignoreFilePath); err == nil {
		readFile, err := os.Open(clonignoreFilePath)

		if err != nil {
			fmt.Println("Error reading .clonignore file")
		}

		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			clonignoreFilesPathsToSkip = append(clonignoreFilesPathsToSkip, localPathPrefix+strings.TrimSpace(fileScanner.Text()))
		}

		readFileCloseErr := readFile.Close()

		if readFileCloseErr != nil {
			fmt.Println("Error closing .clonignore file")
		}
	}

	filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			//contents, err := os.ReadFile(path)
			file, err := os.Open(path)
			reader := bufio.NewReader(file)
			content, _ := io.ReadAll(reader)
			encodedBase64 := base64.StdEncoding.EncodeToString(content)
			encryptedData := Encrypt(encodedBase64)

			if err == nil {
				// Should remote /remote from the beginning of the local path to conform with remote path
				localFilePathOnRemote := strings.Replace(path, localPathPrefix, "", 1)

				// If checksum does not match add this file to arrays with files t
				if !ValidateChecksum(encryptedData, remoteFiles[remotePathPrefix+localFilePathOnRemote]) {
					// Skip file which exists in .clonignore file
					if !ShouldIgnoreFile(clonignoreFilesPathsToSkip, path) {
						filesToUpdate = append(filesToUpdate, path)
					}

				}

				filteredRemoteFilesPaths := make([]string, 0)

				for _, remoteFilePath := range remoteFilesPaths {
					if remotePathPrefix+localFilePathOnRemote != remoteFilePath {
						filteredRemoteFilesPaths = append(filteredRemoteFilesPaths, remoteFilePath)
					}
				}

				remoteFilesPaths = filteredRemoteFilesPaths
			} else {
				ExitErrorf("Error reading file %q", path)
			}
		}

		return nil
	})

	checkFilesResult.FilesToUpload = filesToUpdate
	checkFilesResult.FilesToDelete = remoteFilesPaths

	return checkFilesResult
}

// Encrypt : Helper function for files encryption
// Inspired by Joseph Spurrier https://gist.github.com/josephspurrier/8304f09562d81babb494
func Encrypt(dataString string) []byte {
	data := []byte(dataString)
	key := []byte(GetEnvWithKey("AWS_ENCRYPTION_KEY"))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext
}

// Decrypt : Helper function for files decryption
// Inspired by Joseph Spurrier https://gist.github.com/josephspurrier/8304f09562d81babb494
func Decrypt(data []byte) []byte {
	key := []byte(GetEnvWithKey("AWS_ENCRYPTION_KEY"))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(data) < aes.BlockSize {
		panic("Text is too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data
}

// EncryptFile Read file, convert it to base64 and encrypt
func EncryptFile(filepath string) []byte {
	file, err := os.Open(filepath)
	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	encodedBase64 := base64.StdEncoding.EncodeToString(content)
	encryptedData := Encrypt(encodedBase64)

	if err != nil {
		ExitErrorf("An error with encrypting %q", filepath)
	}

	return encryptedData
}

// DecryptFile Read file, convert it to base64 and encrypt
func DecryptFile(buffer []byte) []byte {
	decrypted := Decrypt(buffer)

	decrypted = bytes.Trim(decrypted, "\x00")
	decryptedData := make([]byte, base64.StdEncoding.DecodedLen(len(decrypted)))

	_, _ = base64.StdEncoding.Decode(decryptedData, decrypted)

	return decryptedData
}
