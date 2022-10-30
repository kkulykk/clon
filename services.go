package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

// GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
