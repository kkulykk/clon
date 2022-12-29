package services

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
)

// DownloadFile : AWS S3 helper method to download a file from a bucket
func DownloadFile(sess *session.Session, bucket string, remoteFilePath string, localDirectoryPath string) {
	fileName := GetFileNameByPath(remoteFilePath)
	localFilePath := localDirectoryPath + fileName

	downloader := s3manager.NewDownloader(sess)
	buf := aws.NewWriteAtBuffer([]byte{})
	numBytes, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(remoteFilePath),
		})

	decryptedData := DecryptFile(buf.Bytes())

	err = os.WriteFile(localFilePath, decryptedData, os.FileMode(0644))
	if err != nil {
		ExitErrorf("Error saving file %q, %v", fileName, err)
	}

	if err != nil {
		ExitErrorf("Unable to download item %q, %v", remoteFilePath, err)
	}

	fmt.Printf("Downloaded %q [%d bytes]\n", fileName, numBytes)
}

// UploadFileWithChecksum : AWS S3 helper method to upload new file to a bucket with upload file integrity check (checksum)
func UploadFileWithChecksum(sess *session.Session, bucket string, localFilePath string, remoteDirectoryPath string) {
	svc := s3.New(sess)
	fileName := GetFileNameByPath(localFilePath)
	remoteFilePath := remoteDirectoryPath + fileName

	encryptedFile := EncryptFile(localFilePath)

	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remoteFilePath),
		Body:   bytes.NewReader(encryptedFile),
	})

	if err != nil {
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

	remoteFileChecksum := strings.Trim(*(s3obj.ETag), "\"")
	localFileChecksum := fmt.Sprintf("%x", md5.Sum(encryptedFile))

	if localFileChecksum != remoteFileChecksum {
		ExitErrorf("Checksum mismatch of file %q on %q bucket", remoteFilePath, bucket)
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

// RemotePathExists : AWS S3 helper method to check if the path in the remote exists
func RemotePathExists(sess *session.Session, remotePath string) bool {
	svc := s3.New(sess)
	bucketName := GetBucketNameFromRemotePath(remotePath)
	remoteFilePath := GetRemoteFilePath(remotePath)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(remoteFilePath),
	})
	if err != nil {
		return false
	}
	return true
}

// GetAwsS3ItemMap : constructs and returns a map of keys (relative filenames) to checksums for the given
// bucket-configured s3 service. It is assumed  that the objects have not been multipart-uploaded,
// which will change the checksum.
func GetAwsS3ItemMap(sess *session.Session, bucket string, remotePath string) (map[string]string, error) {
	svc := s3.New(sess)

	remotePathPrefix := GetRemoteFilePathPrefix(remotePath)

	if !strings.HasSuffix(remotePathPrefix, "/") {
		remotePathPrefix = remotePathPrefix + "/"
	}

	loi := s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(GetRemoteFilePathPrefix(remotePath)),
	}

	obj, err := svc.ListObjects(&loi)

	var items = make(map[string]string)

	if err == nil {
		for _, s3obj := range obj.Contents {
			if !strings.HasSuffix(*(s3obj.Key), "/") {
				// Here we get the checksum
				eTag := strings.Trim(*(s3obj.ETag), "\"")
				items[*(s3obj.Key)] = eTag
			}
		}
		return items, nil
	}

	return nil, err
}
