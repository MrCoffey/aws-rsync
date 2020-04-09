package s3

import (
	"fmt"
	"log"
	"os"

	"github.com/MrCoffey/s3-resync/config"
	"github.com/MrCoffey/s3-resync/db"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// SyncObjects fetch all (up to 1,000) the objects of the legacy Bucket using the aws pagination,
// then iterates in each key coping the object into the destination bucket using "public-read" ACL
// and updates it's path in database.
//
// Note: This operation can generate multiple requests to a s3 service.
func SyncObjects(config *config.Values) {
	var secretKey string = config.S3SecretKey
	var accessKeyID string = config.S3AccessKeyID
	var region string = config.S3Region
	var endpoint string = config.S3Endpoint
	var legacyBucket string = config.OriginBucket

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	svc := s3.New(sess)

	s3config := &s3.ListObjectsV2Input{Bucket: aws.String(legacyBucket)}
	err = svc.ListObjectsV2Pages(s3config,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			copyAndUpdate(svc, config, page)
			return true
		})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v\n\n", legacyBucket, err)
	}
}

func copyAndUpdate(svc *s3.S3, config *config.Values, page *s3.ListObjectsV2Output) {
	var legacyBucket string = config.OriginBucket
	var destinationBucket string = config.DestinationBucket

	// fmt.Printf("%v", page)
	for _, item := range page.Contents {
		origin := fmt.Sprintf("%s/%s", legacyBucket, *item.Key)
		destination := fmt.Sprintf("%s/%s", destinationBucket, *item.Key)

		copyObject(svc, origin, destinationBucket, *item.Key)
		db.UpdateInDb(config, origin, destination)
	}
}

// Warning: This method only support coping files up to 5GB
// after that weigth is necesary to copy using a multipart option
func copyObject(svc *s3.S3, origin string, destinationBucket string, key string) {
	fmt.Printf("Coping %s to %s/%s \n\n", origin, destinationBucket, key)

	copyObjectInput := &s3.CopyObjectInput{
		Bucket:     aws.String(destinationBucket),
		CopySource: aws.String(origin),
		Key:        aws.String(key),
		ACL:        aws.String("public-read"),
	}

	result, err := svc.CopyObject(copyObjectInput)
	if err != nil {
		log.Fatal("Copy failed due to: \n", err)
	}

	fmt.Printf("%v \n\n", result)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
