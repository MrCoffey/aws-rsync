package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetObjects(bucketName string) {
	SecretKey := os.Getenv("SECRET_KEY")
	AccessKeyID := os.Getenv("ACCESS_KEY_ID")
	zone := os.Getenv("REGION")
	url := os.Getenv("S3_ENDPOINT")

	fmt.Printf("ZONEEE %v", zone)

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(zone),
		Credentials:      credentials.NewStaticCredentials(AccessKeyID, SecretKey, ""),
		Endpoint:         aws.String(url),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})

	// Create S3 service client
	svc := s3.New(sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucketName, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucketName)
	fmt.Println("")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
