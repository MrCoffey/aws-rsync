package s3

import (
	"fmt"
	"log"
	"os"

	"github.com/MrCoffey/s3-resync/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var secretKey string = os.Getenv("SECRET_KEY")
var accessKeyID string = os.Getenv("ACCESS_KEY_ID")
var zone string = os.Getenv("REGION")
var url string = os.Getenv("S3_ENDPOINT")

func GetObjects(legacyBucket string, destinationBucket string) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(zone),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretKey, ""),
		Endpoint:         aws.String(url),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	svc := s3.New(sess)

	config := &s3.ListObjectsV2Input{Bucket: aws.String(legacyBucket)}
	err = svc.ListObjectsV2Pages(config,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			findOrCreate(svc, page, legacyBucket, destinationBucket)
			return true
		})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v\n", legacyBucket, err)
	}
}

func findOrCreate(svc *s3.S3, page *s3.ListObjectsV2Output, legacyBucket string, destinationBucket string) {
	// fmt.Printf("%v", page)
	for _, item := range page.Contents {
		origin := fmt.Sprintf("%s/%s", legacyBucket, *item.Key)
		destination := fmt.Sprintf("%s/%s", destinationBucket, *item.Key)

		copyObject(svc, origin, destinationBucket, *item.Key)
		db.UpdateInDb(origin, destination)
	}
}

func copyObject(svc *s3.S3, origin string, destinationBucket string, key string) {
	fmt.Printf("Coping %s to %s/%s \n", origin, destinationBucket, key)

	copyObjectInput := &s3.CopyObjectInput{
		Bucket:     aws.String(destinationBucket),
		CopySource: aws.String(origin),
		Key:        aws.String(key),
	}

	result, err := svc.CopyObject(copyObjectInput)
	if err != nil {
		log.Fatal("Copy failed due to: \n", err)
	}

	fmt.Printf("%v \n", result)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
