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
	secretKey := os.Getenv("SECRET_KEY")
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	zone := os.Getenv("REGION")
	url := os.Getenv("S3_ENDPOINT")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(zone),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretKey, ""),
		Endpoint:         aws.String(url),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	svc := s3.New(sess)

	config := &s3.ListObjectsV2Input{Bucket: aws.String(bucketName)}
	err = svc.ListObjectsV2Pages(config,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			// fmt.Println(page)
			findOrCreate(page, bucketName)
			return true
		})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucketName, err)
	}
}

func findOrCreate(page *s3.ListObjectsV2Output, bucketName string) {
	for _, item := range page.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Printf("Path:           %s-%v", bucketName, *item.Key)

		fmt.Println("")
	}

	fmt.Println("Found", *page.KeyCount, "items in bucket", bucketName)
	fmt.Println("")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
