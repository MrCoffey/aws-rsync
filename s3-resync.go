package main

import (
	"fmt"
	"os"

	"github.com/MrCoffey/s3-resync/s3"
)

// Usage:
//    go run s3-resync.go BUCKET_NAME
func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s bucket_name",
			os.Args[0])
	}

	bucket := os.Args[1]

	s3.ListBuckets(bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
