package s3

import (
	"fmt"
	"log"
	"os"

	"github.com/MrCoffey/s3-resync/config"
	"github.com/MrCoffey/s3-resync/db"

	"github.com/jedib0t/go-pretty/table"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type reportData struct {
	ObjectsInLegacyBucket int `default:0`
	CopiedObjects         int `default:0`
	PathsInDdBefore       int `default:0`
	PathsInDdAfter        int `default:0`
	PathsUpdated          int `default:0`
	PathsSkipped          int `default:0`
}

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
	//db.MigrateDB(config)

	pageInx := 0
	err = svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{Bucket: aws.String(legacyBucket)},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageInx++
			copyAndUpdate(svc, config, page, pageInx)
			return true
		})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v\n\n", legacyBucket, err)
	}
}

func copyAndUpdate(svc *s3.S3, config *config.Values, page *s3.ListObjectsV2Output, pageInx int) {
	var legacyBucket string = config.OriginBucket
	var destinationBucket string = config.DestinationBucket
	var report reportData = reportData{ObjectsInLegacyBucket: len(page.Contents)}

	fmt.Printf("\nCurrently working on page No: %d Bucket Name: %s \n\n", pageInx, legacyBucket)

	for _, item := range page.Contents {
		newPathExist := db.FindPathInDb(config, destinationBucket, *item.Key)

		//db.CreateInDb(config, legacyBucket, *item.Key)

		if newPathExist != true {
			if moveObject(svc, config, *item.Key) {
				report.CopiedObjects++
			}
			if db.UpdateInDb(config, legacyBucket, destinationBucket, *item.Key) {
				report.PathsUpdated++
			}
		} else {
			report.PathsSkipped++
		}
	}

	showReport(&report, pageInx)
}

// Warning: This method only support coping files up to 5GB
// after that weigth is necesary to copy using a multipart option
func moveObject(svc *s3.S3, config *config.Values, key string) bool {
	// fmt.Printf("Coping %s to %s/%s \n\n", origin, destinationBucket, key)
	destinationBucket := config.DestinationBucket
	origin := fmt.Sprintf("%s/%s", config.DestinationBucket, key)

	copyObjectInput := &s3.CopyObjectInput{
		Bucket:     aws.String(destinationBucket),
		CopySource: aws.String(origin),
		Key:        aws.String(key),
		ACL:        aws.String("public-read"),
	}

	_, err := svc.CopyObject(copyObjectInput)
	if err != nil {
		log.Fatal("Copy failed due to: \n", err)
		return false
	}

	// deleteObjectInput := &s3.DeleteObjectInput{
	// 	Bucket: aws.String(origin),
	// 	Key:    aws.String(key),
	// }
	// _, err = svc.DeleteObject(deleteObjectInput)
	// if err != nil {
	// 	log.Fatal("Deleting process failed due to: \n", err)
	// 	return false
	// }

	return true
}

func showReport(report *reportData, pageInx int) {
	page := fmt.Sprintf("SYNC REPORT FOR PAGE: %d", pageInx)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{page, "COUNT"})
	t.AppendRows([]table.Row{
		{"LEGACY OBJECTS", report.ObjectsInLegacyBucket},
		{"COPIED OBJECTS", report.CopiedObjects},
		{"PATHS IN DB BEFORE", report.PathsInDdBefore},
		{"PATHS IN DB AFTER", report.PathsInDdAfter},
		{"PATHS UPDATED", report.PathsUpdated},
		{"PATHS SKIPPED", report.PathsSkipped},
	})
	t.Render()
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
