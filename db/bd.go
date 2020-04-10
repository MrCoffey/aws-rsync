package db

import (
	"fmt"

	"github.com/MrCoffey/s3-resync/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Object is a S3 Model in database
type Object struct {
	gorm.Model
	Path   string
	Bucket string
}

func MigrateDB(config *config.Values) {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n\n")
	}
	defer db.Close()

	db.AutoMigrate(&Object{})

}

// ListPaths all the records
func ListPaths(config *config.Values) {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n\n")
	}
	defer db.Close()

	count := 0
	objects := db.Table("objects").Select("count(distinct(path))").Count(&count)

	fmt.Printf("%v \n\n", objects)
}

// FindPathInDb lookup for a provided path
func FindPathInDb(config *config.Values, bucket string, key string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database \n\n")
	}
	defer db.Close()

	if db.Where("bucket = ? AND path = ?", bucket, key).Find(&Object{}).RecordNotFound() {
		fmt.Printf("Failed to find KEY: %s %s \n", bucket, key)
		return false
	}

	return true
}

// UpdateInDb Updates a provided path
func UpdateInDb(config *config.Values, legacyBucket string, destinationBucket string, key string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()

	fmt.Printf("Updating path %s to: %s \n\n", legacyBucket, destinationBucket)

	object := Object{}

	if err := db.Model(&object).Updates(Object{Bucket: destinationBucket}).RecordNotFound(); err {
		fmt.Println("failed to update object with path: \n", legacyBucket)
		return false
	}
	return true
}

// CreateInDb Inserts a provided path
func CreateInDb(config *config.Values, bucketName string, key string) bool {
	var databaseURL string = config.DatabaseURL

	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		panic("failed to connect database\n\n")
	}
	defer db.Close()

	fmt.Printf("Inserting path %s/%s:\n\n", bucketName, key)

	object := Object{Path: key, Bucket: bucketName}
	if err := db.Create(&object).Error; err != nil {
		fmt.Println("failed to create object with path: \n\n", key)
		return false
	}

	return true
}
