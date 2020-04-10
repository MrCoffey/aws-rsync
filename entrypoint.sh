#!/bin/sh

go run main.go \
    -origin-bucket=${LEGACY_BUCKET_NAME} \
    -destination-bucket=${NEW_BUCKET_NAME} \
    -database-url=${DATABASE_URL} \
    -s3-secret-key=${SECRET_KEY} \
    -s3-access-key-id=${ACCESS_KEY_ID} \
    -s3-region=${REGION} \
    -s3-endpoint=${S3_ENDPOINT} \
    -test-mode=true