package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func DownloadFile(sess *session.Session, bucket string, remoteFilePath string, localDirectoryPath string) {
	fileName := getFileNameByPath(remoteFilePath)
	localFilePath := localDirectoryPath + fileName
	file, err := os.Create(localFilePath)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

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

func UploadFile(sess *session.Session, bucket string, localFilePath string, remoteDirectoryPath string) {
	fileName := getFileNameByPath(localFilePath)
	remoteFilePath := remoteDirectoryPath + fileName
	file, err := os.Open(localFilePath)

	if err != nil {
		ExitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

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

func CreateBucket(sess *session.Session, bucketName string) {
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
}

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

func GetBucketsList(sess *session.Session) {
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)

	if err != nil {
		ExitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func GetBucketItems(sess *session.Session, bucket string) {
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})

	if err != nil {
		ExitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

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

	return nil
}
