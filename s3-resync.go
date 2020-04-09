package main

import (
	"fmt"
	"os"

	"github.com/MrCoffey/s3-resync/s3"
)

// Usage:
//    go run s3-resync.go LEGACY_BUCKET_NAME DESTINATION_BUCKET_NAME
func main() {
	if len(os.Args) != 3 {
		exitErrorf("Bucket name required\nUsage: go run s3-resync.go LEGACY_BUCKET_NAME DESTINATION_BUCKET_NAME")
	}
	legacyBucket := os.Args[1]
	destinationBucket := os.Args[2]

	s3.GetObjects(legacyBucket, destinationBucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
