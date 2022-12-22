package services

import (
	"crypto/md5"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"os"
	"strings"
)

// DownloadFile : AWS S3 helper method to download a file from a bucket
func DownloadFile(sess *session.Session, bucket string, remoteFilePath string, localDirectoryPath string) {
	fileName := GetFileNameByPath(remoteFilePath)
	localFilePath := localDirectoryPath + fileName
	file, err := os.Create(localFilePath)

	if err != nil {
		fmt.Println(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ExitErrorf("Unable to close the file")
		}
	}(file)

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(remoteFilePath),
		})

	if err != nil {
		ExitErrorf("Unable to download item %q, %v", remoteFilePath, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

// UploadFileWithChecksum : AWS S3 helper method to upload new file to a bucket with upload file integrity check (checksum)
func UploadFileWithChecksum(sess *session.Session, bucket string, localFilePath string, remoteDirectoryPath string) {
	svc := s3.New(sess)
	fileName := GetFileNameByPath(localFilePath)
	remoteFilePath := remoteDirectoryPath + fileName
	file, err := os.Open(localFilePath)

	if err != nil {
		ExitErrorf("Unable to open file %q, %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ExitErrorf("Unable to close the file")
		}
	}(file)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remoteFilePath),
		Body:   file,
	})

	if err != nil {
		// Print the error and exit.
		ExitErrorf("Unable to upload %q to %q, %v", remoteFilePath, bucket, err)
	}

	headObj := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remoteFilePath),
	}
	s3obj, err := svc.HeadObject(&headObj)

	if err != nil {
		ExitErrorf("Unable to get checksum of uploaded file")
	}

	fileContent, err := ioutil.ReadFile(localFilePath)
	remoteFileChecksum := strings.Trim(*(s3obj.ETag), "\"")
	localFileChecksum := fmt.Sprintf("%x", md5.Sum(fileContent))

	if localFileChecksum != remoteFileChecksum {
		ExitErrorf("Checksum mismatch of file %q on %q bucket", remoteFilePath, bucket)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", remoteFilePath, bucket)
}

// UploadFile : AWS S3 helper method to upload new file to a bucket
func UploadFile(sess *session.Session, bucket string, localFilePath string, remoteDirectoryPath string) {
	fileName := GetFileNameByPath(localFilePath)
	remoteFilePath := remoteDirectoryPath + fileName
	file, err := os.Open(localFilePath)

	if err != nil {
		ExitErrorf("Unable to open file %q, %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ExitErrorf("Unable to close the file")
		}
	}(file)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remoteFilePath),
		Body:   file,
	})

	if err != nil {
		// Print the error and exit.
		ExitErrorf("Unable to upload %q to %q, %v", remoteFilePath, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", remoteFilePath, bucket)
}

// CreateBucket : AWS S3 helper method to create a new bucket
func CreateBucket(sess *session.Session, bucketName string) error {
	svc := s3.New(sess)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		ExitErrorf("Unable to create bucket %q, %v", bucketName, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	fmt.Printf("Successfully created remote %q\n", bucketName)

	return nil
}

// DeleteBucketFile : AWS S3 helper method to delete file under the path in a bucket
// TODO: Remove two delete file functions
func DeleteBucketFile(sess *session.Session, bucket string, remoteFilePath string) {
	svc := s3.New(sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(remoteFilePath)})

	if err != nil {
		ExitErrorf("Unable to delete object %q from bucket %q, %v", remoteFilePath, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remoteFilePath),
	})

	if err == nil {
		fmt.Printf("File %q on bucket %q has been successfully deleted\n", remoteFilePath, bucket)
	}
}

// GetBucketsList : AWS S3 helper method to get list of available buckets
func GetBucketsList(sess *session.Session) {
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)

	if err != nil {
		ExitErrorf("Unable to list buckets, %v", err)
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

// GetBucketItems : AWS S3 helper method to get info about items under the path
func GetBucketItems(sess *session.Session, path string) {
	svc := s3.New(sess)
	bucket := GetBucketNameFromRemotePath(path)
	prefix := GetRemoteFilePathPrefix(path)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(prefix)})

	if err != nil {
		ExitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	fmt.Println("Found", len(resp.Contents), "items in remote", path)
	fmt.Println("")

	fmt.Println("--------------------------------------------------")

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", GetFileNameByPath(*item.Key))
		fmt.Println("Path:         ", *item.Key)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("--------------------------------------------------")
	}
}

// RemoveBucket : AWS S3 helper method to delete bucket instance
func RemoveBucket(sess *session.Session, bucket string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully deleted remote %q\n", bucket)

	return nil
}

// DeleteFile : AWS S3 helper method to delete file under the path in a bucket
func DeleteFile(sess *session.Session, path string) error {
	svc := s3.New(sess)

	bucket := GetBucketNameFromRemotePath(path)
	fileName := GetRemoteFilePath(path)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(fileName)})
	if err != nil {
		ExitErrorf("Unable to delete object %q from bucket %q, %v", fileName, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return err
	}

	fmt.Printf("File %q successfully deleted\n", fileName)

	return nil
}

// DeleteAll : AWS S3 helper method to delete all data under the path in a bucket
func DeleteAll(sess *session.Session, bucket string) error {
	svc := s3.New(sess)

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		ExitErrorf("Unable to delete objects from bucket %q, %v", bucket, err)
	}

	fmt.Printf("Deleted all object(s) from bucket: %s\n", bucket)

	return nil
}

// DeleteDirectory : AWS S3 helper method to delete directory in a bucket
func DeleteDirectory(sess *session.Session, path string) error {
	svc := s3.New(sess)

	bucket := GetBucketNameFromRemotePath(path)
	directory := GetRemoteFilePathPrefix(path)

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(directory),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		ExitErrorf("Unable to delete objects under given directory: %q, %v", directory, err)
	}

	fmt.Printf("Deleted all object(s) from directory: %q\n", directory)

	return nil
}

// GetBucketFileSize : AWS S3 helper method to get a size of a file in a bucket
func GetBucketFileSize(sess *session.Session, bucket string, remoteFilePath string) {
	svc := s3.New(sess)

	headObj := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remoteFilePath),
	}
	result, err := svc.HeadObject(&headObj)

	if err != nil {
		ExitErrorf("Unable to get size of file %q in %q, %v", remoteFilePath, bucket, err)
	}

	fmt.Printf("Size of %q (bucket: %q): %v bytes\n", remoteFilePath, bucket, aws.Int64Value(result.ContentLength))
}
