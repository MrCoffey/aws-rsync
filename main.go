package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MrCoffey/s3-sync/config"
	"github.com/MrCoffey/s3-sync/db"
	"github.com/MrCoffey/s3-sync/s3"
	"github.com/urfave/cli/v2"
)

// Usage Example:
//    go run s3-resync.go -origin-bucket=LEGACY_BUCKET_NAME \
//		-destination-bucket=DESTINATION_BUCKET_NAME  \
//		-database-url=test:root123@/dbname \
//		-s3-secret-key=dasdas
//		-s3-access-key-id=dasdasd \
//		-s3-region=REGION \
//		-s3-endpoint=0.0.0.0:9000
func main() {
	var originBucket string
	var destinationBucket string
	var databaseURL string
	var s3SecretKey string
	var s3AccessKeyID string
	var s3Region string
	var s3Endpoint string
	var testMode bool

	app := &cli.App{
		Name:  "s3-sync",
		Usage: "Moves objects from one bucket to another using the s3 protocol.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "origin-bucket",
				Value:       "",
				Usage:       "Origin bucket name. (Required)",
				Destination: &originBucket,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "destination-bucket",
				Value:       "",
				Usage:       "Destination bucket name. (Required)",
				Destination: &destinationBucket,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "database-url",
				Value:       "",
				Usage:       "MariaDB URL. (Required)",
				Destination: &databaseURL,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "s3-secret-key",
				Value:       "",
				Usage:       "s3 secret key. (Required)",
				Destination: &s3SecretKey,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "s3-access-key-id",
				Value:       "",
				Usage:       "s3 secret key. (Required)",
				Destination: &s3AccessKeyID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "s3-region",
				Value:       "us-east-1",
				Usage:       "s3 Region. Default 'us-east-1'",
				Destination: &s3Region,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "s3-endpoint",
				Value:       "s3.amazon.com",
				Usage:       "s3 Endpoint. Default 's3.amazon.com'",
				Destination: &s3Endpoint,
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "test-mode",
				Value:       false,
				Usage:       "If enabled, this script will attempt to create the sql schema. Default false",
				Destination: &testMode,
				Required:    false,
			},
		},
		Action: func(ctx *cli.Context) error {
			configuration := config.Values{
				OriginBucket:      originBucket,
				DestinationBucket: destinationBucket,
				DatabaseURL:       databaseURL,
				S3SecretKey:       s3SecretKey,
				S3AccessKeyID:     s3AccessKeyID,
				S3Region:          s3Region,
				S3Endpoint:        s3Endpoint,
				TestMode:          testMode,
			}

			if configuration.TestMode {
				fmt.Printf("TEST MODE: Migrating database. \n")
				db.MigrateDB(&configuration)
			}

			s3.SyncObjects(&configuration)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
