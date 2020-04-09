package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MrCoffey/s3-resync/config"
	"github.com/MrCoffey/s3-resync/s3"
)

// Usage Example:
//    go run s3-resync.go -origin-bucket=LEGACY_BUCKET_NAME \
//		-destination-bucket=DESTINATION_BUCKET_NAME  \
//		-database-url=test:test@host \
//		-s3-secret-key=dasdas
//		-s3-access-key-id=dasdasd \
//		-s3-region=REGION \
//		-s3-endpoint=0.0.0.0:9000
func main() {
	originBucket := flag.String("origin-bucket", "", "Origin bucket name. (Required)")
	destinationBucket := flag.String("destination-bucket", "", "Destination bucket name. (Required)")
	databaseURL := flag.String("database-url", "", "Destination bucket name. (Required)")
	s3SecretKey := flag.String("s3-secret-key", "", "s3 secret key. (Required)")
	s3AccessKeyID := flag.String("s3-access-key-id", "", "s3 access key id. (Required)")

	s3Region := flag.String("s3-region", "us-east-1", "s3 Region.")
	s3Endpoint := flag.String("s3-endpoint", "s3.amazon.com", "s3 Endpoint.")
	flag.Parse()

	if *originBucket == "" || *destinationBucket == "" || *databaseURL == "" || *s3SecretKey == "" || *s3AccessKeyID == "" {
		fmt.Printf("Missing argument \n\n Usage: \n\n go run s3-resync.go")
		flag.PrintDefaults()
		os.Exit(1)
	}

	configuration := config.Values{
		OriginBucket:      *originBucket,
		DestinationBucket: *destinationBucket,
		DatabaseURL:       *databaseURL,
		S3SecretKey:       *s3SecretKey,
		S3AccessKeyID:     *s3AccessKeyID,
		S3Region:          *s3Region,
		S3Endpoint:        *s3Endpoint,
	}

	s3.SyncObjects(&configuration)
}
