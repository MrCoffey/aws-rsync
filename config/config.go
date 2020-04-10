package config

type Values struct {
	OriginBucket      string
	DestinationBucket string
	DatabaseURL       string
	S3SecretKey       string
	S3AccessKeyID     string
	S3Region          string
	S3Endpoint        string
	TestMode          bool
}
